package identity

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

// RoleGuard produces middleware restricting access to the listed roles.
// Implemented by gateway.Middleware; abstracted here to avoid an import cycle.
type RoleGuard interface {
	WithRole(roles ...string) func(http.Handler) http.Handler
}

// RegisterAuthenticatedRoutes mounts JWT-protected endpoints. The caller is
// expected to wrap this group with the Auth middleware.
func (h *Handler) RegisterAuthenticatedRoutes(r chi.Router, gw RoleGuard) {
	// Firm
	r.Get("/api/v1/firms/current", h.getCurrentFirm)
	r.Patch("/api/v1/firms/current", h.updateCurrentFirm)

	// Users
	r.Get("/api/v1/users", h.listUsers)
	r.Get("/api/v1/users/me", h.getMe)
	r.Patch("/api/v1/users/me", h.updateMe)
	r.Get("/api/v1/users/{userId}", h.getUser)
	r.With(gw.WithRole("FirmAdmin")).Patch("/api/v1/users/{userId}", h.updateUser)
	r.With(gw.WithRole("FirmAdmin")).Post("/api/v1/users/{userId}/deactivate", h.deactivateUser)

	// Clients
	r.Get("/api/v1/clients", h.listClients)
	r.Post("/api/v1/clients", h.createClient)
	r.Get("/api/v1/clients/{clientId}", h.getClient)
	r.Patch("/api/v1/clients/{clientId}", h.updateClient)

	// Invitations (FirmAdmin-only create/list/cancel; public token endpoints mounted separately)
	r.With(gw.WithRole("FirmAdmin")).Get("/api/v1/invitations", h.listInvitations)
	r.With(gw.WithRole("FirmAdmin")).Post("/api/v1/invitations", h.createInvitation)
	r.With(gw.WithRole("FirmAdmin")).Delete("/api/v1/invitations/{invitationId}", h.cancelInvitation)
}

// RegisterInvitationPublicRoutes mounts invitation-token endpoints that do not
// require a JWT (used by invitees before they have an account).
func (h *Handler) RegisterInvitationPublicRoutes(r chi.Router) {
	r.Get("/api/v1/invitations/validate/{token}", h.validateInvitation)
	r.Post("/api/v1/invitations/accept", h.acceptInvitation)
}

func (h *Handler) getCurrentFirm(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	firm, err := h.svc.GetFirm(r.Context(), claims.FirmID)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, firm)
}

type updateFirmRequest struct {
	Name                *string `json:"name"`
	LogoURL             *string `json:"logo_url"`
	Timezone            *string `json:"timezone"`
	BillingContactEmail *string `json:"billing_contact_email"`
}

func (h *Handler) updateCurrentFirm(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	var req updateFirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}
	firm, err := h.svc.UpdateFirm(r.Context(), claims.FirmID, UpdateFirmInput(req))
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, firm)
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	limit, offset := parsePagination(r)
	users, err := h.svc.ListUsers(r.Context(), claims.FirmID, limit, offset)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"items": users})
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	user, err := h.svc.GetUser(r.Context(), claims.UserID)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, user)
}

type updateUserRequest struct {
	DisplayName           *string `json:"display_name"`
	Role                  *string `json:"role"`
	NotificationFrequency *string `json:"notification_frequency"`
}

func (h *Handler) updateMe(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}
	// Users cannot change their own role.
	req.Role = nil
	user, err := h.svc.UpdateUser(r.Context(), claims.UserID, UpdateUserInput{
		DisplayName:           req.DisplayName,
		NotificationFrequency: req.NotificationFrequency,
	})
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid user id"))
		return
	}
	user, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid user id"))
		return
	}
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}
	user, err := h.svc.UpdateUser(r.Context(), id, UpdateUserInput(req))
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) deactivateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid user id"))
		return
	}
	if err := h.svc.DeactivateUser(r.Context(), id); err != nil {
		platform.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type createClientRequest struct {
	Name                string `json:"name"`
	Industry            string `json:"industry"`
	PrimaryContactEmail string `json:"primary_contact_email"`
}

func (h *Handler) createClient(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	var req createClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}
	if req.Name == "" {
		platform.WriteError(w, platform.ErrValidation("missing required fields", "name is required"))
		return
	}
	client, err := h.svc.CreateClient(r.Context(), claims.FirmID, CreateClientInput(req))
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusCreated, client)
}

func (h *Handler) listClients(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	limit, offset := parsePagination(r)
	clients, err := h.svc.ListClients(r.Context(), claims.FirmID, limit, offset)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"items": clients})
}

func (h *Handler) getClient(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "clientId"))
	if err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid client id"))
		return
	}
	client, err := h.svc.GetClient(r.Context(), id)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, client)
}

type updateClientRequest struct {
	Name                *string `json:"name"`
	Industry            *string `json:"industry"`
	PrimaryContactEmail *string `json:"primary_contact_email"`
}

func (h *Handler) updateClient(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "clientId"))
	if err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid client id"))
		return
	}
	var req updateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}
	client, err := h.svc.UpdateClient(r.Context(), id, UpdateClientInput(req))
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, client)
}

type createInvitationRequest struct {
	Email        string `json:"email"`
	AssignedRole string `json:"assigned_role"`
}

func (h *Handler) createInvitation(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	var req createInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}
	if req.Email == "" || req.AssignedRole == "" {
		platform.WriteError(w, platform.ErrValidation("missing required fields", "email and assigned_role are required"))
		return
	}
	inv, err := h.svc.CreateInvitation(r.Context(), claims.FirmID, claims.UserID, CreateInvitationInput{
		Email:        req.Email,
		AssignedRole: req.AssignedRole,
	})
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusCreated, inv)
}

func (h *Handler) listInvitations(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	limit, offset := parsePagination(r)
	invs, err := h.svc.ListInvitations(r.Context(), claims.FirmID, limit, offset)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, map[string]any{"items": invs})
}

func (h *Handler) cancelInvitation(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "invitationId"))
	if err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid invitation id"))
		return
	}
	if err := h.svc.CancelInvitation(r.Context(), id); err != nil {
		platform.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) validateInvitation(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	inv, err := h.svc.ValidateInvitationToken(r.Context(), token)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusOK, inv)
}

type acceptInvitationRequest struct {
	Token       string `json:"token"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

func (h *Handler) acceptInvitation(w http.ResponseWriter, r *http.Request) {
	var req acceptInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}
	if req.Token == "" || req.Password == "" {
		platform.WriteError(w, platform.ErrValidation("missing required fields", "token and password are required"))
		return
	}
	result, err := h.svc.AcceptInvitation(r.Context(), req.Token, req.DisplayName, req.Password)
	if err != nil {
		platform.WriteError(w, err)
		return
	}
	platform.WriteJSON(w, http.StatusCreated, map[string]any{
		"user":          result.User,
		"access_token":  result.Tokens.AccessToken,
		"refresh_token": result.Tokens.RefreshToken,
	})
}

func parsePagination(r *http.Request) (int32, int32) {
	limit := int32(50)
	offset := int32(0)
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = int32(n)
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 && n <= math.MaxInt32 {
			offset = int32(n)
		}
	}
	return limit, offset
}
