package service

import (
	"context"
	"net/http"

	"{{ .Computed.common_module_final }}/log"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthCheck handles the HTTP health endpoint.
func (*{{ .Computed.service_name_capitalized }}Service) HealthCheck(writer http.ResponseWriter, _ *http.Request) {
	log.Debug("healthcheck")
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("{}"))
}

// Check implements the standard gRPC health check.
func (*{{ .Computed.service_name_capitalized }}Service) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch streams health updates for the gRPC health endpoint.
func (*{{ .Computed.service_name_capitalized }}Service) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

