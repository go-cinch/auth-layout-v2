package model

import (
	"github.com/golang-module/carbon/v2"
)

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID        uint64           `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:auto increment id" json:"id,string"`            // auto increment id
	CreatedAt *carbon.DateTime `gorm:"column:created_at;type:timestamp(3) without time zone;comment:create time" json:"created_at"`               // create time
	UpdatedAt *carbon.DateTime `gorm:"column:updated_at;type:timestamp(3) without time zone;comment:update time" json:"updated_at"`               // update time
	RoleID    *uint64          `gorm:"column:role_id;type:bigint;index:user_role_id_idx;comment:role id" json:"role_id,string"`                   // role id
	Action    *string          `gorm:"column:action;type:text;comment:action" json:"action"`                                                      // action
	Username  *string          `gorm:"column:username;type:character varying(191);uniqueIndex:user_username_uq;comment:username" json:"username"` // username
	Code      string           `gorm:"column:code;type:character(8);not null;uniqueIndex:user_code_uq;comment:code" json:"code"`                  // code
	Password  *string          `gorm:"column:password;type:text;comment:password" json:"password"`                                                // password
	Platform  *string          `gorm:"column:platform;type:character varying(50);comment:platform" json:"platform"`                               // platform
	LastLogin *carbon.DateTime `gorm:"column:last_login;type:timestamp(3) without time zone;comment:last login" json:"last_login"`                // last login

	Locked *int16 `gorm:"column:locked;type:smallint;default:0;comment:locked" json:"locked"` // locked
	{{ if .Computed.enable_user_lock_final }}
	LockExpire *int64 `gorm:"column:lock_expire;type:bigint;comment:lock expire" json:"lock_expire,string"` // lock expire
	{{ end }}
	{{ if or .Computed.enable_captcha_final .Computed.enable_user_lock_final }}
	Wrong *uint64 `gorm:"column:wrong;type:bigint;comment:wrong" json:"wrong,string"` // wrong
	{{ end }}
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
