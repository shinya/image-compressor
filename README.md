# Image Compressor

A high-performance image compression service that converts PNG, JPG (JPEG), and GIF images to WebP format with optional compression rate and resolution parameters.

[日本語版はこちら](README_ja.md)

## Features

- **Image Format Support**: PNG, JPG (JPEG), GIF → WebP conversion
- **Compression Options**: Configurable quality, width, and height parameters
- **Web Interface**: Simple and intuitive single-page HTML interface
- **API Endpoints**: RESTful API for programmatic access
- **High Performance**: Built with Go for optimal performance
- **Docker Support**: Easy deployment with Docker and Docker Compose

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/shinya/image-compressor.git
cd image-compressor

# Start the service
docker-compose up -d

# Access the web interface
open http://localhost:8080
```

### Manual Installation

```bash
# Prerequisites
sudo apt-get update
sudo apt-get install -y gcc libc6-dev libwebp-dev

# Build the application
go build -o image-compressor ./cmd/server

# Run the application
./image-compressor
```

## API Usage

### Compress Image

**Endpoint**: `POST /api/compress`

**Form Data**:
- `image`: Image file (PNG, JPG, JPEG, GIF)
- `quality`: Compression quality (0-100, optional)
- `width`: Target width (optional)
- `height`: Target height (optional)

**Response**:
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

### Download Compressed File

**Endpoint**: `GET /api/download/{filename}`

Returns the compressed WebP file.

## Configuration

Environment variables can be used to configure the application:

```bash
PORT=8080                    # Server port (default: 8080)
DOWNLOAD_DIR=./downloads     # Download directory (default: ./downloads)
MAX_FILE_SIZE=10485760       # Maximum file size in bytes (default: 10MB)
DEFAULT_QUALITY=80          # Default compression quality (default: 80)
DEFAULT_WIDTH=1920          # Default width (default: 1920)
DEFAULT_HEIGHT=1080         # Default height (default: 1080)
```

## Development

### Prerequisites

- Go 1.24 or later
- GCC and development libraries
- libwebp-dev

### Local Development

```bash
# Install dependencies
go mod download

# Run the application
go run ./cmd/server

# Run tests
go test ./...

# Build for development
go build -o image-compressor ./cmd/server
```

### Docker Development

```bash
# Build the development image
docker build -t image-compressor-dev .

# Run with volume mount for development
docker run -p 8080:8080 -v $(pwd):/app image-compressor-dev
```

## Project Structure

```
image-compressor/
├── cmd/server/              # Main application entry point
│   └── main.go             # Server startup and configuration
├── internal/
│   ├── api/                # HTTP handlers and routing
│   │   └── handlers.go     # API endpoints implementation
│   ├── config/             # Configuration management
│   │   └── config.go       # Environment and app configuration
│   └── service/            # Business logic and image processing
│       └── image_service.go # Image compression logic
├── web/static/             # Frontend files
│   ├── index.html          # Main HTML page
│   ├── style.css           # CSS styles
│   └── app.js              # JavaScript functionality
├── downloads/              # Compressed file storage
├── scripts/                # Utility scripts
│   └── pre-build.sh       # Pre-build validation script
├── .github/workflows/      # GitHub Actions
│   └── test.yml           # Automated testing workflow
├── Dockerfile              # Production Docker image
├── docker-compose.yml      # Development environment
├── go.mod                  # Go module dependencies
├── go.sum                  # Go module checksums
├── .gitignore              # Git ignore patterns
├── README.md               # English documentation
└── README_ja.md           # Japanese documentation
```

## Technology Stack

- **Backend**: Go 1.24 + Gin web framework
- **Image Processing**: 
  - `github.com/chai2010/webp` for WebP encoding
  - `github.com/disintegration/imaging` for image manipulation
- **Frontend**: HTML5 + CSS3 + Vanilla JavaScript
- **Container**: Docker + Docker Compose
- **Testing**: GitHub Actions for automated testing

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

