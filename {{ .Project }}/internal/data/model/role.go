package model

const TableNameRole = "role"

// Role mapped from table <role>
type Role struct {
	ID     uint64  `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id,string"`
	Name   *string `gorm:"column:name;type:character varying(50)" json:"name"`
	Word   *string `gorm:"column:word;type:character varying(50)" json:"word"`
	Action *string `gorm:"column:action;type:text" json:"action"`
}

// TableName Role's table name
func (*Role) TableName() string {
	return TableNameRole
}

