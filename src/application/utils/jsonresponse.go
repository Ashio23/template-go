package utils

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func XMLResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	xml.NewEncoder(w).Encode(data)
}
