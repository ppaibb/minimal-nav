package models

import "time"

// Announcement 系统公告数据模型
type Announcement struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content" binding:"required"`
	DetailMD  string    `gorm:"type:text" json:"detail_md"` // 长篇 Markdown 详情内容
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
