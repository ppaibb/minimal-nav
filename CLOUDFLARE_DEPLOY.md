# ⚡ Minimal Nav - Cloudflare 零成本免服务器部署指南

本项目已全面适配 **Cloudflare Pages + Workers + D1 数据库**。
你无需购买任何云服务器、无需配置 Nginx / Docker，即可拥有全球 Anycast 边缘网络加速与全套导航管理功能，**终身 0 成本（完全免费）**！

---

## 🌟 核心特性

- 🚀 **0 服务器部署**：前端直接托管在 Cloudflare Pages，API 运行在 Cloudflare Workers。
- 🗄️ **云原生 D1 数据库**：全球分布式 SQLite 数据库，读写毫秒级响应。
- 🔒 **安全口令管理**：基于 Web Crypto HMAC-SHA256 的无状态 Token 验证。
- ⚡ **边缘工具集成**：支持 Favicon 智能抓取、URL 连通性测速、Chrome 书签一键导入。
- 🌐 **免费自定义域名 & SSL**：一键绑定你的域名，自动配置 HTTPS。

---

## 🛠️ 方式一：Cloudflare 控制台 5 分钟图形化部署（最推荐）

### 第一步：创建 Cloudflare D1 数据库
1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/)。
2. 在左侧菜单点击 **存储和数据库 (Storage & Databases)** -> **D1 SQL 数据库**。
3. 点击 **创建数据库 (Create Database)**，输入名称：`minimal-nav-db`，点击保存。
4. 进入刚创建的数据库，点击 **控制台 (Console)** 标签页。
5. 打开项目中的 `d1/schema.sql`，复制全部 SQL 内容粘贴到控制台，点击 **执行 (Execute)** 完成表结构与初始数据初始化。

### 第二步：部署 Cloudflare Pages
1. 在左侧菜单点击 **计算 (Workers) 和 Pages** -> **创建应用程序** -> 切换到 **Pages** 标签页。
2. 选择 **连接到 Git (Connect to Git)**，授权并选择你的 `minimal-nav` 仓库（分支选择 `cloudflare-serverless`）。
3. 构建配置填写如下：
   - **框架预设 (Framework preset)**: `Vite`
   - **根目录 (Root directory)**: `frontend`
   - **构建命令 (Build command)**: `npm run build`
   - **构建输出目录 (Build output directory)**: `dist`
4. 点击 **保存并部署 (Save and Deploy)**。

### 第三步：绑定 D1 数据库与环境变量
1. 部署完成后，进入该 Pages 项目的 **设置 (Settings)** -> **Functions (函数)**。
2. 找到 **D1 数据库绑定 (D1 Database Bindings)**，点击 **添加绑定 (Add binding)**：
   - **变量名称 (Variable name)**: `DB` （必须大写）
   - **D1 数据库 (D1 database)**: 选择第一步创建的 `minimal-nav-db`
3. 找到 **环境变量 (Environment variables)**，点击 **添加变量 (Add variable)**：
   - **变量名称**: `ADMIN_PASSWORD`
   - **值**: 设置你的管理员密码（如不设置，默认为 `admin123`）
4. 在 **部署 (Deployments)** 页面点击 **创建重新部署 (Create deployment / Retry deployment)**，即可大功告成！🎉

---

## 💻 方式二：Wrangler 命令行一键部署

如果你习惯本地 CLI 操作，可直接使用 Cloudflare 官方的 `wrangler` 工具：

### 1. 登录 Cloudflare
```bash
npx wrangler login
```

### 2. 创建 D1 数据库并初始化
```bash
# 创建 D1 数据库
npx wrangler d1 create minimal-nav-db

# 将生成的 database_id 填入 wrangler.toml 文件中

# 执行 SQL 初始化表结构与种子数据
npx wrangler d1 execute minimal-nav-db --remote --file=./d1/schema.sql
```

### 3. 构建前端并发布到 Pages
```bash
# 1. 构建前端
cd frontend && npm install && npm run build && cd ..

# 2. 发布到 Cloudflare Pages
npx wrangler pages deploy frontend/dist --project-name=minimal-nav
```

---

## 🧪 本地开发与测试预览

在本地开发时，你无需连接远程 Cloudflare，即可享受完整的 D1 本地模拟器：

```bash
# 1. 安装根目录与前端依赖
npm install
cd frontend && npm install && cd ..

# 2. 初始化本地 D1 数据库
npm run d1:init:local

# 3. 构建前端并在本地启动 Pages + D1 仿真服务
npm run build:frontend
npm run preview
```
打开浏览器访问控制台输出的本地地址（默认 `http://localhost:8788`），即可完整体验所有前后台功能！

---

## ⚙️ 常见问题 (FAQ)

- **Q: 默认后台密码是什么？**
  - A: 默认为 `admin123`。你可以在 Pages 的环境变量中配置 `ADMIN_PASSWORD` 为任意自定义复杂密码。
- **Q: 如何备份和迁移数据？**
  - A: 登录管理后台后，进入系统设置/数据备份即可导出完整 JSON，或直接导入 Chrome/Edge 书签 HTML。也可以随时使用 `wrangler d1 export` 进行数据库备份。
