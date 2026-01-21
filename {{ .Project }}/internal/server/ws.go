package server

import (
	"net/http"

	"{{ .Computed.module_name_final }}/internal/service"
)

// NewWSHandler registers WebSocket routes and returns a ServeMux.
func NewWSHandler(svc *service.{{ .Computed.service_name_capitalized }}Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", svc.Ws)

	return mux
}
