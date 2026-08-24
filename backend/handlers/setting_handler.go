package handlers

import (
	"minimal-nav/backend/db"
	"minimal-nav/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// GetSettings 获取所有系统设置
func GetSettings(c *gin.Context) {
	var list []models.Setting
	db.DB.Find(&list)

	settingsMap := map[string]string{
		"site_name": "Minimal Nav",
		"site_desc": "统一汇聚团队核心工具、部署控制台、设计协作及文档中心，即时检索快速直达。",
	}

	for _, item := range list {
		if item.Value != "" {
			settingsMap[item.Key] = item.Value
		}
	}

	Success(c, settingsMap)
}

// UpdateSettings 批量更新系统设置 (受口令保护)
func UpdateSettings(c *gin.Context) {
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		Error(c, 400, "请求参数格式错误: "+err.Error())
		return
	}

	for k, v := range input {
		setting := models.Setting{Key: k, Value: v}
		db.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Create(&setting)
	}

	// 返回更新后的完整设置
	GetSettings(c)
}
