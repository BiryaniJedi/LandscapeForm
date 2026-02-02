package	db

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
)

func TestDB(t *testing.T) *sql.DB {
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

	// Drop all objects in the public schema (preserves extensions and their functions)
	_, err = db.Exec(`
		DO $$ DECLARE
			r RECORD;
		BEGIN
			-- Drop all tables
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
				EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
			-- Drop all user-defined functions (exclude extension functions)
			FOR r IN (SELECT p.proname, pg_get_function_identity_arguments(p.oid) as args
					  FROM pg_proc p
					  INNER JOIN pg_namespace ns ON (p.pronamespace = ns.oid)
					  LEFT JOIN pg_depend d ON (d.objid = p.oid AND d.deptype = 'e')
					  WHERE ns.nspname = 'public'
					  AND p.prokind = 'f'
					  AND d.objid IS NULL) LOOP
				EXECUTE 'DROP FUNCTION IF EXISTS ' || quote_ident(r.proname) || '(' || r.args || ') CASCADE';
			END LOOP;
		END $$;
	`)
	if err != nil {
		t.Fatalf("reset schema: %v", err)
	}

	// Load schema.sql
	schemaPath := filepath.Join("..", "..", "db", "db", "init", "schema.sql")
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
