# Minimal Nav

极简冷峻风格的个人与团队导航系统。基于 Cloudflare 全无服务器架构（Pages + Functions + D1 数据库）构建，无需云服务器或容器环境，零成本永久运行。

---

## 核心特性

- **冷峻极简排版**：采用瑞士平面设计风格（Swiss Line Layout），以极细线条与严谨字阶构建空间，无多余装饰。
- **全边缘 Serverless 架构**：
  - 前端：托管于 Cloudflare Pages 全球 CDN
  - 后端：基于 TypeScript 与 Hono 框架驱动的 Pages Functions
  - 数据库：基于 Cloudflare D1 (分布式 SQLite)
- **链接管理**：支持多分类过滤、排序、置顶及 Favicon 自动解析。
- **公告系统**：3 行垂直轮播与分页控制器，支持 Markdown 详情展示。
- **网络工具**：提供 Favicon 提取与实时 HTTP 连通性探测 (Ping)。
- **数据迁移**：支持全量 JSON 数据导入导出，以及 Chrome / Edge HTML 书签一键解析导入。
- **无状态鉴权**：基于 Web Crypto HMAC-SHA256 签名，支持多节点边缘分布式验证。

---

## 架构示意

```
[ 用户浏览器 ]
      │
[ Cloudflare 边缘网络 ]
   ├── [ Cloudflare Pages (Vue 3 静态前端) ]
   └── [ Pages Functions (Hono API - /api/*) ]
            │
       [ Cloudflare D1 (分布式 SQLite) ]
```

---

## 部署流程 (Cloudflare)

### 1. 创建 D1 数据库
1. 登录 Cloudflare 控制台，进入 **存储和数据库** -> **D1 SQL 数据库**，创建数据库 `minimal-nav-db`。
2. 在数据库 **控制台 (Console)** 中执行 `d1/schema.sql` 完成表结构与初始数据初始化。

### 2. 部署 Pages
1. 在控制台进入 **Workers 和 Pages** -> **创建应用程序** -> **Pages** -> **连接到 Git**。
2. 选择本仓库，设置构建参数：
   - 框架预设: `Vite`
   - 根目录: `frontend`
   - 构建命令: `npm run build`
   - 输出目录: `dist`
3. 保存并部署。

### 3. 绑定 D1 与环境变量
1. 进入 Pages 项目的 **设置** -> **函数 (Functions)**。
2. 添加 **D1 数据库绑定**：变量名设为 `DB`，选择 `minimal-nav-db`。
3. 添加 **环境变量**：变量名设为 `ADMIN_PASSWORD`，填入管理员密码（默认为 `admin123`）。
4. 重新部署即可上线。

---

## 本地开发

```bash
# 安装依赖
npm install
cd frontend && npm install && cd ..

# 初始化本地 D1 数据库
npm run d1:init:local

# 构建并启动本地仿真预览
npm run build:frontend
npm run preview
```
本地访问地址：`http://localhost:8788`。

---

## 目录结构

```
minimal-nav/
├── d1/schema.sql        # D1 数据表结构与种子数据
├── functions/api/       # Cloudflare Functions (Hono API)
├── frontend/            # Vue 3 前端源码
├── CLOUDFLARE_DEPLOY.md # 详细部署说明
├── wrangler.toml        # Cloudflare 配置文件
└── package.json
```

---

## 开源协议

[MIT License](LICENSE)
