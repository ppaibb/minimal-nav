# Minimal Nav

极简冷峻风格的个人与团队导航系统。支持 **Cloudflare Serverless 全托管** 与 **本地 / 私有服务器自建部署** 两种运行形态，满足从零成本云端托管到企业内网私有化的全部场景。

---

## 核心特性

- **冷峻极简排版**：采用瑞士平面设计风格（Swiss Line Layout），以极细几何线条与严谨无衬线字阶构建界面，克制无冗余修饰。
- **双模部署支持**：
  - **Cloudflare Serverless 模式**：基于 Pages + Functions (Hono) + D1 (分布式 SQLite)，零服务器成本，全球边缘 CDN 毫秒级响应。
  - **私有服务器 / Docker 模式**：支持独立 VPS、Docker Compose 或局域网内网单机运行，数据完全自主掌控。
- **链接与分类管理**：支持多分类过滤、自定义排序、置顶固定，以及 Favicon 自动探测与多源回退。
- **多行垂直公告栏**：3 行平滑分页轮播，支持 Markdown 详情弹窗与独立页面。
- **网络工具集成**：提供 Favicon 快速提取与实时 HTTP 连通性探测 (Ping)。
- **数据迁移与备份**：
  - 支持全量 JSON 备份导入与导出。
  - 原生兼容 Chrome / Edge / Firefox 导出的 HTML 书签文件，支持一键分类解析导入。
- **无状态安全鉴权**：基于 Web Crypto HMAC-SHA256 签名机制，支持跨节点分布式验证。

---

## 架构示意

### 1. Cloudflare Serverless 模式
```
[ 用户浏览器 ]
      │
[ Cloudflare 边缘网络 ]
   ├── [ Cloudflare Pages (Vue 3 静态前端) ]
   └── [ Pages Functions (Hono API - /api/*) ]
            │
       [ Cloudflare D1 (分布式 SQLite) ]
```

### 2. 传统私有化 / Docker 模式
```
[ 用户浏览器 ] ──> [ Nginx (端口 80/443) ] ──> [ 前端静态资源 /dist ]
                                  │ (反向代理 /api)
                                  ▼
                         [ Go 后端服务 (Gin) ] ──> [ 本地 nav.db (SQLite) ]
```

---

## 部署方案

### 方案 A：Cloudflare Serverless 部署（推荐，0 成本）

#### 1. 创建 D1 数据库
1. 登录 Cloudflare 控制台，进入 **存储和数据库** -> **D1 SQL 数据库**，创建名为 `minimal-nav-db` 的数据库。
2. 进入控制台 (Console)，复制项目中的 `d1/schema.sql` 执行，完成数据表与初始数据初始化。

#### 2. 部署 Pages
1. 在控制台进入 **Workers 和 Pages** -> **创建应用程序** -> **Pages** -> **连接到 Git**。
2. 选择本仓库（分支选择 `cloudflare-serverless`），构建参数设置：
   - 框架预设：`Vite`
   - 根目录：`frontend`
   - 构建命令：`npm run build`
   - 输出目录：`dist`
3. 保存并部署。

#### 3. 绑定 D1 与密码配置
1. 进入该 Pages 项目的 **设置** -> **函数 (Functions)**。
2. 添加 **D1 数据库绑定**：变量名设为 `DB`（大写），选择 `minimal-nav-db`。
3. 添加 **环境变量**：变量名设为 `ADMIN_PASSWORD`，填入管理密码（默认 `admin123`）。
4. 重新部署即可上线。

---

### 方案 B：本地 / 私有服务器 Docker 部署（内网 / 独立 VPS）

在主分支（`master`）或自建环境中直接通过 Docker Compose 一键启动：

```bash
# 启动前端与 Go 后端服务
docker compose up -d
```

默认服务端口：
- 前端与管理界面：`http://localhost:80`
- 数据持久化路径：`./data/nav.db`

---

## 本地开发与预览

```bash
# 1. 安装依赖
npm install
cd frontend && npm install && cd ..

# 2. 初始化本地仿真 D1 数据库
npm run d1:init:local

# 3. 编译并启动本地 Pages + D1 仿真服务
npm run build:frontend
npm run preview
```
本地访问地址：`http://localhost:8788`。

---

## 目录结构

```
minimal-nav/
├── d1/schema.sql        # Cloudflare D1 数据表结构与种子数据
├── functions/api/       # Cloudflare Pages Functions (Hono API)
├── frontend/            # Vue 3 前端源码
├── CLOUDFLARE_DEPLOY.md # Cloudflare 专属部署手册
├── wrangler.toml        # Cloudflare 配置文件
└── package.json         # 根工程脚本配置
```

---

## 开源协议

[MIT License](LICENSE)
