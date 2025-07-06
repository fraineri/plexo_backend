package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fraineri/plexo_backend/internal/app_info/entities"
	"github.com/fraineri/plexo_backend/internal/app_info/handlers/dtos"
	"github.com/fraineri/plexo_backend/internal/app_info/usecases"
	"github.com/gorilla/mux"
)

type AppInfo = entities.AppInfo
type AppInfoResponseDTO = dtos.AppInfoResponseDTO

type AppInfoHandler struct {
	getAppInfoUsecase usecases.GetAppInfoUsecase
	disableAppUsecase usecases.DisableAppUsecase
	enableAppUsecase  usecases.EnableAppUsecase
}

func NewAppInfoHandler(
	getAppInfoUsecase usecases.GetAppInfoUsecase,
	disableAppUsecase usecases.DisableAppUsecase,
	enableAppUsecase usecases.EnableAppUsecase,
) *AppInfoHandler {
	return &AppInfoHandler{
		getAppInfoUsecase: getAppInfoUsecase,
		disableAppUsecase: disableAppUsecase,
		enableAppUsecase:  enableAppUsecase,
	}
}

func (h *AppInfoHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/app-info", h.getAppInfo).Methods("GET")
	router.HandleFunc("/app-info/disable", h.disableApp).Methods("POST")
	router.HandleFunc("/app-info/enable", h.enableApp).Methods("POST")
}

func (h *AppInfoHandler) getAppInfo(w http.ResponseWriter, r *http.Request) {
	result, err := h.getAppInfoUsecase.Execute(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get app info: %v", err), http.StatusInternalServerError)
		return
	}

	response := AppInfoResponseDTO{
		Name:    result.Name,
		Version: result.Version,
		Status:  result.Status,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func (h *AppInfoHandler) disableApp(w http.ResponseWriter, r *http.Request) {
	result, err := h.disableAppUsecase.Execute(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to disable app: %v", err), http.StatusInternalServerError)
		return
	}

	response := AppInfoResponseDTO{
		Name:    result.Name,
		Version: result.Version,
		Status:  result.Status,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func (h *AppInfoHandler) enableApp(w http.ResponseWriter, r *http.Request) {
	result, err := h.enableAppUsecase.Execute(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to enable app: %v", err), http.StatusInternalServerError)
		return
	}

	response := AppInfoResponseDTO{
		Name:    result.Name,
		Version: result.Version,
		Status:  result.Status,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
