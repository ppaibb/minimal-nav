# 阶段 1: 前端静态资源构建
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# 阶段 2: Go 单二进制编译 (内嵌前端静态资源)
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
COPY --from=frontend-builder /app/frontend/dist ./dist

RUN go build -ldflags="-s -w" -o minimal-nav .

# 阶段 3: 极简纯净运行镜像 (体积仅 ~20MB)
FROM alpine:latest
WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata
COPY --from=backend-builder /app/minimal-nav /app/minimal-nav

RUN mkdir -p /app/data
ENV DB_PATH=/app/data/nav.db
ENV PORT=8080

EXPOSE 8080

CMD ["/app/minimal-nav"]
