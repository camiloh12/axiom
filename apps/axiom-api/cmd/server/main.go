package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/axiom-platform/axiom/apps/axiom-api/internal/gateway"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/identity"
	"github.com/axiom-platform/axiom/apps/axiom-api/internal/platform"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := platform.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := platform.NewDBPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("connected to database")

	var privKey *rsa.PrivateKey
	var pubKey *rsa.PublicKey
	if cfg.JWTPrivKey != "" {
		// TODO: parse PEM from config
		slog.Info("using configured JWT keys")
	}
	if privKey == nil {
		slog.Warn("generating ephemeral JWT keys — tokens will not survive restarts")
		privKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			slog.Error("failed to generate RSA key", "error", err)
			os.Exit(1)
		}
		pubKey = &privKey.PublicKey
	}

	jwtIssuer := identity.NewJWTIssuer(privKey, pubKey)
	identitySvc := identity.NewService(pool, jwtIssuer)
	identityHandler := identity.NewHandler(identitySvc, jwtIssuer)
	gw := gateway.NewMiddleware(jwtIssuer)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Location"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	identityHandler.RegisterRoutes(r)
	identityHandler.RegisterInvitationPublicRoutes(r)

	r.Group(func(r chi.Router) {
		r.Use(gw.Auth)
		identityHandler.RegisterAuthenticatedRoutes(r, gw)
	})

	slog.Info("starting server", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
