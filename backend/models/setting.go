package models

// Setting 系统设置键值对模型
type Setting struct {
	Key   string `gorm:"primaryKey;type:varchar(64)" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}
