# 构建阶段：使用官方 Go 镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go 模块文件并下载依赖（利用 Docker 缓存机制）
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# 复制项目源代码
COPY . .

# 构建应用（指定 CGO 禁用以保证可移植性）
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener main.go

# 运行阶段：使用轻量的 Alpine 镜像
FROM alpine:3.20

# 安装 CA 证书（用于 HTTPS 请求）
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/url-shortener .
# 复制前端静态文件
COPY --from=builder /app/web ./web

# 暴露应用端口（与代码中保持一致）
EXPOSE 9808

# 启动命令
CMD ["./url-shortener"]