package health

import (
	"net/http"
	"template-go/src/application/utils"
)

type HealthBffToFront struct {
	Status bool   `json:"status"`
	Path   string `json:"path"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	data := HealthBffToFront{Status: true, Path: r.URL.Path}
	utils.JSONResponse(w, data, http.StatusOK)
}
