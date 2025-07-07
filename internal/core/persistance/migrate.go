package persistance

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) {
	log.Println("Checking for pending database migrations...")

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	migrationsPath := filepath.Join(basePath, "..", "..", "..", "migrations")
	sourceURL := "file://" + migrationsPath

	log.Printf("Looking for migrations in: %s", sourceURL)
	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil {

		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply.")
		} else {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
	} else {
		log.Println("Database migrations applied successfully.")
	}
}
