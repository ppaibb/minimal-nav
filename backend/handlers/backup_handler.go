package handlers

import (
	"io"
	"regexp"
	"strings"
	"time"

	"minimal-nav/backend/db"
	"minimal-nav/backend/models"

	"github.com/gin-gonic/gin"
)

// BackupData 备份数据结构
type BackupData struct {
	Version       string                `json:"version"`
	ExportedAt    string                `json:"exported_at"`
	Links         []models.Link         `json:"links"`
	Announcements []models.Announcement `json:"announcements"`
}

// ExportBackup 导出备份数据为 JSON
func ExportBackup(c *gin.Context) {
	var links []models.Link
	db.DB.Find(&links)

	var announcements []models.Announcement
	db.DB.Find(&announcements)

	backup := BackupData{
		Version:       "1.0",
		ExportedAt:    time.Now().Format(time.RFC3339),
		Links:         links,
		Announcements: announcements,
	}

	c.Header("Content-Disposition", "attachment; filename=minimal-nav-backup.json")
	c.Header("Content-Type", "application/json")
	c.JSON(200, backup)
}

// ImportBackupRequest 导入备份请求
type ImportBackupRequest struct {
	Mode          string                `json:"mode"` // "merge" 或 "overwrite"
	Links         []models.Link         `json:"links"`
	Announcements []models.Announcement `json:"announcements"`
}

// ImportBackup 从 JSON 恢复备份数据
func ImportBackup(c *gin.Context) {
	var req ImportBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "无效的备份数据格式")
		return
	}

	if req.Mode == "overwrite" {
		db.DB.Exec("DELETE FROM links")
		db.DB.Exec("DELETE FROM announcements")
	}

	// 插入链接 (重置 ID)
	for _, l := range req.Links {
		newLink := models.Link{
			Title:    l.Title,
			URL:      l.URL,
			Category: l.Category,
			Icon:     l.Icon,
		}
		if newLink.Category == "" {
			newLink.Category = "Default"
		}
		db.DB.Create(&newLink)
	}

	// 插入公告
	for _, a := range req.Announcements {
		newAnnouncement := models.Announcement{
			Content:  a.Content,
			IsActive: a.IsActive,
		}
		db.DB.Create(&newAnnouncement)
	}

	Success(c, gin.H{
		"imported_links":         len(req.Links),
		"imported_announcements": len(req.Announcements),
	})
}

// BookmarkItem 解析出的书签条目
type BookmarkItem struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Category string `json:"category"`
}

// parseNetscapeBookmarks 解析 Netscape Bookmark HTML
func parseNetscapeBookmarks(htmlContent string) []BookmarkItem {
	var items []BookmarkItem
	lines := strings.Split(htmlContent, "\n")

	currentCategory := "常用书签"
	categoryStack := []string{currentCategory}

	h3Regex := regexp.MustCompile(`(?i)<H3[^>]*>(.*?)</H3>`)
	aRegex := regexp.MustCompile(`(?i)<A\s+HREF="([^"]+)"[^>]*>(.*?)</A>`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 检查进入新目录
		if strings.Contains(strings.ToUpper(trimmed), "<DL>") {
			// 进入下一层
		} else if strings.Contains(strings.ToUpper(trimmed), "</DL>") {
			// 退出当前层
			if len(categoryStack) > 1 {
				categoryStack = categoryStack[:len(categoryStack)-1]
				currentCategory = categoryStack[len(categoryStack)-1]
			}
		}

		// 检查分类标题 <H3>
		if match := h3Regex.FindStringSubmatch(trimmed); len(match) > 1 {
			catName := strings.TrimSpace(match[1])
			if catName != "" && !strings.EqualFold(catName, "Bookmarks Bar") && !strings.EqualFold(catName, "书签栏") {
				currentCategory = catName
				categoryStack = append(categoryStack, currentCategory)
			}
		}

		// 检查链接 <A HREF="...">
		if match := aRegex.FindStringSubmatch(trimmed); len(match) > 2 {
			linkURL := strings.TrimSpace(match[1])
			linkTitle := strings.TrimSpace(match[2])

			// 过滤非 http/https 协议（如 javascript: 或 place:）
			if strings.HasPrefix(linkURL, "http://") || strings.HasPrefix(linkURL, "https://") {
				if linkTitle == "" {
					linkTitle = linkURL
				}
				items = append(items, BookmarkItem{
					Title:    linkTitle,
					URL:      linkURL,
					Category: currentCategory,
				})
			}
		}
	}

	return items
}

// ImportBookmarks 导入 Chrome/Edge 书签 HTML 文件
func ImportBookmarks(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		Error(c, 400, "请上传有效的 HTML 书签文件")
		return
	}

	src, err := file.Open()
	if err != nil {
		Error(c, 500, "无法读取上传的文件")
		return
	}
	defer src.Close()

	contentBytes, err := io.ReadAll(src)
	if err != nil {
		Error(c, 500, "读取书签内容失败")
		return
	}

	bookmarks := parseNetscapeBookmarks(string(contentBytes))
	if len(bookmarks) == 0 {
		Error(c, 400, "未能从文件中解析出有效的网址书签")
		return
	}

	// 批量插入数据库
	count := 0
	for _, bm := range bookmarks {
		// 生成推荐 Favicon
		favicon := "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=" + bm.URL + "&size=64"
		link := models.Link{
			Title:    bm.Title,
			URL:      bm.URL,
			Category: bm.Category,
			Icon:     favicon,
		}
		db.DB.Create(&link)
		count++
	}

	Success(c, gin.H{
		"imported_count": count,
	})
}
