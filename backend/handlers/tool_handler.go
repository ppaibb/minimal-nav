package handlers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// FetchFavicon 探测网站 Favicon
func FetchFavicon(c *gin.Context) {
	targetURL := c.Query("url")
	if targetURL == "" {
		Error(c, 400, "URL 参数不能为空")
		return
	}

	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	u, err := url.Parse(targetURL)
	if err != nil || u.Host == "" {
		Error(c, 400, "无效的 URL 格式")
		return
	}

	domain := u.Hostname()

	// 优先提供 Google Favicon 64px 高清源
	googleFavicon := fmt.Sprintf("https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=%s&size=64", url.QueryEscape(targetURL))
	directFavicon := fmt.Sprintf("%s://%s/favicon.ico", u.Scheme, u.Host)

	Success(c, gin.H{
		"domain":  domain,
		"favicon": googleFavicon,
		"direct":  directFavicon,
	})
}

// PingRequest 探测请求
type PingRequest struct {
	URL string `json:"url" binding:"required"`
}

// PingURL 连通性探测与延迟测量
func PingURL(c *gin.Context) {
	var req PingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "请提供要探测的 URL")
		return
	}

	targetURL := req.URL
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 支持内部自签证书
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	startTime := time.Now()
	resp, err := client.Get(targetURL)
	latency := time.Since(startTime).Milliseconds()

	if err != nil {
		Success(c, gin.H{
			"url":         req.URL,
			"healthy":     false,
			"status_code": 0,
			"latency_ms":  latency,
			"error":       err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	isHealthy := resp.StatusCode >= 200 && resp.StatusCode < 400

	Success(c, gin.H{
		"url":         req.URL,
		"healthy":     isHealthy,
		"status_code": resp.StatusCode,
		"latency_ms":  latency,
		"error":       "",
	})
}
