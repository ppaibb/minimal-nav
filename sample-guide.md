# 团队云原生网关与服务接入规范

欢迎使用内部服务网关！本文档介绍如何通过统一入口接入微服务与部署应用。

## 1. 快速入门

执行以下命令快速拉取基础镜像：

```bash
docker pull registry.internal.net/base/node:20-alpine
docker run -d -p 3000:3000 --name web-service registry.internal.net/base/node:20-alpine
```

## 2. 核心路由配置表格

| 服务名称 | 内部端口 | 公网映射路径 | 健康检查地址 |
| :--- | :--- | :--- | :--- |
| **Auth-Center** | `8081` | `/api/v1/auth` | `/healthz` |
| **Data-Engine** | `8082` | `/api/v1/data` | `/actuator/health` |
| **Search-API** | `8083` | `/api/v1/search` | `/ping` |

> **提示：** 如需开通公网域名解析，请在飞书提交审批工单。
