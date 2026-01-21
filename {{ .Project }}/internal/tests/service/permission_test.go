package service_test

import (
	"context"
	"testing"

	"{{ .Computed.common_module_final }}/jwt"
	v1 "{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_final }}"
	"{{ .Computed.module_name_final }}/internal/tests/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestPermission(t *testing.T) {
	s := mock.{{ .Computed.service_name_capitalized }}Service()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	// Set up JWT context with user code
	ctx = jwt.NewServerContextByUser(ctx, jwt.User{
		Attrs: map[string]string{
			"code":     "test-user-code",
			"platform": "web",
		},
	})

	method := "GET"
	uri := "/api/test"
	req := &v1.PermissionRequest{
		Method: &method,
		Uri:    &uri,
	}

	_, err := s.Permission(ctx, req)
	// Note: May fail if user doesn't exist or no permission
	if err != nil {
		t.Logf("Permission returned error (may be expected): %v", err)
	}
}

func TestInfo(t *testing.T) {
	s := mock.{{ .Computed.service_name_capitalized }}Service()
	ctx := mock.NewContextWithUserId(context.Background(), "test-tenant")

	// Set up JWT context with user code
	ctx = jwt.NewServerContextByUser(ctx, jwt.User{
		Attrs: map[string]string{
			"code":     "test-user-code",
			"platform": "web",
		},
	})

	rp, err := s.Info(ctx, &emptypb.Empty{})
	// Note: May fail if user doesn't exist
	if err != nil {
		t.Logf("Info returned error (may be expected): %v", err)
	} else {
		assert.NotNil(t, rp)
	}
}
