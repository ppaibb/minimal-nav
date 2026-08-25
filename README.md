# 🧭 Minimal Nav · 极简冷峻云原生导航系统

<p align="center">
  <img src="https://img.shields.io/badge/Architecture-Cloudflare_Serverless-F38020?style=flat-square&logo=cloudflare" alt="Cloudflare" />
  <img src="https://img.shields.io/badge/Frontend-Vue_3_+_Vite-4FC08D?style=flat-square&logo=vue.js" alt="Vue 3" />
  <img src="https://img.shields.io/badge/Backend-Hono_API-E36002?style=flat-square&logo=hono" alt="Hono" />
  <img src="https://img.shields.io/badge/Database-Cloudflare_D1_(SQLite)-0051C3?style=flat-square&logo=sqlite" alt="D1" />
  <img src="https://img.shields.io/badge/Cost-100%25_FREE-success?style=flat-square" alt="Free" />
</p>

> 专为极客与团队打造的**极简冷峻风格（Swiss & Editorial Layout）导航系统**。  
> **100% Serverless 架构**：无需购买服务器、无需配置 Nginx / Docker，免费托管于 Cloudflare 全球边缘网络，**终身 0 成本、毫秒级响应、永久可用**。

---

## ✨ 核心特性

- 🏛️ **冷峻建筑与出版物设计美学**：
  - 严格遵循瑞士平面排版（Swiss Line Layout），使用极细直线几何分割与精致无衬线字阶。
  - **坚决拒绝**过度包装的“药丸胶囊”与“浮夸大阴影”，保持高级、干练与克制。
- ⚡ **100% Cloudflare Serverless 全球边缘生态**：
  - **前端托管**：基于 Cloudflare Pages，全球 Anycast CDN 极速分发。
  - **后端计算**：基于 TypeScript + Hono 框架驱动的 Pages Functions，毫秒级冷启动。
  - **云原生存储**：基于 Cloudflare D1 分布式 SQLite 数据库，读写性能极致。
- 🔗 **全功能导航链接管理**：
  - 支持多分类过滤、自定义排序、置顶固定。
  - 智能 Favicon 探测：自动抓取 Google 64px 高清站点图标与降级回退。
- 📢 **纯净多行垂直公告广播**：
  - 3 行垂直平滑过渡翻页轮播，支持 Markdown 富文本详情弹窗/独立页面。
- 🛠️ **内置边缘网络工具**：
  - **Favicon 解析**：输入 URL 自动获取高清图标。
  - **实时 Ping 测速**：基于边缘网络实时探测目标网站的 HTTP 状态码与响应延迟（ms）。
- 🔒 **无状态安全鉴权**：
  - 基于 Web Crypto HMAC-SHA256 签名 Token，无需内存 Session，完美适配多节点边缘计算。
- 💾 **数据备份与无缝迁移**：
  - 支持一键导出/恢复完整 JSON 数据。
  - **原生兼容 Chrome / Edge / Firefox 书签**：上传浏览器导出的 HTML 书签文件即可自动解析分类并批量导入。

---

## 🏗️ 架构概览

```
                             ┌────────────────────────┐
                             │       用户浏览器       │
                             └───────────┬────────────┘
                                         │
                        ┌────────────────┴────────────────┐
                        │    Cloudflare Global Edge CDN   │
                        ▼                                 ▼
            ┌───────────────────────┐         ┌────────────────────────┐
            │   Cloudflare Pages    │         │   Cloudflare Functions │
            │ (Vue 3 静态前端资源)  │         │ (Hono API - /api/*)    │
            └───────────────────────┘         └───────────┬────────────┘
                                                          │
                                                          ▼
                                              ┌────────────────────────┐
                                              │  Cloudflare D1 (SQLite)│
                                              │ (Links/Announce/Config)│
                                              └────────────────────────┘
```

---

## 🚀 5 分钟极速部署到 Cloudflare（完全免费）

### 1. 创建 Cloudflare D1 数据库
1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/)；
2. 进入 **存储和数据库 (Storage & Databases)** -> **D1 SQL 数据库** -> **创建数据库**，名称填写：`minimal-nav-db`；
3. 进入该数据库的 **控制台 (Console)** 标签页，将本项目 [`d1/schema.sql`](d1/schema.sql) 的全部内容粘贴进去并点击 **执行 (Execute)**。

### 2. 部署 Cloudflare Pages
1. 在 Cloudflare 控制台进入 **计算 (Workers) 和 Pages** -> **创建应用程序** -> **Pages** -> **连接到 Git**；
2. 选择本仓库，配置构建参数：
   - **框架预设 (Framework preset)**: `Vite`
   - **根目录 (Root directory)**: `frontend`
   - **构建命令 (Build command)**: `npm run build`
   - **输出目录 (Build output directory)**: `dist`
3. 点击 **保存并部署**。

### 3. 绑定 D1 数据库与设置管理员密码
1. 部署完成后，在 Pages 项目中进入 **设置 (Settings)** -> **Functions (函数)**；
2. 在 **D1 数据库绑定** 点击 **添加绑定**：
   - **变量名称 (Variable name)**: `DB` （必须全部大写）
   - **D1 数据库**: 选择刚创建的 `minimal-nav-db`
3. 在 **环境变量 (Environment variables)** 点击 **添加变量**：
   - **变量名称**: `ADMIN_PASSWORD`
   - **值**: 设置你的自定义管理密码（若不填，默认为 `admin123`）
4. 返回 **部署 (Deployments)** 页面点击 **创建重新部署** 即可完整上线！🎉

---

## 💻 本地开发与离线模拟

本项目支持纯本地离线开发，包含完整的 Cloudflare D1 数据库仿真器：

```bash
# 1. 克隆仓库并安装依赖
git clone https://github.com/your-username/minimal-nav.git
cd minimal-nav
npm install
cd frontend && npm install && cd ..

# 2. 初始化本地 D1 数据库
npm run d1:init:local

# 3. 构建前端并在本地启动 Pages + D1 仿真服务
npm run build:frontend
npm run preview
```
本地访问：`http://localhost:8788`。

---

## 📂 项目结构

```
minimal-nav/
├── d1/
│   └── schema.sql        # Cloudflare D1 表结构定义与初始化种子数据
├── functions/
│   └── api/
│       └── [[route]].ts  # Hono 驱动的 Cloudflare Pages Functions 边缘 API 全量实现
├── frontend/             # Vue 3 + Vite + Tailwind CSS 纯静态前端
│   ├── src/
│   │   ├── pages/        # 首页 Index / 公告详情 Docs / 管理后台 Admin
│   │   └── utils/        # 站点配置与网络请求
│   └── package.json
├── CLOUDFLARE_DEPLOY.md  # 详尽的 Cloudflare 图形化 & CLI 部署手册
├── wrangler.toml         # Cloudflare Wrangler 配置文件
└── package.json          # 根工程脚本配置
```

---

## 📄 开源协议

本项目基于 [MIT License](LICENSE) 开源。欢迎 Star ⭐️ 与 PR！
