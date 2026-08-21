package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	tokenStore = make(map[string]time.Time)
	tokenMutex sync.RWMutex
)

// getAdminPassword 获取管理员口令 (默认 admin123)
func getAdminPassword() string {
	pwd := os.Getenv("ADMIN_PASSWORD")
	if pwd == "" {
		return "admin123"
	}
	return pwd
}

// generateToken 生成随机安全 Token
func generateToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// LoginRequest 登录请求
type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

// Login 管理员登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "请输入管理员口令")
		return
	}

	if req.Password != getAdminPassword() {
		Error(c, 401, "管理员口令错误")
		return
	}

	token := generateToken()
	tokenMutex.Lock()
	tokenStore[token] = time.Now().Add(24 * 7 * time.Hour) // 7天有效期
	tokenMutex.Unlock()

	Success(c, gin.H{
		"token":      token,
		"expires_in": 24 * 7 * 3600,
	})
}

// CheckAuth 校验当前 Token 是否有效
func CheckAuth(c *gin.Context) {
	Success(c, gin.H{
		"authenticated": true,
	})
}

// AuthMiddleware 管理后台保护中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Admin-Token")
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			Error(c, 401, "请先解锁管理权限")
			c.Abort()
			return
		}

		tokenMutex.RLock()
		expireAt, exists := tokenStore[token]
		tokenMutex.RUnlock()

		if !exists || time.Now().After(expireAt) {
			Error(c, 401, "登录会话已过期，请重新输入口令")
			c.Abort()
			return
		}

		c.Next()
	}
}
