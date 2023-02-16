# Step 1: Modules caching
FROM golang:1.19 AS base

# ENV GOPROXY=https://goproxy.cn

# Move to working directory /build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# Step 2: Builder
FROM golang:1.18 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
COPY --from=base /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN go build -o /bin/chatgpt-wecom ./cmd/app

# Step 3: Final
FROM alpine:latest
WORKDIR /home/works/program
COPY ./conf/chatgpt.conf ./chatgpt.conf
COPY --from=builder /bin/chatgpt-wecom .

EXPOSE 8000
CMD ["./chatgpt-wecom"]
