package biz

import (
	"context"
	"strings"

	"{{ .Computed.common_module_final }}/page/v2"
	"{{ .Computed.common_module_final }}/utils"

	"{{ .Computed.module_name_final }}/internal/conf"
)

// Action defines a fine-grained permission rule for a resource/menu/button.
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

type FindActionCache struct {
	Page page.Page `json:"page"`
	List []*Action `json:"list"`
}

type ActionRepo interface {
	Create(ctx context.Context, item *Action) error
	Update(ctx context.Context, item *Action) error
	Delete(ctx context.Context, ids []uint64) error

	Find(ctx context.Context, p *page.Page, filter *Action) (rp []*Action, total int64, err error)
	GetByIDs(ctx context.Context, ids []uint64) (rp []*Action, err error)
	GetTree(ctx context.Context) (rp []*Action, err error)

	// FindByCode returns actions by comma-separated codes.
	FindByCode(ctx context.Context, codes string) []Action
	// CodeExists validates that all comma-separated codes exist.
	CodeExists(ctx context.Context, codes string) error
}

type ActionUseCase struct {
	c     *conf.Bootstrap
	repo  ActionRepo
	tx    Transaction
	cache Cache
}

func NewActionUseCase(c *conf.Bootstrap, repo ActionRepo, tx Transaction, cache Cache) *ActionUseCase {
	return &ActionUseCase{
		c:    c,
		repo: repo,
		tx:   tx,
		cache: cache.WithPrefix(strings.Join([]string{
			c.Name, "action",
		}, "_")),
	}
}

func (uc *ActionUseCase) Create(ctx context.Context, item *Action) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) error {
			return uc.repo.Create(ctx, item)
		})
	})
}

func (uc *ActionUseCase) Find(ctx context.Context, p *page.Page, filter *Action) (rp []*Action, total int64, err error) {
	// Ensure stable cache keys by excluding Total (it is output, not input).
	var hashPage page.Page
	if p != nil {
		hashPage = page.Page{Num: p.Num, Size: p.Size, Disable: p.Disable}
	}
	keyObj := struct {
		Page   page.Page `json:"page"`
		Filter *Action   `json:"filter"`
	}{
		Page:   hashPage,
		Filter: filter,
	}

	action := strings.Join([]string{"find", utils.StructMd5(keyObj)}, "_")
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.find(ctx, action, p, filter)
	})
	if err != nil {
		return nil, 0, err
	}
	var cache FindActionCache
	utils.JSON2Struct(&cache, str)
	if p != nil {
		*p = cache.Page
	}
	return cache.List, cache.Page.Total, nil
}

func (uc *ActionUseCase) find(ctx context.Context, action string, p *page.Page, filter *Action) (res string, err error) {
	list, total, err := uc.repo.Find(ctx, p, filter)
	if err != nil {
		return "", err
	}
	var cache FindActionCache
	if p != nil {
		cache.Page = *p
	} else {
		cache.Page.Total = total
	}
	cache.List = list
	res = utils.Struct2JSON(cache)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return res, nil
}

func (uc *ActionUseCase) Update(ctx context.Context, item *Action) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			return uc.repo.Update(ctx, item)
		})
	})
}

func (uc *ActionUseCase) Delete(ctx context.Context, ids ...uint64) error {
	return uc.tx.Tx(ctx, func(ctx context.Context) error {
		return uc.cache.Flush(ctx, func(ctx context.Context) (err error) {
			return uc.repo.Delete(ctx, ids)
		})
	})
}

func (uc *ActionUseCase) GetByIDs(ctx context.Context, ids []uint64) (rp []*Action, err error) {
	// No cache: this is typically used internally in bulk flows and is cheap enough.
	return uc.repo.GetByIDs(ctx, ids)
}

func (uc *ActionUseCase) GetTree(ctx context.Context) (rp []*Action, err error) {
	// Tree is derived from actions; it is safe to cache and is flushed on any mutation.
	action := "tree"
	str, err := uc.cache.Get(ctx, action, func(ctx context.Context) (string, error) {
		return uc.getTree(ctx, action)
	})
	if err != nil {
		return nil, err
	}
	utils.JSON2Struct(&rp, str)
	return rp, nil
}

func (uc *ActionUseCase) getTree(ctx context.Context, action string) (res string, err error) {
	list, err := uc.repo.GetTree(ctx)
	if err != nil {
		return "", err
	}
	res = utils.Struct2JSON(list)
	uc.cache.Set(ctx, action, res, len(list) == 0)
	return res, nil
}

