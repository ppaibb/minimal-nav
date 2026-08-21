# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 设置国内/官方 Go 代理加速模块下载
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o minimal-nav-backend .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates 以支持 HTTPS 请求
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/minimal-nav-backend /app/minimal-nav-backend

# 创建数据存储目录
RUN mkdir -p /app/data
ENV DB_PATH=/app/data/nav.db
ENV PORT=8080

EXPOSE 8080

CMD ["/app/minimal-nav-backend"]
