package service

import (
	"context"

	"{{ .Computed.common_module_final }}/copierx"
	params "{{ .Computed.common_module_final }}/proto/params"
	"{{ .Computed.common_module_final }}/utils"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_final }}"
	"{{ .Computed.module_name_final }}/internal/biz"
)

func (s *{{ .Computed.service_name_capitalized }}Service) CreateWhitelist(ctx context.Context, req *v1.CreateWhitelistRequest) (*emptypb.Empty, error) {
	if s.whitelist == nil {
		return nil, biz.ErrInternal(ctx, "whitelist usecase not configured")
	}
	if req == nil {
		return nil, biz.ErrIllegalParameter(ctx, "whitelist")
	}
	if req.Category > int32(^uint16(0)>>1) {
		return nil, biz.ErrIllegalParameter(ctx, "category")
	}

	w := &biz.Whitelist{
		Category: int16(req.Category),
		Resource: req.Resource,
	}
	if err := s.whitelist.Create(ctx, w); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) FindWhitelist(ctx context.Context, req *v1.FindWhitelistRequest) (*v1.FindWhitelistReply, error) {
	if s.whitelist == nil {
		return nil, biz.ErrInternal(ctx, "whitelist usecase not configured")
	}

	cond := &biz.FindWhitelist{}
	if req != nil {
		if req.Page != nil {
			_ = copierx.Copy(&cond.Page, req.Page)
		}
		if req.Category != nil {
			if *req.Category > int32(^uint16(0)>>1) {
				return nil, biz.ErrIllegalParameter(ctx, "category")
			}
			v := int16(*req.Category)
			cond.Category = &v
		}
		cond.Resource = req.Resource
	}

	list, err := s.whitelist.Find(ctx, cond)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindWhitelistReply{
		Page: &params.Page{},
		List: make([]*v1.Whitelist, 0, len(list)),
	}
	_ = copierx.Copy(&rp.Page, &cond.Page)
	for i := range list {
		rp.List = append(rp.List, whitelistToPB(&list[i]))
	}
	return rp, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) UpdateWhitelist(ctx context.Context, req *v1.UpdateWhitelistRequest) (*emptypb.Empty, error) {
	if s.whitelist == nil {
		return nil, biz.ErrInternal(ctx, "whitelist usecase not configured")
	}
	if req == nil || req.Id == 0 {
		return nil, biz.ErrIllegalParameter(ctx, "id")
	}

	w := &biz.UpdateWhitelist{
		ID:       req.Id,
		Resource: req.Resource,
	}
	if req.Category != nil {
		if *req.Category > int32(^uint16(0)>>1) {
			return nil, biz.ErrIllegalParameter(ctx, "category")
		}
		v := int16(*req.Category)
		w.Category = &v
	}

	if err := s.whitelist.Update(ctx, w); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) DeleteWhitelist(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	if s.whitelist == nil {
		return nil, biz.ErrInternal(ctx, "whitelist usecase not configured")
	}
	if req == nil || len(req.Ids) == 0 {
		return &emptypb.Empty{}, nil
	}
	if err := s.whitelist.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func whitelistToPB(w *biz.Whitelist) *v1.Whitelist {
	if w == nil {
		return &v1.Whitelist{}
	}
	rp := &v1.Whitelist{
		Id:       w.ID,
		Category: int32(w.Category),
		Resource: w.Resource,
	}
	return rp
}
