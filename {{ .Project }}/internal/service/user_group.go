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

func (s *{{ .Computed.service_name_capitalized }}Service) CreateUserGroup(ctx context.Context, req *v1.CreateUserGroupRequest) (*emptypb.Empty, error) {
	if s.userGroup == nil {
		return nil, biz.ErrInternal(ctx, "user_group usecase not configured")
	}
	if req == nil {
		return nil, biz.ErrIllegalParameter(ctx, "user_group")
	}

	g := &biz.UserGroup{
		Name: req.Name,
		Word: req.Word,
	}
	if req.Action != nil {
		g.Action = *req.Action
	}
	if len(req.Users) > 0 {
		g.Users = make([]biz.User, 0, len(req.Users))
		for _, id := range req.Users {
			g.Users = append(g.Users, biz.User{Id: id})
		}
	}

	if err := s.userGroup.Create(ctx, g); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) FindUserGroup(ctx context.Context, req *v1.FindUserGroupRequest) (*v1.FindUserGroupReply, error) {
	if s.userGroup == nil {
		return nil, biz.ErrInternal(ctx, "user_group usecase not configured")
	}

	cond := &biz.FindUserGroup{}
	if req != nil {
		if req.Page != nil {
			_ = copierx.Copy(&cond.Page, req.Page)
		}
		cond.Name = req.Name
		cond.Word = req.Word
		cond.Action = req.Action
	}

	list, err := s.userGroup.Find(ctx, cond)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindUserGroupReply{
		Page: &params.Page{},
		List: make([]*v1.UserGroup, 0, len(list)),
	}
	_ = copierx.Copy(&rp.Page, &cond.Page)
	for i := range list {
		rp.List = append(rp.List, userGroupToPB(&list[i]))
	}
	return rp, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) UpdateUserGroup(ctx context.Context, req *v1.UpdateUserGroupRequest) (*emptypb.Empty, error) {
	if s.userGroup == nil {
		return nil, biz.ErrInternal(ctx, "user_group usecase not configured")
	}
	if req == nil || req.Id == 0 {
		return nil, biz.ErrIllegalParameter(ctx, "id")
	}

	g := &biz.UpdateUserGroup{
		Id:     req.Id,
		Name:   req.Name,
		Word:   req.Word,
		Action: req.Action,
		Users:  req.Users,
	}
	if err := s.userGroup.Update(ctx, g); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) DeleteUserGroup(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	if s.userGroup == nil {
		return nil, biz.ErrInternal(ctx, "user_group usecase not configured")
	}
	if req == nil || len(req.Ids) == 0 {
		return &emptypb.Empty{}, nil
	}
	if err := s.userGroup.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func userGroupToPB(g *biz.UserGroup) *v1.UserGroup {
	if g == nil {
		return &v1.UserGroup{}
	}
	rp := &v1.UserGroup{
		Id:      g.Id,
		Name:    g.Name,
		Word:    g.Word,
		Actions: actionsToPB(g.Actions),
		Users:   make([]*v1.User, 0, len(g.Users)),
	}
	for i := range g.Users {
		u := g.Users[i]
		rp.Users = append(rp.Users, userToPB(&u))
	}
	return rp
}
