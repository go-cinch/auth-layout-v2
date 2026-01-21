package server

import (
	{{ if .Computed.enable_i18n_final }}
	"embed"
	{{ end }}

	"github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)

{{ if .Computed.enable_i18n_final }}
//go:embed middleware/locales
var locales embed.FS
{{ end }}
