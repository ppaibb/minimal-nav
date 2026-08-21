package handlers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"minimal-nav/backend/db"
	"minimal-nav/backend/models"

	"github.com/gin-gonic/gin"
)

// GetLinks 获取所有导航链接
func GetLinks(c *gin.Context) {
	var links []models.Link
	result := db.DB.Order("id asc").Find(&links)
	if result.Error != nil {
		Error(c, 500, "获取链接失败: "+result.Error.Error())
		return
	}
	Success(c, links)
}

// CreateLink 添加一条新链接 (自动填充 Favicon)
func CreateLink(c *gin.Context) {
	var input models.Link
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if input.Category == "" {
		input.Category = "Default"
	}

	// 若未指定图标，自动根据 URL 补充 Favicon
	if input.Icon == "" && input.URL != "" {
		targetURL := input.URL
		if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
			targetURL = "https://" + targetURL
		}
		input.Icon = fmt.Sprintf("https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=%s&size=64", url.QueryEscape(targetURL))
	}

	result := db.DB.Create(&input)
	if result.Error != nil {
		Error(c, 500, "创建链接失败: "+result.Error.Error())
		return
	}

	Success(c, input)
}

// UpdateLink 更新导航链接信息
func UpdateLink(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		Error(c, 400, "无效的 ID 参数")
		return
	}

	var existing models.Link
	if err := db.DB.First(&existing, id).Error; err != nil {
		Error(c, 404, "未找到该链接")
		return
	}

	var input models.Link
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	existing.Title = input.Title
	existing.URL = input.URL
	if input.Category != "" {
		existing.Category = input.Category
	}
	existing.Icon = input.Icon

	if err := db.DB.Save(&existing).Error; err != nil {
		Error(c, 500, "更新链接失败: "+err.Error())
		return
	}

	Success(c, existing)
}

// DeleteLink 删除指定 ID 的链接
func DeleteLink(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		Error(c, 400, "无效的 ID 参数")
		return
	}

	result := db.DB.Delete(&models.Link{}, id)
	if result.Error != nil {
		Error(c, 500, "删除链接失败: "+result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		Error(c, 404, "未找到该链接")
		return
	}

	Success(c, gin.H{"deleted_id": id})
}
