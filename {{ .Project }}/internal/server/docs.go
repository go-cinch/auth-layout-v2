package server

import (
	"net/http"
)

// DocsHandler serves static files from docs directory.
func DocsHandler() http.Handler {
	// Service runs from cmd/{{ .Computed.service_name_final }}, so docs is at ../../docs
	return http.StripPrefix("/docs/", http.FileServer(http.Dir("../../docs")))
}
