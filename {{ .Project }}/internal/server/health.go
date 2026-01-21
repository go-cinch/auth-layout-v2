package server

import (
	"net/http"

	"{{ .Computed.module_name_final }}/internal/service"
)

// HealthHandler exposes the health endpoints served by the HTTP mux.
func HealthHandler(svc *service.{{ .Computed.service_name_capitalized }}Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", svc.HealthCheck)
	return mux
}
