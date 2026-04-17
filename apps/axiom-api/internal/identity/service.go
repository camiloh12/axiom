package identity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity/queries"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	pool    *pgxpool.Pool
	queries *queries.Queries
	jwt     *JWTIssuer
}

func NewService(pool *pgxpool.Pool, jwt *JWTIssuer) *Service {
	return &Service{
		pool:    pool,
		queries: queries.New(pool),
		jwt:     jwt,
	}
}

type RegisterFirmInput struct {
	FirmName   string
	AdminEmail string
	AdminName  string
	Password   string
	Country    string
	StaffCount string
	AuditTypes []string
}

type RegisterFirmResult struct {
	Firm   FirmDTO
	User   UserDTO
	Tokens *TokenPair
}

type FirmDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type UserDTO struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"display_name"`
	Role  string    `json:"role"`
}

func (s *Service) RegisterFirm(ctx context.Context, input RegisterFirmInput) (*RegisterFirmResult, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	slug := generateSlug(input.FirmName)

	auditTypesJSON, err := json.Marshal(input.AuditTypes)
	if err != nil {
		return nil, fmt.Errorf("marshal audit types: %w", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	firm, err := qtx.CreateFirm(ctx, queries.CreateFirmParams{
		Name:                input.FirmName,
		Slug:                slug,
		BillingContactEmail: input.AdminEmail,
		Country:             input.Country,
		StaffCountRange:     pgText(input.StaffCount),
		PrimaryAuditTypes:   auditTypesJSON,
	})
	if err != nil {
		if isDuplicateKey(err) {
			return nil, platform.ErrConflict("firm with this name already exists")
		}
		return nil, fmt.Errorf("create firm: %w", err)
	}

	if _, err := tx.Exec(ctx, "SELECT set_config('app.current_firm_id', $1, true)", firm.ID.String()); err != nil {
		return nil, fmt.Errorf("set firm context: %w", err)
	}

	user, err := qtx.CreateUser(ctx, queries.CreateUserParams{
		FirmID:                pgUUID(firm.ID),
		Email:                 input.AdminEmail,
		DisplayName:           input.AdminName,
		Role:                  queries.UserRoleFirmAdmin,
		AuthMethod:            queries.AuthMethodPassword,
		PasswordHash:          pgText(string(hash)),
		NotificationFrequency: queries.NotificationFrequencyRealTime,
	})
	if err != nil {
		if isDuplicateKey(err) {
			return nil, platform.ErrConflict("email already registered")
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	tokens, err := s.jwt.Issue(user.ID, firm.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("issue tokens: %w", err)
	}

	return &RegisterFirmResult{
		Firm:   FirmDTO{ID: firm.ID, Name: firm.Name, Slug: firm.Slug},
		User:   UserDTO{ID: user.ID, Email: user.Email, Name: user.DisplayName, Role: string(user.Role)},
		Tokens: tokens,
	}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, platform.ErrUnauthorized("invalid email or password")
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	if !user.IsActive {
		return nil, platform.ErrUnauthorized("account is deactivated")
	}

	if !user.PasswordHash.Valid {
		return nil, platform.ErrUnauthorized("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(password)); err != nil {
		return nil, platform.ErrUnauthorized("invalid email or password")
	}

	firmID := uuid.Nil
	if user.FirmID.Valid {
		firmID = user.FirmID.Bytes
	}

	return s.jwt.Issue(user.ID, firmID, string(user.Role))
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, slug)
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	slug = slug + "-" + uuid.New().String()[:8]
	return slug
}

func pgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func pgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func isDuplicateKey(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}
