# Minimal Nav

极简冷峻风格的个人与团队导航系统。支持 **Cloudflare Serverless 全托管** 与 **本地 / 私有服务器单二进制一键部署** 两种运行形态，满足从零成本云端托管到企业内网私有化的全部场景。

---

## 核心特性

- **冷峻极简排版**：采用瑞士平面设计风格（Swiss Line Layout），以极细几何线条与严谨无衬线字阶构建界面，克制无冗余修饰。
- **双模部署支持**：
  - **Cloudflare Serverless 模式**：基于 Pages + Functions (Hono) + D1 (分布式 SQLite)，零服务器成本，全球边缘 CDN 毫秒级响应。
  - **单二进制 0 依赖自建模式**：Go 内嵌前端编译，单个二进制文件（约 20MB）双击即运行，内置 SQLite 数据库，无需安装 Node.js / Nginx / Docker。
  - **轻量 Docker 模式**：提供一体化多阶段极简镜像，单容器体积仅 ~20MB。
- **链接与分类管理**：支持多分类过滤、自定义排序、置顶固定，以及 Favicon 自动探测与多源回退。
- **多行垂直公告栏**：3 行平滑分页轮播，支持 Markdown 详情弹窗与独立页面。
- **网络工具集成**：提供 Favicon 快速提取与实时 HTTP 连通性探测 (Ping)。
- **数据迁移与备份**：
  - 支持全量 JSON 备份导入与导出。
  - 原生兼容 Chrome / Edge / Firefox 导出的 HTML 书签文件，支持一键分类解析导入。
- **系统设置与备案号**：支持后台自定义站点名称、副标题描述及工信部 ICP 备案号直达跳转。

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

### 2. 本地 / 私有服务器单二进制运行模式
```
                    ┌──────────────────────────────────────────────┐
                    │          单一二进制可执行文件 minimal-nav        │
                    │                                              │
  用户浏览器  ──────> │  ├── 内嵌前端资源 (Vue 3 SPA)                │
 (访问 :8080)        │  ├── 后端 API 引擎 (Gin 路由)                 │
                    │  └── 纯 Go SQLite 驱动 (无 CGO 依赖)          │
                    └──────────────────────┬───────────────────────┘
                                           │ 自动读写
                                           ▼
                                本地 SQLite (nav.db)
```

---

## 部署方案

### 方案 A：单二进制 0 依赖一键运行（本地 / 独立 VPS 推荐）

直接下载或编译单个可执行文件，即可在任意机器上单文件运行：

```bash
# 1. 一键打包编译单二进制文件 (Windows 生成 minimal-nav.exe，Linux 生成 minimal-nav)
go run scripts/build.go

# 2. 运行
./minimal-nav
```
启动后自动在同级目录下读写 `nav.db`，浏览器打开 `http://127.0.0.1:8080` 即可使用。

---

### 方案 B：Cloudflare Serverless 部署（云端全托管，0 成本）

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

### 方案 C：Docker 单容器部署

```bash
# 启动轻量单容器
docker compose up -d
```
访问地址：`http://localhost:8080`。

---

## 目录结构

```
minimal-nav/
├── backend/             # Go 后端源码与内嵌静态资源服务
├── frontend/            # Vue 3 极简前端源码
├── d1/schema.sql        # Cloudflare D1 数据表结构与种子数据
├── scripts/build.go     # 单二进制一体化打包构建脚本
├── Dockerfile           # 极简单镜像多阶段构建文件
├── docker-compose.yml   # 极简单容器配置
└── package.json
```

---

## 开源协议

[MIT License](LICENSE)
