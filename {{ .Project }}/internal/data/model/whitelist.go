package model

const TableNameWhitelist = "whitelist"

// Whitelist mapped from table <whitelist>
type Whitelist struct {
	ID       uint64  `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id,string"`
	Category *int16  `gorm:"column:category;type:smallint;not null" json:"category"`
	Resource *string `gorm:"column:resource;type:text;not null" json:"resource"`
}

// TableName Whitelist's table name
func (*Whitelist) TableName() string {
	return TableNameWhitelist
}

