package biz

import "context"

// HotspotRepo provides cached lookup helpers used by the auth domain.
//
// It is intended for "hot" read paths (login/status/permission checks) where we
// want to avoid repeated DB reads. Implementations should keep data fresh via
// periodic/manual Refresh calls.
type HotspotRepo interface {
	// Refresh clears and/or rebuilds the in-memory hotspot caches.
	Refresh(ctx context.Context) error

	// User lookups.
	GetUserByCode(ctx context.Context, code string) *User
	GetUserByUsername(ctx context.Context, username string) *User

	// Permission/Action helpers (optional callers).
	GetActionByCode(ctx context.Context, code string) *Action
	GetUserPermissions(ctx context.Context, userID uint64) ([]*Action, error)
	CheckPermission(ctx context.Context, userID uint64, resource, method string) (bool, error)
}
