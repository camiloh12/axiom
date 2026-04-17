package platform

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestDB creates a temporary test database, runs migrations, and returns a pool.
// The database is dropped when the test completes.
func TestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	adminURL := "postgres://axiom_svc:localdev@localhost:5432/axiom_db?sslmode=disable"
	adminPool, err := pgxpool.New(ctx, adminURL)
	if err != nil {
		t.Fatalf("connect to admin db: %v", err)
	}

	dbName := fmt.Sprintf("axiom_test_%s", t.Name())
	cleanName := ""
	for _, c := range dbName {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
			cleanName += string(c)
		} else if c >= 'A' && c <= 'Z' {
			cleanName += string(c - 'A' + 'a')
		} else {
			cleanName += "_"
		}
	}
	// Postgres identifiers are limited to 63 bytes.
	if len(cleanName) > 63 {
		cleanName = cleanName[:63]
	}
	dbName = cleanName

	_, _ = adminPool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	_, err = adminPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("create test database: %v", err)
	}
	adminPool.Close()

	testURL := fmt.Sprintf("postgres://axiom_svc:localdev@localhost:5432/%s?sslmode=disable", dbName)

	m, err := migrate.New("file://../../migrations", testURL)
	if err != nil {
		t.Fatalf("create migrator: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("run migrations: %v", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		t.Fatalf("close migrator source: %v", srcErr)
	}
	if dbErr != nil {
		t.Fatalf("close migrator db: %v", dbErr)
	}

	pool, err := pgxpool.New(ctx, testURL)
	if err != nil {
		t.Fatalf("connect to test db: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
		cleanupPool, err := pgxpool.New(context.Background(), adminURL)
		if err == nil {
			_, _ = cleanupPool.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
			cleanupPool.Close()
		}
	})

	return pool
}
