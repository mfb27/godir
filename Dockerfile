# 构建阶段
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o bin/app main.go

# 运行阶段
FROM alpine:latest

# 安装ca证书以便HTTPS请求
RUN apk --no-cache add ca-certificates

# 创建用户
RUN adduser -D -s /bin/sh godir

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/app .
COPY --from=builder /app/config ./config

# 更改所有权
RUN chown -R godir:godir ./

# 使用非root用户运行
USER godir

# 暴露端口
EXPOSE 8080

# 运行应用
ENTRYPOINT ["./app", "-c", "config/dev.yml"]