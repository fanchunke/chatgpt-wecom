# Step 1: Modules caching
FROM golang:1.19 AS base
ENV GOPROXY=https://goproxy.cn

# Move to working directory /build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# Step 2: Builder
FROM golang:1.19 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=1

COPY --from=base /go/pkg /go/pkg
RUN sed -i "s@http://deb.debian.org@http://mirrors.aliyun.com@g" /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends gcc-x86-64-linux-gnu libc6-dev-amd64-cross
COPY . /app
WORKDIR /app
RUN CC=x86_64-linux-gnu-gcc go build -ldflags '-linkmode external -extldflags "-static"' -o /bin/chatgpt-wecom ./cmd/app

# Step 3: Final
FROM alpine:latest
WORKDIR /home/works/program
COPY ./conf/chatgpt.conf ./chatgpt.conf
COPY --from=builder /bin/chatgpt-wecom .

EXPOSE 8000
CMD ["./chatgpt-wecom", "-conf=chatgpt.conf"]