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

INSERT OR IGNORE INTO announcements (id, content, is_active) VALUES 
(1, '团队常用开发、协作、数据库、IDE 与 AI 助手导航已整理上线。', 1),
(2, '支持搜索、分类过滤与一键直达，点击右上角后台可自定义拓展。', 1);

INSERT OR IGNORE INTO settings (key, value) VALUES 
('site_name', '团队常用工具推荐'),
('site_desc', '团队常用开发、协作、数据库、IDE 与 AI 工具导航');
