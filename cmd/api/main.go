package main

import (
	"log"

	"github.com/fraineri/plexo_backend/internal/core/persistance"
	"github.com/fraineri/plexo_backend/internal/core/settings"
	"github.com/fraineri/plexo_backend/internal/server"
)

func main() {
	log.Println("Loading settings...")
	settings, err := settings.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading settings: %v", err)
	}
	log.Printf(
		"Settings loaded successfully. [DEBUG = %v, ENVIRONMENT = %s]",
		settings.Project.Debug, settings.Project.Environment,
	)

	persistance.RunMigrations(settings.Database.GetPostgresMigrationURL())

	log.Println("Initializing database connection...")
	db, err := persistance.NewConnection(settings.Database.GetPostgresWriteDSN())
	if err != nil {
		log.Fatalf("Error initializing database connection: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()
	log.Println("Database connection initialized successfully.")

	// Create and start the server
	srv := server.NewServer(db, settings)
	srv.Start(":8080")
}
