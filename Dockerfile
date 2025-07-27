# ビルドステージ
FROM golang:1.24-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git gcc musl-dev libwebp-dev

# 作業ディレクトリを設定
WORKDIR /app

# go.modとgo.sumをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/server

# 実行ステージ
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk --no-cache add ca-certificates libwebp

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# 静的ファイルをコピー
COPY --from=builder /app/web ./web

# 必要なディレクトリを作成
RUN mkdir -p downloads

# ポートを公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./main"] 