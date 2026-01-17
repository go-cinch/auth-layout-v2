package biz

import (
	"context"
	"strings"

	"{{ .Computed.common_module_final }}/page/v2"
	"{{ .Computed.common_module_final }}/utils"

	"{{ .Computed.module_name_final }}/internal/conf"
)

// UserGroup is a group of users with associated permission action codes.
// It enables group-based user management (many users can belong to many groups).
type UserGroup struct {
	Id      uint64   `json:"id,string"`
	Users   []User   `json:"users"`
	Name    string   `json:"name"`
	Word    string   `json:"word"`
	Action  string   `json:"action"`
	Actions []Action `json:"actions"`
}

type FindUserGroup struct {
	Page   page.Page `json:"page"`
	Name   *string   `json:"name"`
	Word   *string   `json:"word"`
	Action *string   `json:"action"`
}

type FindUserGroupCache struct {
	Page page.Page   `json:"page"`
	List []UserGroup `json:"list"`
}

type UpdateUserGroup struct {
	Id     uint64  `json:"id,string"`
	Name   *string `json:"name,omitempty"`
	Word   *string `json:"word,omitempty"`
	Action *string `json:"action,omitempty"`
	// Users is a comma-separated list of user ids, used for bulk replace.
	Users *string `json:"users,omitempty"`
}

type UserGroupRepo interface {
	Create(ctx context.Context, item *UserGroup) error
	Find(ctx context.Context, condition *FindUserGroup) []UserGroup
	Update(ctx context.Context, item *UpdateUserGroup) error
	Delete(ctx context.Context, ids ...uint64) error
}

type UserGroupUseCase struct {
	c     *conf.Bootstrap
	repo  UserGroupRepo
	tx    Transaction
	cache Cache
}

func NewUserGroupUseCase(c *conf.Bootstrap, repo UserGroupRepo, tx Transaction, cache Cache) *UserGroupUseCase {
	return &UserGroupUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		// Keep prefix short; group cache should be flushed along with user-related changes.
		cache: cache.WithPrefix(strings.Join([]string{c.Name, "group"}, "_")),
	}
}

func (uc *UserGroupUseCase) Create(ctx context.Context, item *UserGroup) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *UserGroupUseCase) Find(ctx context.Context, condition *FindUserGroup) (rp []UserGroup, err error) {
	action := strings.Join([]string{"find", utils.StructMd5(condition)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, condition)
	})
	if err != nil {
		return nil, err
	}
	var cache FindUserGroupCache
	utils.JSON2Struct(&cache, str)
	condition.Page = cache.Page
	return cache.List, nil
}

func (uc *UserGroupUseCase) find(ctx context.Context, action string, condition *FindUserGroup) (res string, err error) {
	list := uc.repo.Find(ctx, condition)
	var cache FindUserGroupCache
	cache.List = list
	cache.Page = condition.Page
	res = utils.Struct2JSON(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return res, nil
}

func (uc *UserGroupUseCase) Update(ctx context.Context, item *UpdateUserGroup) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Update(ctx, item)
		})
	})
}

func (uc *UserGroupUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Delete(ctx, ids...)
		})
	})
}
