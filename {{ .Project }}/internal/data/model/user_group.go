package model

const TableNameUserGroup = "user_group"

// UserGroup mapped from table <user_group>
type UserGroup struct {
	ID     uint64  `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:auto increment id" json:"id,string"` // auto increment id
	Name   *string `gorm:"column:name;type:character varying(50);comment:name" json:"name"`                                // name
	Word   *string `gorm:"column:word;type:character varying(50);uniqueIndex:user_group_word_uq;comment:word" json:"word"` // word
	Action *string `gorm:"column:action;type:text;comment:action" json:"action"`                                           // action
}

// TableName UserGroup's table name
func (*UserGroup) TableName() string {
	return TableNameUserGroup
}

const TableNameUserUserGroupRelation = "user_user_group_relation"

// UserUserGroupRelation mapped from table <user_user_group_relation>
type UserUserGroupRelation struct {
	UserID      *uint64 `gorm:"column:user_id;type:bigint;primaryKey;not null;comment:user id" json:"user_id,string"`                                                                    // user id
	UserGroupID *uint64 `gorm:"column:user_group_id;type:bigint;primaryKey;not null;index:user_user_group_relation_user_group_id_idx;comment:user group id" json:"user_group_id,string"` // user group id
}

// TableName UserUserGroupRelation's table name
func (*UserUserGroupRelation) TableName() string {
	return TableNameUserUserGroupRelation
}
