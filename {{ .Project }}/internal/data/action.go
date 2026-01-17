package data

import (
	"context"
	"errors"
	"sort"
	"strings"

	"{{ .Computed.common_module_final }}/copierx"
	"{{ .Computed.common_module_final }}/id"
	"{{ .Computed.common_module_final }}/log"
	"{{ .Computed.common_module_final }}/page/v2"
	"{{ .Computed.common_module_final }}/utils"
	"gorm.io/gorm"

	"{{ .Computed.module_name_final }}/internal/biz"
	"{{ .Computed.module_name_final }}/internal/data/model"
)

type actionRepo struct {
	data *Data
}

func NewActionRepo(data *Data) biz.ActionRepo {
	return &actionRepo{
		data: data,
	}
}

func (ro actionRepo) Create(ctx context.Context, item *biz.Action) (err error) {
	if item == nil {
		return biz.ErrIllegalParameter(ctx, "action")
	}

	db := gorm.G[model.Action](ro.data.DB(ctx))

	// Normalize and check if word exists (word is nullable).
	if item.Word != nil {
		word := strings.TrimSpace(*item.Word)
		item.Word = &word
		count, err := db.Where("word = ?", word).Count(ctx, "*")
		if err != nil {
			log.WithError(err).Error("check action word exists failed")
			return err
		}
		if count > 0 {
			return biz.ErrDuplicateField(ctx, "word", word)
		}
	}

	if item.ID == 0 {
		item.ID = ro.data.ID(ctx)
	}

	var m model.Action
	_ = copierx.Copy(&m, item)
	m.ID = item.ID
	if m.Word != nil {
		word := strings.TrimSpace(*m.Word)
		m.Word = &word
	}

	// Always generate code from ID for uniqueness/consistency.
	code := id.NewCode(m.ID)
	m.Code = &code

	// Default resource is "*" when empty.
	if m.Resource == nil || strings.TrimSpace(*m.Resource) == "" {
		res := "*"
		m.Resource = &res
	}

	err = db.Create(ctx, &m)
	if err != nil {
		log.WithError(err).Error("create action failed")
		return err
	}

	// Best-effort: propagate generated code back to caller.
	item.Code = &code
	return nil
}

func (ro actionRepo) Update(ctx context.Context, item *biz.Action) (err error) {
	if item == nil || item.ID == 0 {
		return biz.ErrIllegalParameter(ctx, "id")
	}

	if item.Word != nil {
		word := strings.TrimSpace(*item.Word)
		item.Word = &word
	}

	db := gorm.G[model.Action](ro.data.DB(ctx))

	m, err := db.Where("id = ?", item.ID).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return biz.ErrRecordNotFound(ctx)
	}
	if err != nil {
		log.WithError(err).Error("get action failed")
		return err
	}

	change := make(map[string]interface{})
	utils.CompareDiff(m, item, &change)
	delete(change, "id")
	delete(change, "code") // code is derived from id; don't allow updating it here

	if len(change) == 0 {
		return biz.ErrDataNotChange(ctx)
	}

	// Check word uniqueness if word is being updated.
	if item.Word != nil {
		word := strings.TrimSpace(*item.Word)
		oldWord := ""
		if m.Word != nil {
			oldWord = strings.TrimSpace(*m.Word)
		}
		if word != oldWord {
			count, err := db.Where("word = ? AND id <> ?", word, item.ID).Count(ctx, "*")
			if err != nil {
				log.WithError(err).Error("check action word uniqueness failed")
				return err
			}
			if count > 0 {
				return biz.ErrDuplicateField(ctx, "word", word)
			}
		}
	}

	err = ro.data.DB(ctx).
		Model(&model.Action{}).
		Where("id = ?", item.ID).
		Updates(change).
		Error
	if err != nil {
		log.WithError(err).Error("update action failed")
	}
	return err
}

func (ro actionRepo) Delete(ctx context.Context, ids []uint64) (err error) {
	if len(ids) == 0 {
		return nil
	}

	db := gorm.G[model.Action](ro.data.DB(ctx))
	_, err = db.Where("id IN ?", ids).Delete(ctx)
	if err != nil {
		log.WithError(err).Error("delete action failed")
	}
	return err
}

func (ro actionRepo) Find(ctx context.Context, p *page.Page, filter *biz.Action) (rp []*biz.Action, total int64, err error) {
	db := gorm.G[model.Action](ro.data.DB(ctx))
	q := db.Where("1 = 1")

	if filter != nil {
		if filter.ID != 0 {
			q = q.Where("id = ?", filter.ID)
		}
		if filter.Code != nil && strings.TrimSpace(*filter.Code) != "" {
			q = q.Where("code LIKE ?", "%"+strings.TrimSpace(*filter.Code)+"%")
		}
		if filter.Name != nil && strings.TrimSpace(*filter.Name) != "" {
			q = q.Where("name LIKE ?", "%"+strings.TrimSpace(*filter.Name)+"%")
		}
		if filter.Word != nil && strings.TrimSpace(*filter.Word) != "" {
			q = q.Where("word LIKE ?", "%"+strings.TrimSpace(*filter.Word)+"%")
		}
		if filter.Resource != nil && strings.TrimSpace(*filter.Resource) != "" {
			q = q.Where("resource LIKE ?", "%"+strings.TrimSpace(*filter.Resource)+"%")
		}
	}

	total, err = q.Count(ctx, "*")
	if err != nil {
		log.WithError(err).Error("count action failed")
		return nil, 0, err
	}
	if p != nil {
		p.Total = total
		if !p.Disable {
			limit, offset := p.Limit()
			q = q.Limit(limit).Offset(offset)
		}
	}

	list, err := q.Order("id DESC").Find(ctx)
	if err != nil {
		log.WithError(err).Error("find action failed")
		return nil, total, err
	}

	rp = make([]*biz.Action, 0, len(list))
	for i := range list {
		it := new(biz.Action)
		_ = copierx.Copy(it, list[i])
		rp = append(rp, it)
	}
	return rp, total, nil
}

