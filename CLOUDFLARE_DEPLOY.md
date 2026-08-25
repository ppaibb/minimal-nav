# Cloudflare 部署指南

本项目支持通过 Cloudflare Pages 与 D1 数据库实现零服务器全托管部署。

---

## 方式一：Cloudflare 控制台部署（推荐）

### 1. 创建 D1 数据库
1. 登录 [Cloudflare 控制台](https://dash.cloudflare.com/)。
2. 进入 **存储和数据库** -> **D1 SQL 数据库**，点击 **创建数据库**，输入名称 `minimal-nav-db`。
3. 进入该数据库的 **控制台 (Console)** 标签页，复制项目中的 `d1/schema.sql` 内容粘贴并执行。

### 2. 部署 Pages 项目
1. 进入 **Workers 和 Pages** -> **创建应用程序** -> **Pages** -> **连接到 Git**。
2. 选择本仓库（分支选择 `cloudflare-serverless`），填写构建配置：
   - 框架预设: `Vite`
   - 根目录: `frontend`
   - 构建命令: `npm run build`
   - 输出目录: `dist`
3. 点击 **保存并部署**。

### 3. 配置数据库绑定与环境变量
1. 部署完成后，进入该 Pages 项目的 **设置** -> **函数 (Functions)**。
2. 在 **D1 数据库绑定** 中点击 **添加绑定**：
   - 变量名称: `DB` (必须大写)
   - D1 数据库: 选择 `minimal-nav-db`
3. 在 **环境变量** 中点击 **添加变量**：
   - 变量名称: `ADMIN_PASSWORD`
   - 值: 设置自定义管理员密码（若留空则默认为 `admin123`）
4. 进入 **部署 (Deployments)** 页面，点击最新部署项右侧菜单选择 **重新部署 (Retry deployment)** 即可生效。

---

## 方式二：Wrangler 命令行部署

```bash
# 1. 登录 Cloudflare
npx wrangler login

# 2. 创建 D1 数据库
npx wrangler d1 create minimal-nav-db
# 将生成的 database_id 填入 wrangler.toml 中

# 3. 初始化远端数据表
npx wrangler d1 execute minimal-nav-db --remote --file=./d1/schema.sql

# 4. 构建前端并发布
cd frontend && npm install && npm run build && cd ..
npx wrangler pages deploy frontend/dist --project-name=minimal-nav
```

---

## 本地开发与测试

```bash
# 1. 安装依赖
npm install
cd frontend && npm install && cd ..

# 2. 初始化本地 D1 数据库
npm run d1:init:local

# 3. 构建并启动本地仿真
npm run build:frontend
npm run preview
```
本地访问地址：`http://localhost:8788`。

---

## 常见问题

- **默认管理员密码**：默认为 `admin123`，可在 Pages 环境变量中通过 `ADMIN_PASSWORD` 自定义。
- **数据备份**：管理后台支持导出全量 JSON 数据及导入浏览器 HTML 书签文件。
