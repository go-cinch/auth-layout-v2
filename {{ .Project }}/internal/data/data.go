package data

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"{{ .Computed.common_module_final }}/id"
	"{{ .Computed.common_module_final }}/log"
	"{{ .Computed.common_module_final }}/utils"
	glog "{{ .Computed.common_module_final }}/plugins/gorm/log"
	"{{ .Computed.common_module_final }}/plugins/gorm/tenant/v2"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"{{ .Computed.module_name_final }}/internal/biz"
	"{{ .Computed.module_name_final }}/internal/conf"
	"{{ .Computed.module_name_final }}/internal/db"
)

// Data wraps all data sources used by the service.
type Data struct {
	Tenant    *tenant.Tenant
	sonyflake *id.Sonyflake
}

// NewData initializes the configured database connection via tenant v2.
func NewData(c *conf.Bootstrap) (*Data, func(), error) {
	gormTenant, err := NewDB(c)
	if err != nil {
		return nil, nil, err
	}

	sonyflake, err := NewSonyflake(c)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.Info("closing database connections")
	}

	return &Data{
		Tenant:    gormTenant,
		sonyflake: sonyflake,
	}, cleanup, nil
}

// NewDB initializes tenant-aware database connection using tenant v2 package.
// Supports both MySQL and PostgreSQL via internal driver detection.
func NewDB(c *conf.Bootstrap) (*tenant.Tenant, error) {
	if c == nil || c.Db == nil {
		err := errors.New("db config is required")
		log.WithError(err).Error("initialize db failed")
		return nil, err
	}

	dbConf := c.Db
	driver := strings.ToLower(strings.TrimSpace(dbConf.Driver))
	dsn := strings.TrimSpace(dbConf.Dsn)
	if driver == "" {
		err := errors.New("db driver is required")
		log.WithError(err).Error("initialize db failed")
		return nil, err
	}
	if dsn == "" {
		err := errors.New("db DSN is required")
		log.WithError(err).Error("initialize db failed")
		return nil, err
	}

	level := log.NewLevel(c.Log.Level)
	// Force to warn level when show sql is false.
	if level > log.WarnLevel && !c.Log.ShowSQL {
		level = log.WarnLevel
	}

	ops := []func(*tenant.Options){
		tenant.WithDriver(driver),
		tenant.WithDSN("", dsn), // Empty string for default tenant
		tenant.WithSQLFile(db.SQLFiles),
		tenant.WithSQLRoot(db.SQLRoot),
		tenant.WithSkipMigrate(!dbConf.Migrate), // Skip migration if Migrate is false
		tenant.WithConfig(&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			QueryFields: true,
			Logger: glog.New(
				glog.WithColorful(false),
				glog.WithSlow(200),
				glog.WithLevel(level),
			),
		}),
		tenant.WithMaxIdle(10),
		tenant.WithMaxOpen(100),
	}

	gormTenant, err := tenant.New(ops...)
	if err != nil {
		log.WithError(err).Error("create tenant failed")
		return nil, err
	}

	// Always call Migrate() to initialize the database connection.
	// WithSkipMigrate controls whether SQL migrations are actually executed.
	if err := gormTenant.Migrate(); err != nil {
		log.WithError(err).Error("migrate tenant failed")
		return nil, err
	}

	log.Info("initialize db success, driver: %s", driver)
	return gormTenant, nil
}

type contextTxKey struct{}

// Tx is transaction wrapper.
func (d *Data) Tx(ctx context.Context, handler func(ctx context.Context) error) error {
	return d.Tenant.DB(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return handler(ctx)
	})
}

// DB returns a tenant-aware GORM DB instance from context.
// If a transaction is present in the context, it returns the transaction DB.
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.Tenant.DB(ctx)
}

// NewTransaction creates a new Transaction from Data.
func NewTransaction(d *Data) biz.Transaction {
	return d
}

