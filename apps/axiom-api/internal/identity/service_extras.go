package identity

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity/queries"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

type ClientDTO struct {
	ID                  uuid.UUID `json:"id"`
	FirmID              uuid.UUID `json:"firm_id"`
	Name                string    `json:"name"`
	Industry            string    `json:"industry,omitempty"`
	PrimaryContactEmail string    `json:"primary_contact_email,omitempty"`
}

type InvitationDTO struct {
	ID           uuid.UUID `json:"id"`
	FirmID       uuid.UUID `json:"firm_id"`
	Email        string    `json:"email"`
	AssignedRole string    `json:"assigned_role"`
	Status       string    `json:"status"`
	ExpiresAt    time.Time `json:"expires_at"`
	Token        string    `json:"token,omitempty"` // returned only on creation
}

type UpdateFirmInput struct {
	Name                *string
	LogoURL             *string
	Timezone            *string
	BillingContactEmail *string
}

type UpdateUserInput struct {
	DisplayName           *string
	Role                  *string
	NotificationFrequency *string
}

type CreateClientInput struct {
	Name                string
	Industry            string
	PrimaryContactEmail string
}

type UpdateClientInput struct {
	Name                *string
	Industry            *string
	PrimaryContactEmail *string
}

type CreateInvitationInput struct {
	Email        string
	AssignedRole string
	ExpiresIn    time.Duration // default: 7 days
}

// ctxWithFirm sets the RLS context variable for the current session.
func (s *Service) ctxWithFirm(ctx context.Context, firmID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, "SELECT set_config('app.current_firm_id', $1, false)", firmID.String())
	return err
}

func (s *Service) GetFirm(ctx context.Context, firmID uuid.UUID) (*FirmDTO, error) {
	if err := s.ctxWithFirm(ctx, firmID); err != nil {
		return nil, fmt.Errorf("set firm context: %w", err)
	}
	firm, err := s.queries.GetFirmByID(ctx, firmID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, platform.ErrNotFound("firm not found")
		}
		return nil, fmt.Errorf("get firm: %w", err)
	}
	return &FirmDTO{ID: firm.ID, Name: firm.Name, Slug: firm.Slug}, nil
}

func (s *Service) UpdateFirm(ctx context.Context, firmID uuid.UUID, input UpdateFirmInput) (*FirmDTO, error) {
	if err := s.ctxWithFirm(ctx, firmID); err != nil {
		return nil, fmt.Errorf("set firm context: %w", err)
	}
	params := queries.UpdateFirmParams{
		ID:                  firmID,
		Name:                pgTextPtr(input.Name),
		LogoUrl:             pgTextPtr(input.LogoURL),
		Timezone:            pgTextPtr(input.Timezone),
		BillingContactEmail: pgTextPtr(input.BillingContactEmail),
		Settings:            nil,
	}
	firm, err := s.queries.UpdateFirm(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, platform.ErrNotFound("firm not found")
		}
		return nil, fmt.Errorf("update firm: %w", err)
	}
	return &FirmDTO{ID: firm.ID, Name: firm.Name, Slug: firm.Slug}, nil
}

