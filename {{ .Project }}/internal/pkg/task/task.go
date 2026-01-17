package task

{{- if and .Computed.enable_redis_final .Computed.enable_task_final }}

import (
	"context"
	"strings"
	"time"

	"{{ .Computed.common_module_final }}/log"
	"{{ .Computed.common_module_final }}/utils"
	"{{ .Computed.common_module_final }}/worker"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"

	"{{ .Computed.module_name_final }}/internal/biz"
	"{{ .Computed.module_name_final }}/internal/conf"
)

// ProviderSet is task providers.
var ProviderSet = wire.NewSet(New)

// New initializes the task worker from config.
func New(c *conf.Bootstrap, user *biz.UserUseCase{{ if .Computed.enable_hotspot_final }}, hotspot biz.HotspotRepo{{ end }}) (w *worker.Worker, err error) {
	if c == nil || c.Redis == nil || c.Redis.Dsn == "" {
		return nil, errors.New("redis config is required for tasks")
	}

	w = worker.New(
		worker.WithRedisURI(c.Redis.Dsn),
		worker.WithGroup(c.Name),
		worker.WithHandler(func(ctx context.Context, p worker.Payload) error {
			return process(task{
				ctx:     ctx,
				c:       c,
				payload: p,
				user:    user,
				{{- if .Computed.enable_hotspot_final }}
				hotspot: hotspot,
				{{- end }}
			})
		}),
	)
	if w.Error != nil {
		log.Error(w.Error)
		return nil, errors.New("initialize worker failed")
	}

	if c.Task != nil {
		for id, item := range c.Task.Cron {
			if item == nil {
				continue
			}
			err = w.Cron(
				context.Background(),
				worker.WithRunUUID(id),
				worker.WithRunGroup(item.Name),
				worker.WithRunExpr(item.Expr),
				worker.WithRunTimeout(int(item.Timeout)),
				worker.WithRunMaxRetry(int(item.Retry)),
			)
			if err != nil {
				log.Error(err)
				return nil, errors.New("initialize worker failed")
			}
		}
	}

	log.Info("initialize worker success")

	{{- if .Computed.enable_hotspot_final }}
	// When app restart, clear hotspot (best-effort).
	if c.Task != nil && c.Task.Group != nil && c.Task.Group.RefreshHotspotManual != "" {
		_ = w.Once(
			context.Background(),
			worker.WithRunUUID(strings.Join([]string{c.Task.Group.RefreshHotspotManual}, ".")),
			worker.WithRunGroup(c.Task.Group.RefreshHotspotManual),
			worker.WithRunIn(10*time.Second),
			worker.WithRunReplace(true),
		)
	}
	{{- end }}

	return w, nil
}

type task struct {
	ctx     context.Context
	c       *conf.Bootstrap
	payload worker.Payload
	user    *biz.UserUseCase
	{{- if .Computed.enable_hotspot_final }}
	hotspot biz.HotspotRepo
	{{- end }}
}

func process(t task) (err error) {
	tr := otel.Tracer("task")
	ctx, span := tr.Start(t.ctx, "Task")
	defer span.End()

	if t.c == nil || t.c.Task == nil || t.c.Task.Group == nil {
		log.WithContext(ctx).Warn("task config not loaded")
		return nil
	}

	// Use task group to match tasks instead of UID.
	switch t.payload.Group {
	case t.c.Task.Group.LoginFailed:
		var req biz.LoginTime
		utils.JSON2Struct(&req, t.payload.Payload)
		if t.user != nil {
			err = t.user.WrongPwd(ctx, &req)
		}
	case t.c.Task.Group.LoginLast:
		var req biz.LoginTime
		utils.JSON2Struct(&req, t.payload.Payload)
		if t.user != nil {
			err = t.user.LastLogin(ctx, req.Username)
		}
	{{- if .Computed.enable_hotspot_final }}
	case t.c.Task.Group.RefreshHotspot, t.c.Task.Group.RefreshHotspotManual:
		if t.hotspot != nil {
			err = t.hotspot.Refresh(ctx)
		}
	{{- end }}
	default:
		log.WithContext(ctx).Warn("unknown task group: %s", t.payload.Group)
	}
	return err
}

{{- else }}

import "github.com/google/wire"

// ProviderSet is empty when tasks are disabled.
var ProviderSet = wire.NewSet()

{{- end }}

