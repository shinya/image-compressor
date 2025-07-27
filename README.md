# Image Compression Tool

A web application that converts and compresses PNG, JPG, JPEG, and GIF images to WebP format.

[日本語版はこちら](README_ja.md)

## Features

- **Supported Formats**: PNG, JPG, JPEG, GIF → WebP
- **Compression Settings**: Adjustable quality (1-100), width, and height
- **Drag & Drop**: Intuitive file upload interface
- **Real-time Preview**: Preview of uploaded images
- **Compression Results**: Display original size, compressed size, and compression ratio
- **Download Function**: Download compressed images

## Technology Stack

- **Backend**: Go + Gin
- **Image Processing**: imaging + webp
- **Frontend**: HTML + CSS + JavaScript
- **Container**: Docker + Docker Compose

## Setup

### Prerequisites

- Docker
- Docker Compose

### Quick Start

1. Clone the repository
```bash
git clone <repository-url>
cd image-compressor
```

2. Start with Docker Compose
```bash
docker-compose up --build
```

3. Access in your browser
```
http://localhost:8080
```

## Usage

1. **Image Upload**
   - Drag & drop or click to select an image
   - Supported formats: PNG, JPG, JPEG, GIF
   - Maximum file size: 10MB

2. **Compression Settings** (Optional)
   - **Quality**: Set between 1-100 (default: 80)
   - **Width**: Specify in pixels (leave blank for auto-adjustment)
   - **Height**: Specify in pixels (leave blank for auto-adjustment)

3. **Start Compression**
   - Click "Start Compression" button
   - Loading indicator will be displayed during processing

4. **View Results & Download**
   - Compression results will be displayed
   - Click "Download" button to download the compressed image

## API Specification

### Image Compression API

**Endpoint**: `POST /api/compress`

**Request**:
- `image`: Image file (multipart/form-data)
- `quality`: Quality (1-100, optional)
- `width`: Width in pixels (optional)
- `height`: Height in pixels (optional)

**Response**:
```json
{
  "success": true,
  "data": {
    "original_filename": "example.png",
    "compressed_filename": "example_compressed.webp",
    "original_size": 1024000,
    "compressed_size": 256000,
    "compression_ratio": 25.0,
    "processing_time": 1.23,
    "download_url": "/api/download/example_compressed.webp"
  }
}
```

### Download API

**Endpoint**: `GET /api/download/:filename`

**Response**: Compressed image file

## Environment Variables

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `PORT` | `8080` | Server port number |
| `DOWNLOAD_DIR` | `./downloads` | Directory for compressed files |

## Development

### Local Development Environment

1. Install Go
2. Install dependencies
```bash
go mod tidy
```

3. Start the server
```bash
go run cmd/server/main.go
```

### Dependencies

```bash
go get github.com/gin-gonic/gin
go get github.com/chai2010/webp
go get github.com/disintegration/imaging
```

## Releases

The latest release can be found on the [GitHub releases page](https://github.com/shinya/image-compressor/releases).

### Installation

Download the latest release for your platform from the [releases page](https://github.com/shinya/image-compressor/releases).

## Contributing

Pull requests and issue reports are welcome.
