# 構建階段
FROM golang:1.21 AS builder

WORKDIR /app

# 複製 go.mod 和 go.sum 文件（如果存在）
COPY go.mod ./
COPY go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 更新依賴
RUN go mod tidy

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 運行階段
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 從構建階段複製編譯好的執行檔
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 運行
CMD ["./main"]