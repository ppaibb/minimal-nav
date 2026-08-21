package handlers

import (
	"net/http"
	"strconv"

	"minimal-nav/backend/db"
	"minimal-nav/backend/models"

	"github.com/gin-gonic/gin"
)

// GetLinks 获取所有导航链接
func GetLinks(c *gin.Context) {
	var links []models.Link
	result := db.DB.Order("id asc").Find(&links)
	if result.Error != nil {
		Error(c, http.StatusInternalServerError, 500, "获取链接失败: "+result.Error.Error())
		return
	}
	Success(c, links)
}

// CreateLink 添加一条新链接
func CreateLink(c *gin.Context) {
	var input models.Link
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, http.StatusBadRequest, 400, "请求参数错误: "+err.Error())
		return
	}

	result := db.DB.Create(&input)
	if result.Error != nil {
		Error(c, http.StatusInternalServerError, 500, "创建链接失败: "+result.Error.Error())
		return
	}

	Success(c, input)
}

// DeleteLink 删除指定 ID 的链接
func DeleteLink(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		Error(c, http.StatusBadRequest, 400, "无效的 ID 参数")
		return
	}

	result := db.DB.Delete(&models.Link{}, id)
	if result.Error != nil {
		Error(c, http.StatusInternalServerError, 500, "删除链接失败: "+result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		Error(c, http.StatusNotFound, 404, "未找到该链接")
		return
	}

	Success(c, gin.H{"deleted_id": id})
}
