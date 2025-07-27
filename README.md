# 画像圧縮ツール

PNG、JPG、JPEG、GIF形式の画像をWebP形式に変換・圧縮するWebアプリケーションです。

## 機能

- **対応形式**: PNG、JPG、JPEG、GIF → WebP
- **圧縮設定**: 品質（1-100）、幅、高さの調整が可能
- **ドラッグ&ドロップ**: 直感的なファイルアップロード
- **リアルタイムプレビュー**: アップロードした画像のプレビュー表示
- **圧縮結果表示**: 元のサイズ、圧縮後サイズ、圧縮率を表示
- **ダウンロード機能**: 圧縮後の画像をダウンロード

## 技術スタック

- **バックエンド**: Go + Gin
- **画像処理**: imaging + webp
- **フロントエンド**: HTML + CSS + JavaScript
- **コンテナ**: Docker + Docker Compose

## セットアップ

### 前提条件

- Docker
- Docker Compose

### 起動方法

1. リポジトリをクローン
```bash
git clone <repository-url>
cd image-compressor
```

2. Docker Composeで起動
```bash
docker-compose up --build
```

3. ブラウザでアクセス
```
http://localhost:8080
```

## 使用方法

1. **画像アップロード**
   - ドラッグ&ドロップまたはクリックして画像を選択
   - 対応形式: PNG、JPG、JPEG、GIF
   - 最大ファイルサイズ: 10MB

2. **圧縮設定**（オプション）
   - **品質**: 1-100の範囲で設定（デフォルト: 80）
   - **幅**: ピクセル単位で指定（空欄で自動調整）
   - **高さ**: ピクセル単位で指定（空欄で自動調整）

3. **圧縮実行**
   - 「圧縮開始」ボタンをクリック
   - 処理中はローディング表示

4. **結果確認・ダウンロード**
   - 圧縮結果が表示される
   - 「ダウンロード」ボタンで圧縮後の画像をダウンロード

## API仕様

### 画像圧縮API

**エンドポイント**: `POST /api/compress`

**リクエスト**:
- `image`: 画像ファイル（multipart/form-data）
- `quality`: 品質（1-100、オプション）
- `width`: 幅（ピクセル、オプション）
- `height`: 高さ（ピクセル、オプション）

**レスポンス**:
```json
{
  "success": true,
  "data": {
    "original_filename": "example.png",
    "compressed_filename": "example_compressed.webp",
    "original_size": 1024000,
    "compressed_size": 256000,
    "compression_ratio": 25.0,
    "download_url": "/api/download/example_compressed.webp"
  }
}
```

### ダウンロードAPI

**エンドポイント**: `GET /api/download/:filename`

**レスポンス**: 圧縮された画像ファイル

## 環境変数

| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| `PORT` | `8080` | サーバーのポート番号 |
| `DOWNLOAD_DIR` | `./downloads` | 圧縮後のファイル保存ディレクトリ |

## 開発

### ローカル開発環境

1. Goをインストール
2. 依存関係をインストール
```bash
go mod tidy
```

3. サーバーを起動
```bash
go run cmd/server/main.go
```

### 依存関係

```bash
go get github.com/gin-gonic/gin
go get github.com/chai2010/webp
go get github.com/disintegration/imaging
```

## ライセンス

MIT License

## 貢献

プルリクエストやイシューの報告を歓迎します。 