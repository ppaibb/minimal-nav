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
		Logger: logger.Default.LogMode(logger.Warn),
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
			// 1. 开发与 IDE
			{Title: "Visual Studio Code", URL: "https://code.visualstudio.com/", Category: "开发与 IDE", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fcode.visualstudio.com&size=64"},
			{Title: "TRAE CN", URL: "https://www.trae.cn/", Category: "开发与 IDE", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Ftrae.cn&size=64"},
			{Title: "CodeBuddy CN", URL: "https://www.codebuddy.cn/home/", Category: "开发与 IDE", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fcodebuddy.cn&size=64"},
			{Title: "Qoder CN", URL: "https://qoder.com.cn/", Category: "开发与 IDE", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fqoder.com&size=64"},
			{Title: "Node.js 24 LTS", URL: "https://nodejs.org/zh-cn/download", Category: "开发与 IDE", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fnodejs.org&size=64"},
			{Title: "OpenJDK 17", URL: "https://adoptium.net/zh-CN/temurin/releases?version=17&os=any&arch=any", Category: "开发与 IDE", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fadoptium.net&size=64"},
			{Title: "OpenJDK 21", URL: "https://adoptium.net/zh-CN/temurin/releases?version=21&os=any&arch=any", Category: "开发与 IDE", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fadoptium.net&size=64"},
			{Title: "Oracle JDK 8", URL: "https://repo.huaweicloud.com:8443/artifactory/java-local/jdk/8u202-b08/", Category: "开发与 IDE", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fhuaweicloud.com&size=64"},

			// 2. AI 智能助手
			{Title: "TRAE Work", URL: "https://www.trae.cn/ide/download", Category: "AI 智能助手", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Ftrae.cn&size=64"},
			{Title: "WorkBuddy", URL: "https://www.codebuddy.cn/work/", Category: "AI 智能助手", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fcodebuddy.cn&size=64"},
			{Title: "QoderWork CN", URL: "https://qoder.com.cn/qoderwork", Category: "AI 智能助手", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fqoder.com&size=64"},
			{Title: "ChatGPT (Codex Desktop)", URL: "https://learn.chatgpt.com/", Category: "AI 智能助手", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fchatgpt.com&size=64"},
			{Title: "Claude Desktop", URL: "https://claude.com/download", Category: "AI 智能助手", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fclaude.com&size=64"},
			{Title: "Claude Code CLI", URL: "https://claude.com/product/claude-code", Category: "AI 智能助手", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fclaude.com&size=64"},
			{Title: "Codex CLI", URL: "https://learn.chatgpt.com/docs/codex/cli#getting-started", Category: "AI 智能助手", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fchatgpt.com&size=64"},

			// 3. 数据与接口
			{Title: "Apifox", URL: "https://apifox.com/", Category: "数据与接口", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fapifox.com&size=64"},
			{Title: "DBX", URL: "https://dbxio.com/", Category: "数据与接口", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fdbxio.com&size=64"},
			{Title: "Navicat Premium", URL: "https://www.navicat.com.cn/download/navicat-premium", Category: "数据与接口", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fnavicat.com.cn&size=64"},
			{Title: "SQLark (达梦专用)", URL: "https://www.sqlark.com/", Category: "数据与接口", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fsqlark.com&size=64"},
			{Title: "Tiny RDM (Redis)", URL: "https://redis.tinycraft.cc/zh/", Category: "数据与接口", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fredis.tinycraft.cc&size=64"},
			{Title: "MQTTX", URL: "https://mqttx.app/", Category: "数据与接口", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fmqttx.app&size=64"},

			// 4. 网络与运维
			{Title: "网易 UU 远程", URL: "https://uuyc.163.com/", Category: "网络与运维", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fuuyc.163.com&size=64"},
			{Title: "Netcatty 终端", URL: "https://netcatty.app/en/", Category: "网络与运维", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fnetcatty.app&size=64"},
			{Title: "Xshell", URL: "https://www.xshell.com/zh/free-for-home-school/", Category: "网络与运维", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fxshell.com&size=64"},
			{Title: "WireGuard", URL: "https://www.wireguard.com/install/", Category: "网络与运维", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fwireguard.com&size=64"},
			{Title: "OpenVPN", URL: "https://openvpn.net/community/", Category: "网络与运维", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fopenvpn.net&size=64"},

			// 5. 文档与办公
			{Title: "语雀", URL: "https://www.yuque.com/", Category: "文档与办公", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fyuque.com&size=64"},
			{Title: "Sublime Text", URL: "https://www.sublimetext.com/", Category: "文档与办公", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fsublimetext.com&size=64"},

			// 6. Docker 镜像加速
			{Title: "自建镜像加速节点", URL: "https://cr.gua.cx/", Category: "Docker 镜像加速", Icon: "https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fcr.gua.cx&size=64"},
		}
		DB.Create(&initialLinks)
		log.Println("已初始化默认精选导航链接数据")
	}

	var announcementCount int64
	DB.Model(&models.Announcement{}).Count(&announcementCount)
	if announcementCount == 0 {
		initialAnnouncements := []models.Announcement{
			{
				Content: "🔥 AI 编程助手接入指南（Claude Code / Codex / WorkBuddy）",
				DetailMD: `# AI 编程助手接入指南

团队统一接入 724AI 镜像中转网络，为 Claude Code、Codex 及 WorkBuddy 客户端提供稳定免翻墙的极速编码服务。

---

## 准备工作与注意事项

1. **专属 Key**：每个人使用独立的专属 API Key，请勿混用或共享。开通权限或获取 Key 请在飞书联系管理员。
2. **分组说明**：当前 Key 区分 Claude 与 Codex 独立路由分组，如需跨组使用请联系管理员调整。

---

## 一、Claude Code 接入

将 Claude Code 的默认请求地址替换为 724AI 中转站（` + "`" + `api.724ai.org` + "`" + `），直连稳定无需代理。

### 1. 安装 Claude Code
确保已安装 Node.js 18+ 环境：
` + "```bash" + `
npm install -g @anthropic-ai/claude-code
` + "```" + `

### 2. 配置 724AI 接入
编辑配置文件 ` + "`~/.claude/settings.json`" + `：
` + "```json" + `
{
  "env": {
    "ANTHROPIC_API_KEY": "sk-xxxxxxxxxxxxxxxx",
    "ANTHROPIC_BASE_URL": "https://api.724ai.org",
    "CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS": "1"
  }
}
` + "```" + `
> **注意**：` + "`ANTHROPIC_BASE_URL`" + ` 必须填写 ` + "`https://api.724ai.org`" + `，末尾严禁加斜杠 ` + "`/`" + `。

---

## 二、Codex CLI & App 接入

### 1. 安装 Codex
- **桌面客户端**：前往官方页面下载对应平台版本。
- **CLI 命令行**：
` + "```bash" + `
npm install -g @openai/codex
` + "```" + `

---

## 三、WorkBuddy 接入

WorkBuddy 支持兼容 OpenAI 协议的自定义模型，直接使用 724AI 通用接口地址：
- **接口地址**：` + "`https://api.724ai.org/v1`" + `
- **推荐模型**：` + "`gpt-5.6-sol`" + ` / ` + "`gpt-5.6-luna`" + `
`,
				IsActive: true,
			},
			{
				Content: "欢迎使用团队极简导航系统！支持单二进制运行与云原生部署",
				DetailMD: `# 欢迎使用团队极简导航系统

本项目是专为极客与团队打造的极简冷峻风格（Swiss Layout）导航与工具集成平台。

---

## 核心功能介绍

1. **多分类导航**：集中整理了开发工具、AI 助手、数据库、网络运维与常用文档。
2. **实时测速与 Favicon 探测**：内置边缘网络探测引擎，实时测试目标网址连通性并自动提取高清图标。
3. **数据导入与备份**：支持全量 JSON 数据导出，以及 Chrome / Edge 浏览器书签 HTML 一键批量解析导入。
4. **后台管理**：点击右上角后台管理即可进入自定义配置，所有数据实时保存到本地 SQLite 数据库。
`,
				IsActive: true,
			},
		}
		DB.Create(&initialAnnouncements)
		log.Println("已初始化默认公告数据")
	}

	var settingCount int64
	DB.Model(&models.Setting{}).Count(&settingCount)
	if settingCount == 0 {
		initialSettings := []models.Setting{
			{Key: "site_name", Value: "团队常用工具推荐"},
			{Key: "site_desc", Value: "团队常用开发、协作、数据库、IDE 与 AI 工具导航"},
			{Key: "icp_beian", Value: ""},
		}
		DB.Create(&initialSettings)
		log.Println("已初始化默认站点设置数据")
	}
}
