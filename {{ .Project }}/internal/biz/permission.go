package biz

import (
	"context"
	"strings"

	"{{ .Computed.module_name_final }}/internal/conf"
)

{{- if not .Computed.enable_action_final }}

// Action defines an RBAC permission rule. It is part of the core RBAC feature set.
// When the optional Action module is enabled, the full Action use case/type lives
// in `internal/biz/action.go` instead.
type Action struct {
	ID       uint64    `json:"id,string"`
	Code     *string   `json:"code,omitempty"`
	Name     *string   `json:"name,omitempty"`
	Word     *string   `json:"word,omitempty"`
	Resource *string   `json:"resource,omitempty"`
	Menu     *string   `json:"menu,omitempty"`
	Btn      *string   `json:"btn,omitempty"`
	Children []*Action `json:"children,omitempty"`
}

{{- end }}

type Permission struct {
	Resources []string `json:"resources"`
	Menus     []string `json:"menus"`
	Btns      []string `json:"btns"`
}

type CheckPermission struct {
	UserID   uint64 `json:"userId,string"`
	Resource string `json:"resource"`
	Method   string `json:"method"`
	URI      string `json:"uri"`
}

type PermissionRepo interface {
	// GetUserPermissions returns all actions a user can access.
	GetUserPermissions(ctx context.Context, userID uint64) ([]*Action, error)
	// CheckPermission checks if a user is allowed to access the resource.
	// If method is empty, resource is treated as a gRPC operation (exact match).
	// If method is not empty, resource is treated as an HTTP URI and matched against action rules.
	CheckPermission(ctx context.Context, userID uint64, resource, method string) (bool, error)
}

type PermissionUseCase struct {
	c    *conf.Bootstrap
	repo PermissionRepo
}

func NewPermissionUseCase(c *conf.Bootstrap, repo PermissionRepo) *PermissionUseCase {
	return &PermissionUseCase{
		c:    c,
		repo: repo,
	}
}

// Check validates a request against the caller's permissions.
// When Method is set, URI is used as the resource to check.
func (uc *PermissionUseCase) Check(ctx context.Context, req *CheckPermission) (ok bool, err error) {
	if req == nil {
		return false, ErrIllegalParameter(ctx, "permission")
	}
	if req.UserID == 0 {
		return false, ErrIllegalParameter(ctx, "userID")
	}

	resource := strings.TrimSpace(req.Resource)
	method := strings.TrimSpace(req.Method)
	if method != "" {
		resource = strings.TrimSpace(req.URI)
	}

	return uc.repo.CheckPermission(ctx, req.UserID, resource, method)
}

func (uc *PermissionUseCase) CheckPermission(ctx context.Context, userID uint64, resource, method string) (bool, error) {
	return uc.repo.CheckPermission(ctx, userID, resource, method)
}

func (uc *PermissionUseCase) GetUserPermissions(ctx context.Context, userID uint64) ([]*Action, error) {
	if userID == 0 {
		return nil, ErrIllegalParameter(ctx, "userID")
	}
	return uc.repo.GetUserPermissions(ctx, userID)
}

func (uc *PermissionUseCase) GetByUserID(ctx context.Context, userID uint64) (*Permission, error) {
	if userID == 0 {
		return nil, ErrIllegalParameter(ctx, "userID")
	}

	actions, err := uc.repo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	rp := &Permission{
		Resources: make([]string, 0),
		Menus:     make([]string, 0),
		Btns:      make([]string, 0),
	}

	seenRes := make(map[string]struct{})
	seenMenu := make(map[string]struct{})
	seenBtn := make(map[string]struct{})

	addLines := func(dst *[]string, seen map[string]struct{}, s *string) {
		if s == nil {
			return
		}
		for _, line := range strings.Split(*s, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if _, ok := seen[line]; ok {
				continue
			}
			seen[line] = struct{}{}
			*dst = append(*dst, line)
		}
	}

	for _, a := range actions {
		if a == nil {
			continue
		}
		addLines(&rp.Resources, seenRes, a.Resource)
		addLines(&rp.Menus, seenMenu, a.Menu)
		addLines(&rp.Btns, seenBtn, a.Btn)
	}

	return rp, nil
}

