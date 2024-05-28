# 使用 alpine 作为基础镜像
FROM alpine:latest

# 设置 Go 的版本
ENV GO_VERSION=1.22.3

# 安装必要的依赖和工具
RUN apk add --no-cache \
    curl \
    tar \
    bash \
    gcc \
    musl-dev

# 下载并安装 Go
RUN curl -o go${GO_VERSION}.linux-amd64.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# 设置 Go 环境变量
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# 创建工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制应用源代码
COPY . .


RUN echo "build complete"