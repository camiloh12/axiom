package identity

import (
	"encoding/json"
	"net/http"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
	jwt *JWTIssuer
}

func NewHandler(svc *Service, jwt *JWTIssuer) *Handler {
	return &Handler{svc: svc, jwt: jwt}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/api/v1/auth/register", h.register)
	r.Post("/api/v1/auth/login", h.login)
	r.Post("/api/v1/auth/refresh", h.refresh)
}

type registerRequest struct {
	FirmName          string   `json:"firm_name"`
	AdminEmail        string   `json:"admin_email"`
	AdminName         string   `json:"admin_name"`
	Password          string   `json:"password"`
	Country           string   `json:"country"`
	StaffCountRange   string   `json:"staff_count_range"`
	PrimaryAuditTypes []string `json:"primary_audit_types"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}

	if req.FirmName == "" || req.AdminEmail == "" || req.Password == "" {
		platform.WriteError(w, platform.ErrValidation("missing required fields", "firm_name, admin_email, and password are required"))
		return
	}

	result, err := h.svc.RegisterFirm(r.Context(), RegisterFirmInput{
		FirmName:   req.FirmName,
		AdminEmail: req.AdminEmail,
		AdminName:  req.AdminName,
		Password:   req.Password,
		Country:    req.Country,
		StaffCount: req.StaffCountRange,
		AuditTypes: req.PrimaryAuditTypes,
	})
	if err != nil {
		platform.WriteError(w, err)
		return
	}

	platform.WriteJSON(w, http.StatusCreated, map[string]any{
		"firm":          result.Firm,
		"user":          result.User,
		"access_token":  result.Tokens.AccessToken,
		"refresh_token": result.Tokens.RefreshToken,
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}

	tokens, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		platform.WriteError(w, err)
		return
	}

	platform.WriteJSON(w, http.StatusOK, tokens)
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		platform.WriteError(w, platform.ErrBadRequest("invalid request body"))
		return
	}

	tokens, err := h.jwt.Refresh(req.RefreshToken)
	if err != nil {
		platform.WriteError(w, platform.ErrUnauthorized("invalid refresh token"))
		return
	}

	platform.WriteJSON(w, http.StatusOK, tokens)
}
