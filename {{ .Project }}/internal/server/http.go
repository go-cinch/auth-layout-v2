package server

import (
	tenantMiddleware "{{ .Computed.common_module_final }}/middleware/tenant/v2"
	{{ if .Computed.enable_i18n_final }}
	"{{ .Computed.common_module_final }}/i18n"
	i18nMiddleware "{{ .Computed.common_module_final }}/middleware/i18n"
	"golang.org/x/text/language"
	{{ end }}
	"{{ .Computed.common_module_final }}/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	{{ if .Computed.enable_trace_final }}
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	traceMiddleware "{{ .Computed.common_module_final }}/middleware/trace"
	{{ end }}
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"

	v1 "{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_final }}"
	"{{ .Computed.module_name_final }}/internal/conf"
	{{ if contains "header" .Computed.middlewares_final }}
	localMiddleware "{{ .Computed.module_name_final }}/internal/server/middleware"
	{{ end }}
	"{{ .Computed.module_name_final }}/internal/service"
)

// NewHTTPServer creates an HTTP server.
func NewHTTPServer(c *conf.Bootstrap, svc *service.{{ .Computed.service_name_capitalized }}Service) *http.Server {
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		tenantMiddleware.Tenant(), // Default required middleware for multi-tenancy
		{{ if .Computed.enable_i18n_final }}
		i18nMiddleware.Translator(i18n.WithLanguage(language.Make(c.Server.Language)), i18n.WithFs(locales)),
		{{ end }}
		ratelimit.Server(),
		{{ if contains "header" .Computed.middlewares_final }}
		localMiddleware.Header(),
		{{ end }}
	}
	{{ if .Computed.enable_trace_final }}
	if c.Tracer.Enable {
		middlewares = append(middlewares, tracing.Server(), traceMiddleware.ID())
	}
	{{ end }}
	middlewares = append(middlewares, logging.Server(), metadata.Server())
	if c.Server.Validate {
		middlewares = append(middlewares, validate.Validator())
	}

	var opts = []http.ServerOption{
		http.Middleware(middlewares...),
	}
	if c.Server.Http.Network != "" {
		opts = append(opts, http.Network(c.Server.Http.Network))
	}
	if c.Server.Http.Addr != "" {
		opts = append(opts, http.Address(c.Server.Http.Addr))
	}
	if c.Server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Server.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)
	v1.Register{{ .Computed.service_name_capitalized }}HTTPServer(srv, svc)
	{{ if .Computed.enable_health_check_final }}
	srv.HandlePrefix("/healthz", HealthHandler(svc))
	{{ end }}
	if c.Server.Http.Docs {
		srv.HandlePrefix("/docs/", DocsHandler())
	}
	if c.Server.EnablePprof {
		srv.HandlePrefix("/debug/pprof", pprof.NewHandler())
	}
	{{ if .Computed.enable_ws_final }}
	srv.HandlePrefix("/ws", NewWSHandler(svc))
	{{ end }}

	return srv
}

