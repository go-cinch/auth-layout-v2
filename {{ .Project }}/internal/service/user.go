package service

import (
	"context"
	{{- if .Computed.enable_user_lock_final }}
	"strconv"
	"strings"
	{{- end }}

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
	if s.user == nil {
		return nil, biz.ErrInternal(ctx, "user usecase not configured")
	}
	cond := &biz.FindUser{}
	if req != nil {
		if req.Page != nil {
			_ = copierx.Copy(&cond.Page, req.Page)
		}
		cond.StartCreatedAt = req.StartCreatedAt
		cond.EndCreatedAt = req.EndCreatedAt
		cond.StartUpdatedAt = req.StartUpdatedAt
		cond.EndUpdatedAt = req.EndUpdatedAt
		cond.Username = req.Username
		cond.Code = req.Code
		cond.Platform = req.Platform
		cond.Locked = req.Locked
	}

	list, err := s.user.Find(ctx, cond)
	if err != nil {
		return nil, err
	}

	rp := &v1.FindUserReply{
		Page: &params.Page{},
		List: make([]*v1.User, 0, len(list)),
	}
	_ = copierx.Copy(&rp.Page, &cond.Page)
	for i := range list {
		u := list[i]
		rp.List = append(rp.List, userToPB(&u))
	}
	return rp, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*emptypb.Empty, error) {
	if s.user == nil {
		return nil, biz.ErrInternal(ctx, "user usecase not configured")
	}
	if req == nil || req.Id == 0 {
		return nil, biz.ErrIllegalParameter(ctx, "id")
	}

	u := &biz.UpdateUser{
		Id:       req.Id,
		Username: req.Username,
		Password: req.Password,
		Platform: req.Platform,
		Locked:   req.Locked,
		Action:   req.Action,
		RoleId:   req.RoleId,
	}
	{{- if .Computed.enable_user_lock_final }}
	if req.LockExpire != nil {
		u.LockExpire = req.LockExpire
	}
	{{- end }}

	if err := s.user.Update(ctx, u); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func (s *{{ .Computed.service_name_capitalized }}Service) DeleteUser(ctx context.Context, req *params.IdsRequest) (*emptypb.Empty, error) {
	if s.user == nil {
		return nil, biz.ErrInternal(ctx, "user usecase not configured")
	}
	if req == nil || len(req.Ids) == 0 {
		return &emptypb.Empty{}, nil
	}
	if err := s.user.Delete(ctx, utils.Str2Uint64Arr(req.Ids)...); err != nil {
		return nil, err
	}
	s.flushCache(ctx)
	return &emptypb.Empty{}, nil
}

func userToPB(u *biz.User) *v1.User {
	if u == nil {
		return &v1.User{}
	}

	rp := &v1.User{
		Id:       u.Id,
		Username: u.Username,
		Code:     u.Code,
		Platform: u.Platform,
		Locked:   u.Locked,
		RoleId:   u.RoleId,
	}

	// Time fields are rendered as strings in the API.
	rp.CreatedAt = u.CreatedAt.ToDateTimeString()
	rp.UpdatedAt = u.UpdatedAt.ToDateTimeString()
	rp.LastLogin = u.LastLogin.ToDateTimeString()

	{{- if .Computed.enable_user_lock_final }}
	rp.LockExpire = u.LockExpire
	rp.LockMsg = u.LockMsg
	{{- end }}

	// Nested role/actions are not reliably handled by generic copiers.
	rp.Role = roleToPB(&u.Role)
	rp.Actions = actionsToPB(u.Actions)
	return rp
}

func roleToPB(r *biz.Role) *v1.Role {
	if r == nil {
		return &v1.Role{}
	}
	rp := &v1.Role{
		Id:     r.ID,
		Name:   r.Name,
		Word:   r.Word,
		Action: r.Action,
	}
	rp.Actions = actionsToPB(r.Actions)
	return rp
}

func actionsToPB(list []biz.Action) []*v1.Action {
	if len(list) == 0 {
		return []*v1.Action{}
	}
	rp := make([]*v1.Action, 0, len(list))
	for i := range list {
		a := list[i]
		rp = append(rp, actionToPB(&a))
	}
	return rp
}

func actionToPB(a *biz.Action) *v1.Action {
	if a == nil {
		return &v1.Action{}
	}
	return &v1.Action{
		Id:       a.ID,
		Code:     derefString(a.Code),
		Name:     derefString(a.Name),
		Word:     derefString(a.Word),
		Resource: derefString(a.Resource),
		Menu:     derefString(a.Menu),
		Btn:      derefString(a.Btn),
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

{{- if .Computed.enable_user_lock_final }}
func parseLockExpireTime(s string) (int64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, true
	}
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return v, true
	}

	// Accept duration strings like "5m", "24h".
	if dt := carbon.Now().AddDuration(s); dt.Error == nil {
		return dt.StdTime().Unix(), true
	}

	// Accept absolute datetime strings.
	parsed := carbon.Parse(s)
	if parsed.Error == nil {
		return parsed.StdTime().Unix(), true
	}
	return 0, false
}
{{- end }}
