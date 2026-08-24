package db

import (
	"log"
	"os"
	"path/filepath"

	"minimal-nav/backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化 SQLite 数据库并执行自动迁移
func InitDB(dbPath string) {
	if dbPath == "" {
		dbPath = "nav.db"
	}

	// 确保数据库文件所在目录存在
	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("无法创建数据库目录: %v", err)
		}
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接 SQLite 数据库失败: %v", err)
	}

	// 自动迁移数据表
	err = DB.AutoMigrate(&models.Link{}, &models.Announcement{}, &models.Setting{})
	if err != nil {
		log.Fatalf("数据表迁移失败: %v", err)
	}

	log.Println("SQLite 数据库初始化及表迁移成功！")

	// 插入种子数据（如果表中无数据）
	seedInitialData()
}

func seedInitialData() {
	var linkCount int64
	DB.Model(&models.Link{}).Count(&linkCount)
	if linkCount == 0 {
		initialLinks := []models.Link{
			{Title: "GitHub", URL: "https://github.com", Category: "开发协作"},
			{Title: "Vercel", URL: "https://vercel.com", Category: "部署运维"},
			{Title: "Tailwind CSS", URL: "https://tailwindcss.com", Category: "设计资源"},
			{Title: "Vue.js", URL: "https://vuejs.org", Category: "开发框架"},
			{Title: "Go Dev", URL: "https://go.dev", Category: "开发框架"},
			{Title: "Cloudflare", URL: "https://cloudflare.com", Category: "部署运维"},
		}
		DB.Create(&initialLinks)
		log.Println("已初始化默认导航链接数据")
	}

	var announcementCount int64
	DB.Model(&models.Announcement{}).Count(&announcementCount)
	if announcementCount == 0 {
		initialAnnouncements := []models.Announcement{
			{Content: "欢迎使用企业极简导航系统！点击右上角管理后台可自定义链接和公告。", IsActive: true},
			{Content: "系统已全量启用极简暗黑/浅色自适应主题与毫秒级即时响应。", IsActive: true},
		}
		DB.Create(&initialAnnouncements)
		log.Println("已初始化默认公告数据")
	}
}
