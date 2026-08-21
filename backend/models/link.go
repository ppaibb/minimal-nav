package models

// Link 导航链接数据模型
type Link struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `gorm:"size:255;not null" json:"title" binding:"required"`
	URL      string `gorm:"size:1024;not null" json:"url" binding:"required"`
	Category string `gorm:"size:100;default:'Default'" json:"category"`
	Icon     string `gorm:"size:255" json:"icon,omitempty"`
}
