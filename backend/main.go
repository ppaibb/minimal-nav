package main

import (
	"log"
	"net/http"
	"os"

	"minimal-nav/backend/db"
	"minimal-nav/backend/handlers"

	"github.com/gin-gonic/gin"
)

// corsMiddleware 配置跨域中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Admin-Token, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// registerRoutes 注册 API 路由
func registerRoutes(r *gin.RouterGroup) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		handlers.Success(c, gin.H{
			"status": "up",
		})
	})

	// 1. 认证接口
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", handlers.Login)
		authGroup.GET("/check", handlers.AuthMiddleware(), handlers.CheckAuth)
	}

	// 2. 导航链接管理 (读公开，写保护)
	r.GET("/links", handlers.GetLinks)
	r.POST("/links", handlers.AuthMiddleware(), handlers.CreateLink)
	r.PUT("/links/:id", handlers.AuthMiddleware(), handlers.UpdateLink)
	r.DELETE("/links/:id", handlers.AuthMiddleware(), handlers.DeleteLink)

	// 3. 公告管理 (读公开，写保护)
	r.GET("/announcements", handlers.GetAnnouncements)
	r.GET("/announcements/active", handlers.GetActiveAnnouncements)
	r.GET("/announcements/:id", handlers.GetAnnouncementDetail)
	r.POST("/announcements", handlers.AuthMiddleware(), handlers.CreateAnnouncement)
	r.PUT("/announcements/:id", handlers.AuthMiddleware(), handlers.UpdateAnnouncement)
	r.PUT("/announcements/:id/toggle", handlers.AuthMiddleware(), handlers.ToggleAnnouncement)
	r.DELETE("/announcements/:id", handlers.AuthMiddleware(), handlers.DeleteAnnouncement)

	// 4. 工具能力 (Favicon 提取 & 连通性 Ping)
	toolGroup := r.Group("/tools")
	{
		toolGroup.GET("/favicon", handlers.FetchFavicon)
		toolGroup.POST("/ping", handlers.PingURL)
	}

	// 5. 数据备份与书签导入 (受口令保护)
	backupGroup := r.Group("/backup", handlers.AuthMiddleware())
	{
		backupGroup.GET("/export", handlers.ExportBackup)
		backupGroup.POST("/import", handlers.ImportBackup)
		backupGroup.POST("/import-bookmarks", handlers.ImportBookmarks)
	}

	// 6. 系统站点设置 (读公开，写受口令保护)
	r.GET("/settings", handlers.GetSettings)
	r.PUT("/settings", handlers.AuthMiddleware(), handlers.UpdateSettings)
}

func main() {
	// 1. 初始化数据库
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "nav.db"
	}
	db.InitDB(dbPath)

	// 2. 初始化 Gin 引擎
	r := gin.Default()
	r.Use(corsMiddleware())

	// 3. 注册路由 (同时支持根路径和 /api 前缀，完美适配直连与 Nginx 反代)
	registerRoutes(&r.RouterGroup)
	apiGroup := r.Group("/api")
	registerRoutes(apiGroup)

	// 4. 监听端口
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := host + ":" + port

	log.Printf("Minimal Nav 后端服务已启动，监听地址 http://%s ...\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v\n", err)
	}
}