func (s *Service) ListUsers(ctx context.Context, firmID uuid.UUID, limit, offset int32) ([]UserDTO, error) {
	if err := s.ctxWithFirm(ctx, firmID); err != nil {
		return nil, fmt.Errorf("set firm context: %w", err)
	}
	rows, err := s.queries.ListUsersByFirmID(ctx, queries.ListUsersByFirmIDParams{
		FirmID: pgUUID(firmID),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	result := make([]UserDTO, 0, len(rows))
	for _, u := range rows {
		result = append(result, UserDTO{ID: u.ID, Email: u.Email, Name: u.DisplayName, Role: string(u.Role)})
	}
	return result, nil
}

func (s *Service) GetUser(ctx context.Context, userID uuid.UUID) (*UserDTO, error) {
	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, platform.ErrNotFound("user not found")
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &UserDTO{ID: user.ID, Email: user.Email, Name: user.DisplayName, Role: string(user.Role)}, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID uuid.UUID, input UpdateUserInput) (*UserDTO, error) {
	params := queries.UpdateUserParams{
		ID:          userID,
		DisplayName: pgTextPtr(input.DisplayName),
	}
	if input.Role != nil {
		params.Role = queries.NullUserRole{UserRole: queries.UserRole(*input.Role), Valid: true}
	}
	if input.NotificationFrequency != nil {
		params.NotificationFrequency = queries.NullNotificationFrequency{
			NotificationFrequency: queries.NotificationFrequency(*input.NotificationFrequency),
			Valid:                 true,
		}
	}
	user, err := s.queries.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, platform.ErrNotFound("user not found")
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &UserDTO{ID: user.ID, Email: user.Email, Name: user.DisplayName, Role: string(user.Role)}, nil
}

func (s *Service) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.queries.DeactivateUser(ctx, userID); err != nil {
		return fmt.Errorf("deactivate user: %w", err)
	}
	return nil
}

func (s *Service) CreateClient(ctx context.Context, firmID uuid.UUID, input CreateClientInput) (*ClientDTO, error) {
	if err := s.ctxWithFirm(ctx, firmID); err != nil {
		return nil, fmt.Errorf("set firm context: %w", err)
	}
	client, err := s.queries.CreateClient(ctx, queries.CreateClientParams{
		FirmID:              firmID,
		Name:                input.Name,
		Industry:            pgText(input.Industry),
		PrimaryContactEmail: pgText(input.PrimaryContactEmail),
	})
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	return clientToDTO(client), nil
}

func (s *Service) ListClients(ctx context.Context, firmID uuid.UUID, limit, offset int32) ([]ClientDTO, error) {
	if err := s.ctxWithFirm(ctx, firmID); err != nil {
		return nil, fmt.Errorf("set firm context: %w", err)
	}
	rows, err := s.queries.ListClientsByFirmID(ctx, queries.ListClientsByFirmIDParams{
		FirmID: firmID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list clients: %w", err)
	}
	result := make([]ClientDTO, 0, len(rows))
	for _, c := range rows {
		result = append(result, *clientToDTO(c))
	}
	return result, nil
}

func (s *Service) GetClient(ctx context.Context, clientID uuid.UUID) (*ClientDTO, error) {
	client, err := s.queries.GetClientByID(ctx, clientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, platform.ErrNotFound("client not found")
		}
		return nil, fmt.Errorf("get client: %w", err)
	}
	return clientToDTO(client), nil
}

func (s *Service) UpdateClient(ctx context.Context, clientID uuid.UUID, input UpdateClientInput) (*ClientDTO, error) {
	client, err := s.queries.UpdateClient(ctx, queries.UpdateClientParams{
		ID:                  clientID,
		Name:                pgTextPtr(input.Name),
		Industry:            pgTextPtr(input.Industry),
		PrimaryContactEmail: pgTextPtr(input.PrimaryContactEmail),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, platform.ErrNotFound("client not found")
		}
		return nil, fmt.Errorf("update client: %w", err)
	}
	return clientToDTO(client), nil
}

func (s *Service) CreateInvitation(ctx context.Context, firmID, inviterID uuid.UUID, input CreateInvitationInput) (*InvitationDTO, error) {
	// Launch posture (per docs/superpowers/specs/implementation-plan-design.md §2.1):
	// the auditee-side surface is gated behind CLIENT_HUB_ENABLED. Reject
	// ClientAdmin / ClientUser invitations while the flag is off.
	if !s.flags.ClientHubEnabled() && (input.AssignedRole == "ClientAdmin" || input.AssignedRole == "ClientUser") {
		return nil, platform.ErrValidation(
			"CLIENT_HUB_DISABLED",
			"Client roles cannot be invited while the Client Hub is disabled.",
		)
	}
	if err := s.ctxWithFirm(ctx, firmID); err != nil {
		return nil, fmt.Errorf("set firm context: %w", err)
	}
	token, tokenHash, err := generateInvitationToken()
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	expiresIn := input.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 7 * 24 * time.Hour
	}
	inv, err := s.queries.CreateInvitation(ctx, queries.CreateInvitationParams{
		FirmID:       firmID,
		Email:        input.Email,
		AssignedRole: queries.UserRole(input.AssignedRole),
		TokenHash:    tokenHash,
		ExpiresAt:    time.Now().Add(expiresIn),
		InvitedByID:  inviterID,
	})
	if err != nil {
		return nil, fmt.Errorf("create invitation: %w", err)
	}
	dto := invitationToDTO(inv)
	dto.Token = token
	return dto, nil
}

