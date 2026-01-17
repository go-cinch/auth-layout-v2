package biz

import (
	"context"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewAuthUseCase,
	NewUserUseCase,
	NewPermissionUseCase,
	NewRoleUseCase,
	{{- if .Computed.enable_action_final }}
	NewActionUseCase,
	{{- end }}
	{{- if .Computed.enable_user_group_final }}
	NewUserGroupUseCase,
	{{- end }}
	{{- if .Computed.enable_whitelist_final }}
	NewWhitelistUseCase,
	{{- end }}
)

type Transaction interface {
	Tx(ctx context.Context, handler func(context.Context) error) error
}

type Cache interface {
	// Cache is get redis instance
	Cache() redis.UniversalClient
	// WithPrefix will add cache key prefix
	WithPrefix(prefix string) Cache
	// WithRefresh get data from db skip cache and refresh cache
	WithRefresh() Cache
	// Get is get cache data by key from redis, do write handler if cache is empty
	Get(ctx context.Context, action string, write func(context.Context) (string, error)) (string, error)
	// Set is set data to redis
	Set(ctx context.Context, action, data string, short bool)
	// Del delete key
	Del(ctx context.Context, action string)
	// SetWithExpiration is set data to redis with custom expiration
	SetWithExpiration(ctx context.Context, action, data string, seconds int64)
	// Flush is clean association cache if handler err=nil
	Flush(ctx context.Context, handler func(context.Context) error) error
	// FlushByPrefix clean cache by prefix, without prefix equals flush all by default cache prefix
	FlushByPrefix(ctx context.Context, prefix ...string) (err error)
}

