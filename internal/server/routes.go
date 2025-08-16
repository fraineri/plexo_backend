package server

import (
	"net/http"

	"github.com/fraineri/plexo_backend/internal/app_info/handlers"
	authHandlers "github.com/fraineri/plexo_backend/internal/auth/handlers"
)

func (s *Server) registerRoutes() {
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()

	// --- Initialize Handlers ---
	appInfoHandler := handlers.NewAppInfoHandler(s.db)
	authHandler := authHandlers.NewAuthHandler(s.db, s.settings)

	// Register public routes
	appInfoHandler.RegisterRoutes(apiRouter)
	authHandler.RegisterRoutes(apiRouter)

	// --- Protected Routes ---
	adminRouter := apiRouter.PathPrefix("/admin").Subrouter()
	// Apply the middleware
	adminRouter.Use(s.authMiddleware("admin.dashboard.read"))
	adminRouter.HandleFunc("/dashboard", s.handleAdminDashboard()).Methods("GET")
}

// handleAdminDashboard is an example of a protected handler.
func (s *Server) handleAdminDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Welcome to the protected admin dashboard!")); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}
}
