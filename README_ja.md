# Image Compressor

PNG、JPG（JPEG）、GIF画像をWebP形式に変換する高性能な画像圧縮サービスです。圧縮率と解像度のパラメータをオプションで設定できます。

[English version is here](README.md)

## 機能

- **対応画像形式**: PNG、JPG（JPEG）、GIF → WebP変換
- **圧縮オプション**: 品質、幅、高さのパラメータ設定可能
- **Webインターフェース**: シンプルで直感的な単一ページHTMLインターフェース
- **APIエンドポイント**: プログラムからのアクセス用RESTful API
- **高性能**: Goで構築された最適なパフォーマンス
- **Docker対応**: DockerとDocker Composeによる簡単なデプロイ

## クイックスタート

### Docker Composeを使用（推奨）

```bash
# リポジトリをクローン
git clone https://github.com/shinya/image-compressor.git
cd image-compressor

# サービスを起動
docker-compose up -d

# Webインターフェースにアクセス
open http://localhost:8080
```

### 手動インストール

```bash
# 前提条件
sudo apt-get update
sudo apt-get install -y gcc libc6-dev libwebp-dev

# アプリケーションをビルド
go build -o image-compressor ./cmd/server

# アプリケーションを実行
./image-compressor
```

## API使用方法

### 画像圧縮

**エンドポイント**: `POST /api/compress`

**フォームデータ**:
- `image`: 画像ファイル（PNG、JPG、JPEG、GIF）
- `quality`: 圧縮品質（0-100、オプション）
- `width`: 目標幅（オプション）
- `height`: 目標高さ（オプション）

**レスポンス**:
```json
{
  "success": true,
  "message": "Image compressed successfully",
  "original_size": 1024000,
  "compressed_size": 256000,
  "compression_ratio": 75.0,
  "processing_time": 0.5,
  "output_file": "compressed_1234567890.webp",
  "download_url": "/api/download/compressed_1234567890.webp"
}
```

### 圧縮ファイルのダウンロード

**エンドポイント**: `GET /api/download/{filename}`

圧縮されたWebPファイルを返します。

## 設定

環境変数を使用してアプリケーションを設定できます：

```bash
PORT=8080                    # サーバーポート（デフォルト: 8080）
DOWNLOAD_DIR=./downloads     # ダウンロードディレクトリ（デフォルト: ./downloads）
MAX_FILE_SIZE=10485760       # 最大ファイルサイズ（バイト）（デフォルト: 10MB）
DEFAULT_QUALITY=80          # デフォルト圧縮品質（デフォルト: 80）
DEFAULT_WIDTH=1920          # デフォルト幅（デフォルト: 1920）
DEFAULT_HEIGHT=1080         # デフォルト高さ（デフォルト: 1080）
```

## 開発

### 前提条件

- Go 1.24以降
- GCCと開発ライブラリ
- libwebp-dev

### ローカル開発

```bash
# 依存関係をインストール
go mod download

# アプリケーションを実行
go run ./cmd/server

# テストを実行
go test ./...

# 開発用にビルド
go build -o image-compressor ./cmd/server
```

### Docker開発

```bash
# 開発イメージをビルド
docker build -t image-compressor-dev .

# ボリュームマウントで開発実行
docker run -p 8080:8080 -v $(pwd):/app image-compressor-dev
```

## プロジェクト構造

```
image-compressor/
├── cmd/server/              # メインアプリケーションエントリーポイント
│   └── main.go             # サーバー起動と設定
├── internal/
│   ├── api/                # HTTPハンドラーとルーティング
│   │   └── handlers.go     # APIエンドポイント実装
│   ├── config/             # 設定管理
│   │   └── config.go       # 環境とアプリ設定
│   └── service/            # ビジネスロジックと画像処理
│       └── image_service.go # 画像圧縮ロジック
├── web/static/             # フロントエンドファイル
│   ├── index.html          # メインHTMLページ
│   ├── style.css           # CSSスタイル
│   └── app.js              # JavaScript機能
├── downloads/              # 圧縮ファイルストレージ
├── scripts/                # ユーティリティスクリプト
│   └── pre-build.sh       # ビルド前検証スクリプト
├── .github/workflows/      # GitHub Actions
│   └── test.yml           # 自動テストワークフロー
├── Dockerfile              # 本番Dockerイメージ
├── docker-compose.yml      # 開発環境
├── go.mod                  # Goモジュール依存関係
├── go.sum                  # Goモジュールチェックサム
├── .gitignore              # Git無視パターン
├── README.md               # 英語ドキュメント
└── README_ja.md           # 日本語ドキュメント
```

## 技術スタック

- **バックエンド**: Go 1.24 + Gin Webフレームワーク
- **画像処理**: 
  - `github.com/chai2010/webp` for WebPエンコーディング
  - `github.com/disintegration/imaging` for 画像操作
- **フロントエンド**: HTML5 + CSS3 + バニラJavaScript
- **コンテナ**: Docker + Docker Compose
- **テスト**: GitHub Actions for 自動テスト

## 貢献

貢献を歓迎します！プルリクエストを自由に送信してください。

