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

-- 4. 插入默认初始化种子数据 (如果表为空)
INSERT OR IGNORE INTO links (id, title, url, category, icon) VALUES 
(1, 'GitHub', 'https://github.com', '开发协作', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fgithub.com&size=64'),
(2, 'Cloudflare', 'https://cloudflare.com', '部署运维', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fcloudflare.com&size=64'),
(3, 'Tailwind CSS', 'https://tailwindcss.com', '设计资源', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Ftailwindcss.com&size=64'),
(4, 'Vue.js', 'https://vuejs.org', '开发框架', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fvuejs.org&size=64'),
(5, 'Vercel', 'https://vercel.com', '部署运维', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fvercel.com&size=64'),
(6, 'Go Dev', 'https://go.dev', '开发框架', 'https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fgo.dev&size=64');

INSERT OR IGNORE INTO announcements (id, content, is_active) VALUES 
(1, '欢迎使用企业极简导航系统！点击右上角管理后台可自定义链接和公告。', 1),
(2, '已全量启用 Cloudflare Serverless 全球边缘加速架构与 D1 数据库。', 1);

INSERT OR IGNORE INTO settings (key, value) VALUES 
('site_name', 'Minimal Nav'),
('site_desc', '企业极简导航系统 · Cloudflare Serverless 版');
