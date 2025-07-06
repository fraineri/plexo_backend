package dtos

type AppInfoResponseDTO struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}
