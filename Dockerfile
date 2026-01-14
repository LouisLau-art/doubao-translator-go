# 使用 Go 官方镜像作为构建环境
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译生产版本（静态编译）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o translator

# ------------------------------
# 使用极简镜像作为运行环境
# ------------------------------
FROM alpine:3.18

# 设置环境变量
ENV PORT=5000
ENV GIN_MODE=release

# 添加证书和时区支持
RUN apk add --no-cache ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/translator .
COPY --from=builder /app/static ./static
COPY --from=builder /app/.env.example ./

# 创建非特权用户
RUN adduser -D translator
RUN chown translator:translator /app /app/translator /app/static
USER translator

# 暴露端口
EXPOSE 5000

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:$PORT/api/health || exit 1

# 启动命令
CMD ["./translator"]