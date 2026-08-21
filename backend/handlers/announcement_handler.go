package handlers

import (
	"net/http"
	"strconv"

	"minimal-nav/backend/db"
	"minimal-nav/backend/models"

	"github.com/gin-gonic/gin"
)

// GetAnnouncements 获取所有公告
func GetAnnouncements(c *gin.Context) {
	var list []models.Announcement
	result := db.DB.Order("created_at desc").Find(&list)
	if result.Error != nil {
		Error(c, http.StatusInternalServerError, 500, "获取公告列表失败: "+result.Error.Error())
		return
	}
	Success(c, list)
}

// GetActiveAnnouncements 仅获取当前生效的公告
func GetActiveAnnouncements(c *gin.Context) {
	var list []models.Announcement
	result := db.DB.Where("is_active = ?", true).Order("created_at desc").Find(&list)
	if result.Error != nil {
		Error(c, http.StatusInternalServerError, 500, "获取生效公告失败: "+result.Error.Error())
		return
	}
	Success(c, list)
}

// CreateAnnouncement 发布一条新公告
func CreateAnnouncement(c *gin.Context) {
	var input models.Announcement
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusBadRequest, 400, "请求参数错误: "+err.Error())
		return
	}

	result := db.DB.Create(&input)
	if result.Error != nil {
		Error(c, http.StatusInternalServerError, 500, "发布公告失败: "+result.Error.Error())
		return
	}

	Success(c, input)
}

// ToggleAnnouncement 切换指定公告的生效状态
func ToggleAnnouncement(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		Error(c, http.StatusBadRequest, 400, "无效的 ID 参数")
		return
	}

	var item models.Announcement
	if err := db.DB.First(&item, id).Error; err != nil {
		Error(c, http.StatusNotFound, 404, "未找到该公告")
		return
	}

	item.IsActive = !item.IsActive
	if err := db.DB.Save(&item).Error; err != nil {
		Error(c, http.StatusInternalServerError, 500, "更新状态失败: "+err.Error())
		return
	}

	Success(c, item)
}

// DeleteAnnouncement 删除指定 ID 的公告 (辅助接口)
func DeleteAnnouncement(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		Error(c, http.StatusBadRequest, 400, "无效的 ID 参数")
		return
	}

	result := db.DB.Delete(&models.Announcement{}, id)
	if result.Error != nil {
		Error(c, http.StatusInternalServerError, 500, "删除公告失败: "+result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		Error(c, http.StatusNotFound, 404, "未找到该公告")
		return
	}

	Success(c, gin.H{"deleted_id": id})
}
