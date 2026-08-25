-- Minimal Nav Cloudflare D1 数据库架构

-- 1. 导航链接表
CREATE TABLE IF NOT EXISTS links (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    icon TEXT,
    category TEXT DEFAULT 'Default',
    sort_order INTEGER DEFAULT 0,
    is_pinned INTEGER DEFAULT 0
);

-- 2. 公告表
CREATE TABLE IF NOT EXISTS announcements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    content TEXT NOT NULL,
    detail_md TEXT,
    is_active INTEGER DEFAULT 1,
    sort_order INTEGER DEFAULT 0
);

-- 3. 系统配置表
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 4. 插入默认初始化种子数据 (精简归纳为 5 大核心类别)
INSERT OR REPLACE INTO links (id, title, url, category, icon, sort_order) VALUES 
-- 1. 开发与 IDE
(1, 'Visual Studio Code', 'https://code.visualstudio.com/', '开发与 IDE', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fcode.visualstudio.com&size=64', 1),
(2, 'TRAE CN', 'https://www.trae.cn/', '开发与 IDE', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Ftrae.cn&size=64', 2),
(3, 'CodeBuddy CN', 'https://www.codebuddy.cn/home/', '开发与 IDE', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fcodebuddy.cn&size=64', 3),
(4, 'Qoder CN', 'https://qoder.com.cn/', '开发与 IDE', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fqoder.com&size=64', 4),
(5, 'Node.js 24 LTS', 'https://nodejs.org/zh-cn/download', '开发与 IDE', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fnodejs.org&size=64', 5),
(6, 'OpenJDK 17', 'https://adoptium.net/zh-CN/temurin/releases?version=17&os=any&arch=any', '开发与 IDE', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fadoptium.net&size=64', 6),
(7, 'OpenJDK 21', 'https://adoptium.net/zh-CN/temurin/releases?version=21&os=any&arch=any', '开发与 IDE', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fadoptium.net&size=64', 7),
(8, 'Oracle JDK 8', 'https://repo.huaweicloud.com:8443/artifactory/java-local/jdk/8u202-b08/', '开发与 IDE', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fhuaweicloud.com&size=64', 8),

-- 2. AI 智能助手
(9, 'TRAE Work', 'https://www.trae.cn/ide/download', 'AI 智能助手', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Ftrae.cn&size=64', 9),
(10, 'WorkBuddy', 'https://www.codebuddy.cn/work/', 'AI 智能助手', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fcodebuddy.cn&size=64', 10),
(11, 'QoderWork CN', 'https://qoder.com.cn/qoderwork', 'AI 智能助手', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fqoder.com&size=64', 11),
(12, 'ChatGPT (Codex Desktop)', 'https://learn.chatgpt.com/', 'AI 智能助手', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fchatgpt.com&size=64', 12),
(13, 'Claude Desktop', 'https://claude.com/download', 'AI 智能助手', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fclaude.com&size=64', 13),
(14, 'Claude Code CLI', 'https://claude.com/product/claude-code', 'AI 智能助手', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fclaude.com&size=64', 14),
(15, 'Codex CLI', 'https://learn.chatgpt.com/docs/codex/cli#getting-started', 'AI 智能助手', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fchatgpt.com&size=64', 15),

-- 3. 数据与接口
(16, 'Apifox', 'https://apifox.com/', '数据与接口', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fapifox.com&size=64', 16),
(17, 'DBX', 'https://dbxio.com/', '数据与接口', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fdbxio.com&size=64', 17),
(18, 'Navicat Premium', 'https://www.navicat.com.cn/download/navicat-premium', '数据与接口', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fnavicat.com.cn&size=64', 18),
(19, 'SQLark (达梦专用)', 'https://www.sqlark.com/', '数据与接口', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fsqlark.com&size=64', 19),
(20, 'Tiny RDM (Redis)', 'https://redis.tinycraft.cc/zh/', '数据与接口', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fredis.tinycraft.cc&size=64', 20),
(21, 'MQTTX', 'https://mqttx.app/', '数据与接口', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fmqttx.app&size=64', 21),

-- 4. 网络与运维
(22, '网易 UU 远程', 'https://uuyc.163.com/', '网络与运维', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fuuyc.163.com&size=64', 22),
(23, 'Netcatty 终端', 'https://netcatty.app/en/', '网络与运维', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fnetcatty.app&size=64', 23),
(24, 'Xshell', 'https://www.xshell.com/zh/free-for-home-school/', '网络与运维', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fxshell.com&size=64', 24),
(25, 'WireGuard', 'https://www.wireguard.com/install/', '网络与运维', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fwireguard.com&size=64', 25),
(26, 'OpenVPN', 'https://openvpn.net/community/', '网络与运维', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fopenvpn.net&size=64', 26),

-- 5. 文档与办公
(27, '语雀', 'https://www.yuque.com/', '文档与办公', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fyuque.com&size=64', 27),
(28, 'Sublime Text', 'https://www.sublimetext.com/', '文档与办公', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fsublimetext.com&size=64', 28),

-- 6. Docker 镜像加速
(29, '自建镜像加速节点', 'https://cr.gua.cx/', 'Docker 镜像加速', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fcr.gua.cx&size=64', 29);

INSERT OR REPLACE INTO announcements (id, content, detail_md, is_active, sort_order) VALUES 
(1, 
'🔥 AI 编程助手接入指南（Claude Code / Codex / WorkBuddy）',
'# AI 编程助手接入指南

团队统一接入 724AI 镜像中转网络，为 Claude Code、Codex 及 WorkBuddy 客户端提供稳定免翻墙的极速编码服务。

---

## 准备工作与注意事项

1. **专属 Key**：每个人使用独立的专属 API Key，请勿混用或共享。开通权限或获取 Key 请在飞书联系管理员。
2. **分组说明**：当前 Key 区分 Claude 与 Codex 独立路由分组，如需跨组使用请联系管理员调整。

---

## 一、Claude Code 接入

将 Claude Code 的默认请求地址替换为 724AI 中转站（`api.724ai.org`），直连稳定无需代理。

### 1. 安装 Claude Code
确保已安装 Node.js 18+ 环境：
```bash
npm install -g @anthropic-ai/claude-code
```

### 2. 配置 724AI 接入
编辑配置文件 `~/.claude/settings.json`（Windows 对应 `C:\Users\用户名\.claude\settings.json`）：
```json
{
  "env": {
    "ANTHROPIC_API_KEY": "sk-xxxxxxxxxxxxxxxx",
    "ANTHROPIC_BASE_URL": "https://api.724ai.org",
    "CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS": "1"
  }
}
```
> **注意**：`ANTHROPIC_BASE_URL` 必须填写 `https://api.724ai.org`，末尾严禁加斜杠 `/`。

### 3. 跳过官方网页鉴权
在 `~/.claude.json` 中添加跳过标识：
```json
{
  "hasCompletedOnboarding": true
}
```

---

## 二、Codex CLI & App 接入

### 1. 安装 Codex
- **桌面客户端**：前往官方页面下载对应平台版本。
- **CLI 命令行**：
```bash
npm install -g @openai/codex
```

---

## 三、WorkBuddy 接入

WorkBuddy 支持兼容 OpenAI 协议的自定义模型，直接使用 724AI 通用接口地址：
- **接口地址**：`https://api.724ai.org/v1`
- **推荐模型**：`gpt-5.6-sol` / `gpt-5.6-luna`
',
1, 
1),

(2, 
'欢迎使用团队极简导航系统！支持云原生部署与自定义配置',
'# 欢迎使用团队极简导航系统

本项目是专为极客与团队打造的极简冷峻风格（Swiss Layout）导航与工具集成平台。

---

## 核心功能介绍

1. **多分类导航**：集中整理了开发工具、AI 助手、数据库、网络运维与常用文档。
2. **实时测速与 Favicon 探测**：内置边缘网络探测引擎，实时测试目标网址连通性并自动提取高清图标。
3. **数据导入与备份**：支持全量 JSON 数据导出，以及 Chrome / Edge 浏览器书签 HTML 一键批量解析导入。
4. **后台管理**：点击右上角后台管理即可进入自定义配置，所有数据实时保存到云端 D1 数据库。
',
1, 
2);

INSERT OR IGNORE INTO settings (key, value) VALUES 
('site_name', '团队常用工具推荐'),
('site_desc', '团队常用开发、协作、数据库、IDE 与 AI 工具导航');
