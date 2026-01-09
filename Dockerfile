# 第一阶段：编译阶段
FROM golang:1.21-alpine AS builder

# 设置国内镜像源加速下载
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 先拷贝依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源代码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main main.go

# 第二阶段：运行阶段
FROM alpine:latest

WORKDIR /root/

# 从编译阶段拷贝二进制文件
COPY --from=builder /app/main .
# 拷贝必要的配置文件和静态资源目录
COPY --from=builder /app/etc ./etc
# 如果有 json 目录或其他必要目录也一并拷贝
COPY --from=builder /app/json ./json

# 暴露配置文件中的端口（假设是 8080）
EXPOSE 8080

# 启动程序
CMD ["./main", "server", "s0"]