func (s *Service) ListInvitations(ctx context.Context, firmID uuid.UUID, limit, offset int32) ([]InvitationDTO, error) {
	if err := s.ctxWithFirm(ctx, firmID); err != nil {
		return nil, fmt.Errorf("set firm context: %w", err)
	}
	rows, err := s.queries.ListInvitationsByFirmID(ctx, queries.ListInvitationsByFirmIDParams{
		FirmID: firmID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list invitations: %w", err)
	}
	result := make([]InvitationDTO, 0, len(rows))
	for _, i := range rows {
		result = append(result, *invitationToDTO(i))
	}
	return result, nil
}

func (s *Service) ValidateInvitationToken(ctx context.Context, token string) (*InvitationDTO, error) {
	hash := hashToken(token)
	inv, err := s.queries.GetInvitationByTokenHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, platform.ErrNotFound("invitation not found")
		}
		return nil, fmt.Errorf("get invitation: %w", err)
	}
	if inv.Status != queries.InvitationStatusSent {
		return nil, platform.ErrBadRequest("invitation is no longer valid")
	}
	if time.Now().After(inv.ExpiresAt) {
		return nil, platform.ErrBadRequest("invitation has expired")
	}
	return invitationToDTO(inv), nil
}

type AcceptInvitationResult struct {
	User   UserDTO
	Tokens *TokenPair
}

func (s *Service) AcceptInvitation(ctx context.Context, token, displayName, password string) (*AcceptInvitationResult, error) {
	inv, err := s.ValidateInvitationToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := s.ctxWithFirm(ctx, inv.FirmID); err != nil {
		return nil, fmt.Errorf("set firm context: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.queries.WithTx(tx)
	if _, err := tx.Exec(ctx, "SELECT set_config('app.current_firm_id', $1, true)", inv.FirmID.String()); err != nil {
		return nil, fmt.Errorf("set firm context tx: %w", err)
	}

	user, err := qtx.CreateUser(ctx, queries.CreateUserParams{
		FirmID:                pgUUID(inv.FirmID),
		Email:                 inv.Email,
		DisplayName:           displayName,
		Role:                  queries.UserRole(inv.AssignedRole),
		AuthMethod:            queries.AuthMethodPassword,
		PasswordHash:          pgText(string(hash)),
		NotificationFrequency: queries.NotificationFrequencyDaily,
	})
	if err != nil {
		if isDuplicateKey(err) {
			return nil, platform.ErrConflict("email already registered")
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	if _, err := qtx.AcceptInvitation(ctx, inv.ID); err != nil {
		return nil, fmt.Errorf("accept invitation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	tokens, err := s.jwt.Issue(user.ID, inv.FirmID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("issue tokens: %w", err)
	}

	return &AcceptInvitationResult{
		User:   UserDTO{ID: user.ID, Email: user.Email, Name: user.DisplayName, Role: string(user.Role)},
		Tokens: tokens,
	}, nil
}

func (s *Service) CancelInvitation(ctx context.Context, invitationID uuid.UUID) error {
	if err := s.queries.CancelInvitation(ctx, invitationID); err != nil {
		return fmt.Errorf("cancel invitation: %w", err)
	}
	return nil
}

// --- helpers ---

func pgTextPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func clientToDTO(c queries.Client) *ClientDTO {
	out := &ClientDTO{ID: c.ID, FirmID: c.FirmID, Name: c.Name}
	if c.Industry.Valid {
		out.Industry = c.Industry.String
	}
	if c.PrimaryContactEmail.Valid {
		out.PrimaryContactEmail = c.PrimaryContactEmail.String
	}
	return out
}

func invitationToDTO(i queries.Invitation) *InvitationDTO {
	return &InvitationDTO{
		ID:           i.ID,
		FirmID:       i.FirmID,
		Email:        i.Email,
		AssignedRole: string(i.AssignedRole),
		Status:       string(i.Status),
		ExpiresAt:    i.ExpiresAt,
	}
}

func generateInvitationToken() (token, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	token = hex.EncodeToString(b)
	hash = hashToken(token)
	return token, hash, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
