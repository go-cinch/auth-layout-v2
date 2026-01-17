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

func (s *{{ .Computed.service_name_capitalized }}Service) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*emptypb.Empty, error) {
	if s.role == nil {
		return nil, biz.ErrInternal(ctx, "role usecase not configured")
	}
	if req == nil {
		return nil, biz.ErrIllegalParameter(ctx, "role")
	}

	r := &biz.Role{
		Name: req.Name,
		Word: req.Word,
	}
	if req.Action != nil {
		r.Action = *req.Action
	}
	if err := s.role.Create(ctx, r); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) FindRole(ctx context.Context, req *v1.FindRoleRequest) (*v1.FindRoleReply, error) {
	if s.role == nil {
		return nil, biz.ErrInternal(ctx, "role usecase not configured")
	}
	cond := &biz.FindRole{}
	if req != nil {
		if req.Page != nil {
			_ = copierx.Copy(&cond.Page, req.Page)
		}
		cond.Name = req.Name
		cond.Word = req.Word
		cond.Action = req.Action
	}

	list, err := s.role.Find(ctx, cond)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindRoleReply{
		Page: &params.Page{},
		List: make([]*v1.Role, 0, len(list)),
	}
	_ = copierx.Copy(&rp.Page, &cond.Page)
	for i := range list {
		rp.List = append(rp.List, roleToPB(&list[i]))
	}
	return rp, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*emptypb.Empty, error) {
	if s.role == nil {
		return nil, biz.ErrInternal(ctx, "role usecase not configured")
	}
	if req == nil || req.Id == 0 {
		return nil, biz.ErrIllegalParameter(ctx, "id")
	}

	r := &biz.UpdateRole{
		ID:     req.Id,
		Name:   req.Name,
		Word:   req.Word,
		Action: req.Action,
	}
	if err := s.role.Update(ctx, r); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) DeleteRole(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	if s.role == nil {
		return nil, biz.ErrInternal(ctx, "role usecase not configured")
	}
	if req == nil || len(req.Ids) == 0 {
		return &emptypb.Empty{}, nil
	}
	if err := s.role.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}
