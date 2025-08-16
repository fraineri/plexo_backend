package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fraineri/plexo_backend/internal/auth/persistance"
	"github.com/fraineri/plexo_backend/internal/auth/services"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
	"github.com/fraineri/plexo_backend/internal/core/settings"
	"github.com/gorilla/mux"
)

type Server struct {
	db       *sql.DB
	settings *settings.Settings
	router   *mux.Router

	jwtService  services.JWTService
	authService services.AuthorizationService
}

func NewServer(db *sql.DB, settings *settings.Settings) *Server {
	s := &Server{
		db:       db,
		settings: settings,
		router:   mux.NewRouter(),
	}

	// --- Initialize Services ---
	unitOfWork := uow.NewUnitOfWork(db)
	userRoleRepo := persistance.NewUserRoleRepository(unitOfWork)
	rolePermissionRepo := persistance.NewRolePermissionRepository(unitOfWork)

	s.jwtService = services.NewJWTService(settings.JWT.SecretKey, settings.JWT.SessionTimeMins)
	s.authService = services.NewAuthorizationService(userRoleRepo, rolePermissionRepo)

	// Register all routes
	s.registerRoutes()

	return s
}

// Start runs the HTTP server.
func (s *Server) Start(addr string) {
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, s.router))
}
