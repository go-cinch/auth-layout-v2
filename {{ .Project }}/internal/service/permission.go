package service

import (
	"context"

	"{{ .Computed.common_module_final }}/copierx"
	"{{ .Computed.common_module_final }}/jwt"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_final }}"
	"{{ .Computed.module_name_final }}/internal/biz"
)

func (s *{{ .Computed.service_name_capitalized }}Service) Permission(ctx context.Context, req *v1.PermissionRequest) (*emptypb.Empty, error) {
	if s.user == nil || s.permission == nil {
		return nil, biz.ErrInternal(ctx, "permission service not configured")
	}

	u := jwt.FromServerContext(ctx)
	info := s.user.Info(ctx, u.Attrs["code"])
	if info == nil || info.Id == 0 {
		return nil, biz.ErrNoPermission(ctx)
	}

	check := &biz.CheckPermission{
		UserID: info.Id,
	}
	if req != nil {
		if req.Resource != nil {
			check.Resource = *req.Resource
		}
		if req.Method != nil {
			check.Method = *req.Method
		}
		if req.Uri != nil {
			check.URI = *req.Uri
		}
	}

	ok, err := s.permission.Check(ctx, check)
	if err != nil || !ok {
		return nil, biz.ErrNoPermission(ctx)
	}

	// Append minimal user attrs for auth_request integrations.
	jwt.AppendToReplyHeader(ctx, jwt.User{
		Attrs: map[string]string{
			"code":     info.Code,
			"platform": info.Platform,
		},
	})

	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) Info(ctx context.Context, _ *emptypb.Empty) (*v1.InfoReply, error) {
	if s.user == nil || s.permission == nil {
		return nil, biz.ErrInternal(ctx, "info service not configured")
	}

	u := jwt.FromServerContext(ctx)
	info := s.user.Info(ctx, u.Attrs["code"])
	if info == nil || info.Id == 0 {
		return nil, biz.ErrNoPermission(ctx)
	}

	rp := &v1.InfoReply{
		Permission: &v1.Permission{},
	}

	p, _ := s.permission.GetByUserID(ctx, info.Id)
	_ = copierx.Copy(&rp.Permission, p)
	_ = copierx.Copy(&rp, info)
	return rp, nil
}