func (ro actionRepo) GetByIDs(ctx context.Context, ids []uint64) (rp []*biz.Action, err error) {
	if len(ids) == 0 {
		return []*biz.Action{}, nil
	}

	db := gorm.G[model.Action](ro.data.DB(ctx))
	list, err := db.Where("id IN ?", ids).Order("id ASC").Find(ctx)
	if err != nil {
		log.WithError(err).Error("get actions by ids failed")
		return nil, err
	}

	rp = make([]*biz.Action, 0, len(list))
	for i := range list {
		it := new(biz.Action)
		_ = copierx.Copy(it, list[i])
		rp = append(rp, it)
	}
	return rp, nil
}

func (ro actionRepo) GetTree(ctx context.Context) (rp []*biz.Action, err error) {
	// Build a tree from all menu entries across actions.
	// Each line in `menu` is treated as a menu path; separators supported: "/", ".", ":".
	db := gorm.G[model.Action](ro.data.DB(ctx))
	list, err := db.Select("menu").Find(ctx)
	if err != nil {
		log.WithError(err).Error("get action menus failed")
		return nil, err
	}

	paths := make(map[string]struct{})
	for i := range list {
		if list[i].Menu == nil {
			continue
		}
		for _, line := range strings.Split(*list[i].Menu, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			paths[line] = struct{}{}
		}
	}

	type node struct {
		act      *biz.Action
		children map[string]*node
	}

	root := &node{children: map[string]*node{}}
	for path := range paths {
		parts := actionSplitMenuPath(path)
		if len(parts) == 0 {
			continue
		}
		cur := root
		full := ""
		for _, part := range parts {
			if part == "" {
				continue
			}
			if full == "" {
				full = part
			} else {
				full = full + "/" + part
			}
			child := cur.children[part]
			if child == nil {
				name := part
				word := full
				child = &node{
					act: &biz.Action{
						Name:     &name,
						Word:     &word,
						Children: make([]*biz.Action, 0),
					},
					children: map[string]*node{},
				}
				cur.children[part] = child
			}
			cur = child
		}
	}

	var build func(n *node) []*biz.Action
	build = func(n *node) []*biz.Action {
		if len(n.children) == 0 {
			return nil
		}
		keys := make([]string, 0, len(n.children))
		for k := range n.children {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		out := make([]*biz.Action, 0, len(keys))
		for _, k := range keys {
			child := n.children[k]
			child.act.Children = build(child)
			out = append(out, child.act)
		}
		return out
	}

	rp = build(root)
	return rp, nil
}

func (ro actionRepo) FindByCode(ctx context.Context, codes string) (rp []biz.Action) {
	rp = make([]biz.Action, 0)
	codes = strings.TrimSpace(codes)
	if codes == "" {
		return rp
	}

	arr := actionSplitComma(codes)
	if len(arr) == 0 {
		return rp
	}

	db := gorm.G[model.Action](ro.data.DB(ctx))
	list, err := db.Where("code IN ?", arr).Find(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("find action by code failed")
		return rp
	}

	byCode := make(map[string]model.Action, len(list))
	for _, m := range list {
		if m.Code == nil {
			continue
		}
		c := strings.TrimSpace(*m.Code)
		if c == "" {
			continue
		}
		byCode[c] = m
	}

	for _, code := range arr {
		m, ok := byCode[code]
		if !ok {
			continue
		}
		var a biz.Action
		_ = copierx.Copy(&a, m)
		rp = append(rp, a)
	}
	return rp
}

func (ro actionRepo) CodeExists(ctx context.Context, codes string) (err error) {
	codes = strings.TrimSpace(codes)
	if codes == "" {
		return nil
	}

	arr := actionSplitComma(codes)
	if len(arr) == 0 {
		return nil
	}

	db := gorm.G[model.Action](ro.data.DB(ctx))
	list, qErr := db.Select("code").Where("code IN ?", arr).Find(ctx)
	if qErr != nil {
		log.WithContext(ctx).WithError(qErr).Error("check action code exists failed")
		return qErr
	}

	seen := make(map[string]struct{}, len(list))
	for _, m := range list {
		if m.Code == nil {
			continue
		}
		seen[strings.TrimSpace(*m.Code)] = struct{}{}
	}

	for _, code := range arr {
		if _, ok := seen[code]; ok {
			continue
		}
		err = biz.ErrRecordNotFound(ctx)
		log.WithContext(ctx).WithError(err).Error("invalid code: %s", codes)
		return err
	}
	return nil
}

func actionSplitMenuPath(path string) []string {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}

	// Prefer "/" as it's most common for menu paths.
	if strings.Contains(path, "/") {
		path = strings.Trim(path, "/")
	}

	var parts []string
	switch {
	case strings.Contains(path, "/"):
		parts = strings.Split(path, "/")
	case strings.Contains(path, "."):
		parts = strings.Split(path, ".")
	case strings.Contains(path, ":"):
		parts = strings.Split(path, ":")
	default:
		parts = []string{path}
	}

	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

func actionSplitComma(s string) []string {
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, it := range raw {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		if _, ok := seen[it]; ok {
			continue
		}
		seen[it] = struct{}{}
		out = append(out, it)
	}
	return out
}

