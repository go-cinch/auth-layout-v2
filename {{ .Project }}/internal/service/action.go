package service

import (
	"context"
	"strings"

	"{{ .Computed.common_module_final }}/copierx"
	"{{ .Computed.common_module_final }}/page/v2"
	params "{{ .Computed.common_module_final }}/proto/params"
	"{{ .Computed.common_module_final }}/utils"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "{{ .Computed.module_name_final }}/api/{{ .Computed.service_name_final }}"
	"{{ .Computed.module_name_final }}/internal/biz"
)

func (s *{{ .Computed.service_name_capitalized }}Service) CreateAction(ctx context.Context, req *v1.CreateActionRequest) (*emptypb.Empty, error) {
	if s.action == nil {
		return nil, biz.ErrInternal(ctx, "action usecase not configured")
	}
	if req == nil {
		return nil, biz.ErrIllegalParameter(ctx, "action")
	}

	a := &biz.Action{}
	if name := strings.TrimSpace(req.Name); name != "" {
		a.Name = &name
	}
	if word := strings.TrimSpace(req.Word); word != "" {
		a.Word = &word
	}
	if req.Resource != nil {
		if res := strings.TrimSpace(*req.Resource); res != "" {
			a.Resource = &res
		}
	}
	if req.Menu != nil {
		if menu := strings.TrimSpace(*req.Menu); menu != "" {
			a.Menu = &menu
		}
	}
	if req.Btn != nil {
		if btn := strings.TrimSpace(*req.Btn); btn != "" {
			a.Btn = &btn
		}
	}

	if err := s.action.Create(ctx, a); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) FindAction(ctx context.Context, req *v1.FindActionRequest) (*v1.FindActionReply, error) {
	if s.action == nil {
		return nil, biz.ErrInternal(ctx, "action usecase not configured")
	}

	var p page.Page
	filter := &biz.Action{}
	if req != nil {
		if req.Page != nil {
			_ = copierx.Copy(&p, req.Page)
		}
		filter.Code = req.Code
		filter.Name = req.Name
		filter.Word = req.Word
		filter.Resource = req.Resource
	}

	list, _, err := s.action.Find(ctx, &p, filter)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindActionReply{
		Page: &params.Page{},
		List: make([]*v1.Action, 0, len(list)),
	}
	_ = copierx.Copy(&rp.Page, &p)
	for i := range list {
		rp.List = append(rp.List, actionToPB(list[i]))
	}
	return rp, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) UpdateAction(ctx context.Context, req *v1.UpdateActionRequest) (*emptypb.Empty, error) {
	if s.action == nil {
		return nil, biz.ErrInternal(ctx, "action usecase not configured")
	}
	if req == nil || req.Id == 0 {
		return nil, biz.ErrIllegalParameter(ctx, "id")
	}

	a := &biz.Action{
		ID: req.Id,
	}
	if req.Name != nil {
		v := strings.TrimSpace(*req.Name)
		a.Name = &v
	}
	if req.Word != nil {
		v := strings.TrimSpace(*req.Word)
		a.Word = &v
	}
	if req.Resource != nil {
		v := strings.TrimSpace(*req.Resource)
		a.Resource = &v
	}
	if req.Menu != nil {
		v := strings.TrimSpace(*req.Menu)
		a.Menu = &v
	}
	if req.Btn != nil {
		v := strings.TrimSpace(*req.Btn)
		a.Btn = &v
	}

	if err := s.action.Update(ctx, a); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) DeleteAction(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	if s.action == nil {
		return nil, biz.ErrInternal(ctx, "action usecase not configured")
	}
	if req == nil || len(req.Ids) == 0 {
		return &emptypb.Empty{}, nil
	}
	if err := s.action.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}
