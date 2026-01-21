package service

import (
	"context"

	"go.opentelemetry.io/otel"

	"{{ .Computed.common_module_final }}/copierx"
	params "{{ .Computed.common_module_final }}/proto/params"
	"{{ .Computed.common_module_final }}/utils"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_final }}"
	"{{ .Computed.module_name_final }}/internal/biz"
)

func (s *{{ .Computed.service_name_capitalized }}Service) CreateWhitelist(ctx context.Context, req *v1.CreateWhitelistRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "CreateWhitelist")
	defer span.End()

	w := &biz.Whitelist{}
	copierx.Copy(w, req)
	if err := s.whitelist.Create(ctx, w); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) FindWhitelist(ctx context.Context, req *v1.FindWhitelistRequest) (*v1.FindWhitelistReply, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "FindWhitelist")
	defer span.End()

	r := &biz.FindWhitelist{}
	copierx.Copy(&r.Page, req.GetPage())
	copierx.Copy(r, req)

	list, err := s.whitelist.Find(ctx, r)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindWhitelistReply{
		Page: &params.Page{},
	}
	copierx.Copy(&rp.Page, &r.Page)
	copierx.Copy(&rp.List, list)
	return rp, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) UpdateWhitelist(ctx context.Context, req *v1.UpdateWhitelistRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "UpdateWhitelist")
	defer span.End()

	w := &biz.UpdateWhitelist{}
	copierx.Copy(w, req)

	if err := s.whitelist.Update(ctx, w); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) DeleteWhitelist(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	tr := otel.Tracer("service")
	ctx, span := tr.Start(ctx, "DeleteWhitelist")
	defer span.End()

	if err := s.whitelist.Delete(ctx, utils.Str2Int64Arr(req.GetIds())...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}
