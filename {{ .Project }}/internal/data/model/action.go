package model

const TableNameAction = "action"

// Action mapped from table <action>
type Action struct {
	ID       uint64  `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id,string"`
	Name     *string `gorm:"column:name;type:character varying(50)" json:"name"`
	Code     *string `gorm:"column:code;type:character(8);not null;uniqueIndex:action_code_uq" json:"code"`
	Word     *string `gorm:"column:word;type:character varying(50);uniqueIndex:action_word_uq" json:"word"`
	Resource *string `gorm:"column:resource;type:text" json:"resource"`
	Menu     *string `gorm:"column:menu;type:text" json:"menu"`
	Btn      *string `gorm:"column:btn;type:text" json:"btn"`
}

// TableName Action's table name
func (*Action) TableName() string {
	return TableNameAction
}