// ID generates a unique distributed ID using Sonyflake.
func (d *Data) ID(ctx context.Context) uint64 {
	return d.sonyflake.ID(ctx)
}

// NewSonyflake initializes the Sonyflake ID generator.
func NewSonyflake(c *conf.Bootstrap) (*id.Sonyflake, error) {
	machineID, _ := strconv.ParseUint(c.Server.MachineId, 10, 16)
	sf := id.NewSonyflake(
		id.WithSonyflakeMachineID(uint16(machineID)),
		id.WithSonyflakeStartTime(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
	)
	if sf.Error != nil {
		log.WithError(sf.Error).Error("initialize sonyflake failed")
		return nil, errors.New("initialize sonyflake failed")
	}
	log.
		WithField("machine.id", machineID).
		Info("initialize sonyflake success")
	return sf, nil
}

{{- if not (and .Computed.enable_redis_final .Computed.enable_cache_final) }}
type nopCache struct{}

func (c *nopCache) Cache() redis.UniversalClient { return nil }
func (c *nopCache) WithPrefix(string) biz.Cache { return c }
func (c *nopCache) WithRefresh() biz.Cache      { return c }
func (c *nopCache) Get(ctx context.Context, _ string, write func(context.Context) (string, error)) (string, error) {
	if write == nil {
		return "", nil
	}
	return write(ctx)
}
func (c *nopCache) Set(context.Context, string, string, bool)         {}
func (c *nopCache) Del(context.Context, string)                       {}
func (c *nopCache) SetWithExpiration(context.Context, string, string, int64) {
}
func (c *nopCache) Flush(ctx context.Context, handler func(ctx context.Context) error) error {
	if handler == nil {
		return nil
	}
	return handler(ctx)
}
func (c *nopCache) FlushByPrefix(context.Context, ...string) error { return nil }

// NewCache returns a no-op Cache when Redis cache is disabled.
func NewCache() biz.Cache { return &nopCache{} }
{{- end }}

// NewRedis initializes Redis client from config.
func NewRedis(c *conf.Bootstrap) (redis.UniversalClient, error) {
	return newRedis(c)
}

// newRedis is the shared Redis initialization logic.
func newRedis(c *conf.Bootstrap) (client redis.UniversalClient, err error) {
	if c == nil || c.Redis == nil || c.Redis.Dsn == "" {
		err = errors.New("redis config is required")
		log.WithError(err).Error("initialize redis failed")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var u *url.URL
	u, err = url.Parse(c.Redis.Dsn)
	if err != nil {
		log.Error(err)
		err = errors.New("initialize redis failed")
		return
	}
	if u.User != nil {
		u.User = url.UserPassword(u.User.Username(), "***")
	}
	showDsn, _ := url.PathUnescape(u.String())
	client, err = utils.ParseRedisURI(c.Redis.Dsn)
	if err != nil {
		log.Error(err)
		err = errors.New("initialize redis failed")
		return
	}
	err = client.Ping(ctx).Err()
	if err != nil {
		log.Error(err)
		err = errors.New("initialize redis failed")
		return
	}
	log.
		WithField("redis.dsn", showDsn).
		Info("initialize redis success")
	return
}

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewSonyflake,
	{{- if .Computed.enable_redis_final }}
	NewRedis,
	{{- end }}
	{{- if .Computed.enable_trace_final }}
	NewTracer,
	{{- end }}
	NewCache,
	NewTransaction,
	NewAuthRepo,
	NewUserRepo,
	NewRoleRepo,
	NewPermissionRepo,
	{{- if .Computed.enable_action_final }}
	NewActionRepo,
	{{- end }}
	{{- if .Computed.enable_user_group_final }}
	NewUserGroupRepo,
	{{- end }}
	{{- if .Computed.enable_whitelist_final }}
	NewWhitelistRepo,
	{{- end }}
	{{- if .Computed.enable_hotspot_final }}
	NewHotspotRepo,
	{{- end }}
)
