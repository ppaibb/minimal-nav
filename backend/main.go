package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"minimal-nav/backend/db"
	"minimal-nav/backend/handlers"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var embedDistFS embed.FS

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

// registerAPIRoutes 注册所有后端 API 路由
func registerAPIRoutes(r *gin.RouterGroup) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		handlers.Success(c, gin.H{
			"status":   "up",
			"platform": "single-binary-embed",
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

// setupEmbeddedFrontend 配置内嵌前端静态资源服务与 SPA Fallback
func setupEmbeddedFrontend(r *gin.Engine) {
	distSubFS, err := fs.Sub(embedDistFS, "dist")
	if err != nil {
		log.Printf("⚠️ 未能加载内嵌前端资源: %v\n", err)
		return
	}

	indexHTML, err := fs.ReadFile(distSubFS, "index.html")
	if err != nil {
		log.Printf("⚠️ 未能读取 index.html: %v\n", err)
	}

	fileServer := http.FileServer(http.FS(distSubFS))

	// 处理所有非 API 路由并支持 Vue Router History 模式
	r.NoRoute(func(c *gin.Context) {
		reqPath := c.Request.URL.Path

		// 如果是 API 请求直接返回 404 JSON
		if strings.HasPrefix(reqPath, "/api/") || reqPath == "/api" {
			handlers.Error(c, 404, "API endpoint not found")
			return
		}

		cleanPath := strings.TrimPrefix(path.Clean(reqPath), "/")
		if cleanPath == "" || cleanPath == "." || cleanPath == "index.html" {
			c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
			return
		}

		// 检查嵌入文件系统中是否存在该静态资源文件
		file, err := distSubFS.Open(cleanPath)
		if err == nil {
			stat, statErr := file.Stat()
			_ = file.Close()
			if statErr == nil && !stat.IsDir() {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}

		// 文件不存在时，回退到 index.html（供 Vue Router 处理客户端路由）
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})
}

func main() {
	// 1. 初始化数据库
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "nav.db"
	}
	db.InitDB(dbPath)

	// 2. 初始化 Gin 引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(corsMiddleware())

	// 3. 注册 API 路由 (/api 前缀)
	apiGroup := r.Group("/api")
	registerAPIRoutes(apiGroup)

	// 4. 挂载内嵌前端 SPA 静态文件服务
	setupEmbeddedFrontend(r)

	// 5. 监听端口
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := host + ":" + port

	log.Printf("====================================================\n")
	log.Printf("🧭 Minimal Nav 极简导航系统（单二进制版）已成功启动！\n")
	log.Printf("🌐 请在浏览器中打开: http://127.0.0.1:%s\n", port)
	log.Printf("🗄️ 数据存储路径: %s\n", dbPath)
	log.Printf("====================================================\n")

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v\n", err)
	}
}
