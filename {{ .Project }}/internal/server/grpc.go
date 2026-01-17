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
	"github.com/go-kratos/kratos/v2/transport/grpc"

	v1 "{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_final }}"
	"{{ .Computed.module_name_final }}/internal/conf"
	"{{ .Computed.module_name_final }}/internal/service"
)

// NewGRPCServer creates a gRPC server.
func NewGRPCServer(c *conf.Bootstrap, svc *service.{{ .Computed.service_name_capitalized }}Service) *grpc.Server {
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		tenantMiddleware.Tenant(), // Default required middleware for multi-tenancy
		{{ if .Computed.enable_i18n_final }}
		i18nMiddleware.Translator(i18n.WithLanguage(language.Make(c.Server.Language)), i18n.WithFs(locales)),
		{{ end }}
		ratelimit.Server(),
		logging.Server(),
		metadata.Server(),
	}
	{{ if .Computed.enable_trace_final }}
	if c.Tracer.Enable {
		middlewares = append(middlewares, tracing.Server(), traceMiddleware.ID())
	}
	{{ end }}
	if c.Server.Validate {
		middlewares = append(middlewares, validate.Validator())
	}

	var opts = []grpc.ServerOption{
		grpc.Middleware(middlewares...),
	}
	if c.Server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Server.Grpc.Network))
	}
	if c.Server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Server.Grpc.Addr))
	}
	if c.Server.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Server.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)
	v1.Register{{ .Computed.service_name_capitalized }}Server(srv, svc)

	return srv
}

