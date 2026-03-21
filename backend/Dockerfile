# syntax=docker/dockerfile:1

FROM golang:1.25.6-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
RUN GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrate

#------------------------------------

FROM gcr.io/distroless/base-debian12

WORKDIR /app

# 実行ファイル
COPY --from=builder /app/server ./server
COPY --from=builder /app/migrate ./migrate
# マイグレーションファイルをコピー
COPY --from=builder /app/migrations ./migrations

EXPOSE 8000

ENTRYPOINT ["/app/server"]
