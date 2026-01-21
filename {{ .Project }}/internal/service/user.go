package service

import (
	"context"
	{{- if .Computed.enable_user_lock_final }}
	"strings"
	{{- end }}

	"go.opentelemetry.io/otel"

	"{{ .Computed.common_module_final }}/copierx"
	params "{{ .Computed.common_module_final }}/proto/params"
	"{{ .Computed.common_module_final }}/utils"
	{{- if .Computed.enable_user_lock_final }}
	"github.com/golang-module/carbon/v2"
	{{- end }}
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_final }}"
	"{{ .Computed.module_name_final }}/internal/biz"
)

func (s *{{ .Computed.service_name_capitalized }}Service) FindUser(ctx context.Context, req *v1.FindUserRequest) (*v1.FindUserReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "FindUser")
	defer span.End()

	r := &biz.FindUser{}
	copierx.Copy(&r.Page, req.GetPage())
	copierx.Copy(r, req)

	list, err := s.user.Find(ctx, r)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindUserReply{
		Page: &params.Page{},
	}
	copierx.Copy(&rp.Page, &r.Page)
	copierx.Copy(&rp.List, list)
	return rp, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "UpdateUser")
	defer span.End()

	u := &biz.UpdateUser{}
	copierx.Copy(u, req)
	if req.Locked != nil {
		var locked int16
		if *req.Locked {
			locked = 1
		}
		u.Locked = &locked
	}
	{{- if .Computed.enable_user_lock_final }}
	// Parse lock_expire_time string if provided (priority over lock_expire).
	if req.LockExpireTime != nil && *req.LockExpireTime != "" {
		if expire, ok := parseLockExpireTime(*req.LockExpireTime); ok {
			u.LockExpire = &expire
		}
	}
	{{- end }}

	if err := s.user.Update(ctx, u); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) DeleteUser(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "DeleteUser")
	defer span.End()

	if err := s.user.Delete(ctx, utils.Str2Int64Arr(req.GetIds())...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

{{- if .Computed.enable_user_lock_final }}
func parseLockExpireTime(s string) (int64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, true
	}
	parsed := carbon.Parse(s)
	if parsed.Error == nil {
		return parsed.StdTime().Unix(), true
	}
	return 0, false
}
{{- end }}
