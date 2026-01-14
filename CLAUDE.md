# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Go + Vue Petite Doubao Translator** - a minimal, high-performance web translator using the Volcano Engine Doubao Translation API (ARK platform). It features real-time translation with Markdown/LaTeX rendering and local history management.

## Architecture

### Backend (Go + Gin)
- **Entry point**: `main.go` - Single-file Go application
- **Framework**: Gin 1.11.0 with CORS middleware
- **API integration**: Volcano Engine ARK Doubao Translation API v3
- **Key features**:
  - MD5-based caching with configurable TTL (default: 1 hour)
  - Token bucket rate limiting (30 requests per 2 seconds)
  - Automatic text chunking for documents >800 characters
  - Health check endpoint at `/api/health`
  - Static file serving from `/static` directory

### Frontend (Vue Petite)
- **Framework**: Vue Petite 0.4.1 (lightweight Vue alternative)
- **Entry point**: `static/index.html`
- **Key components**:
  - `static/app.js`: Vue Petite application logic
  - `static/style.css`: Complete styling with dark theme
  - `static/libs/`: Frontend dependencies (petite-vue.js, marked.min.js, MathJax)

### Data Flow
```
User Input → Vue App → POST /api/translate →
Go Server (validate, chunk, cache, rate limit) →
Doubao API → Response → Go Server → JSON →
Vue App (render, cache, display)
```

## Development Commands

### Build & Run
```bash
make dev              # Install dependencies + run dev server (go run .)
make setup            # Install dependencies + download frontend libs
make build            # Compile regular binary
make build-prod       # Compile production binary (CGO_ENABLED=0, GOOS=linux)
make serve            # Build production + run binary
```

### Code Quality
```bash
make fmt              # Run go fmt ./...
make lint             # Run golangci-lint (requires installation)
```

### System Service (Linux)
```bash
make install-service  # Install as systemd service (requires sudo)
make status           # Check systemd service status
make logs             # View journal logs
make uninstall        # Remove systemd service
```

### Production
```bash
make prod             # Build production binary + optionally compress with UPX
```

## Configuration

### Environment Variables (.env)
Required variables:
```env
ARK_API_KEY=your_ark_api_key_here        # Required - API key from Volcano Engine
ARK_API_URL=https://ark.cn-beijing.volces.com/api/v3/responses  # API endpoint
PORT=5000                                # Server port (default: 5000)
```

Optional variables:
```env
GIN_MODE=release                         # Gin mode: debug/release
CACHE_TTL=3600                          # Cache time-to-live in seconds
CACHE_MAX_SIZE=1000                     # Maximum cache entries
MAX_TEXT_LENGTH=5000                    # Text length limit per request
RATE_LIMIT_RPM=30                       # API rate limit (requests per minute)
```

### API Endpoints
- `GET /` - Serves the Vue Petite frontend
- `GET /api/languages` - Returns supported language list (14 languages)
- `POST /api/translate` - Main translation endpoint with caching & rate limiting
- `GET /api/health` - Health check endpoint
- `GET /static/*` - Static file serving
- `GET /libs/*` - Frontend libraries (aliased to /static/libs)

## Key Implementation Details

### Backend Caching
- Uses `sync.Map` for thread-safe caching
- Cache keys are MD5 hashes of `text + source + target`
- Background goroutine cleans expired cache entries every hour

### Rate Limiting
- Token bucket algorithm via `golang.org/x/time/rate`
- Default: 30 requests per 2 seconds
- Returns HTTP 429 when limit exceeded

### Text Processing
- Automatically splits text into ~800 character chunks
- Preserves paragraph boundaries during splitting
- Reassembles translated chunks maintaining original structure

### Frontend Features
- Real-time auto-translation with 500ms debounce
- Markdown rendering via marked.js
- LaTeX formula rendering via MathJax 3
- Local history (50 entries) using localStorage
- Dark theme with responsive design
- Font size adjustment (12px-26px)

## Dependencies

### Backend (go.mod)
- `github.com/gin-gonic/gin` - Web framework
- `github.com/gin-contrib/cors` - CORS middleware
- `github.com/joho/godotenv` - Environment variable loading
- `golang.org/x/time/rate` - Rate limiting

### Frontend (downloaded via make setup)
- `petite-vue.js` - Vue Petite framework
- `marked.min.js` - Markdown parser
- `tex-mml-chtml.js` - MathJax 3 for LaTeX rendering

## File Structure
```
go-translator/
├── main.go                      # Go backend entry point
├── go.mod                       # Go module dependencies
├── Makefile                     # Build/run automation
├── .env.example                 # Example environment config
├── download-libs.sh             # Shell script to download frontend libs
└── static/                      # Frontend assets
    ├── index.html               # Main HTML entry point
    ├── app.js                   # Vue Petite application logic
    ├── style.css                # Complete styling (dark/light themes)
    └── libs/                    # Frontend dependencies
        ├── petite-vue.js        # Vue Petite framework
        ├── marked.min.js        # Markdown parser
        └── mathjax/             # MathJax 3 library
            └── tex-mml-chtml.js # LaTeX rendering
```

## Development Notes

1. **API Key Required**: Must obtain ARK API key from Volcano Engine console
2. **Single-file Backend**: All Go logic is in `main.go` for simplicity
3. **Frontend Dependencies**: Downloaded automatically via `make setup`
4. **Production Build**: Uses `CGO_ENABLED=0` for static binary
5. **Caching**: Enabled by default with 1-hour TTL, configurable via env vars
6. **Rate Limiting**: Protects against API abuse, configurable via env vars

## Troubleshooting

- **API Key Errors**: Check `.env` file exists with correct `ARK_API_KEY`
- **Port Conflicts**: Change `PORT` in `.env` if 5000 is occupied
- **Frontend Issues**: Run `make setup` to download missing libraries
- **Rate Limiting**: Adjust `RATE_LIMIT_RPM` in `.env` if needed
- **Cache Issues**: Clear cache by restarting the application