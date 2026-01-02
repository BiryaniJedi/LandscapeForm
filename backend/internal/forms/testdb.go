package forms

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL not set (did you load .env.testing?)")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	// Ensure DB is reachable
	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	// Reset schema (hard reset)
	_, err = db.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
	if err != nil {
		t.Fatalf("reset schema: %v", err)
	}

	// Load schema.sql
	schemaPath := filepath.Join("..", "..", "db", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema.sql: %v", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		t.Fatalf("apply schema: %v", err)
	}

	// Cleanup when test finishes
	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}
