# DC_Handler Phase 1 - Go Backend

A Go-based content handler service for WordPress automation with Ollama integration (local LLM).

## Quick Start

### Prerequisites
- **Go 1.21+** - [Download](https://golang.org/dl)
- **Docker & Docker Compose** (optional but recommended)
- **Ollama** - [Download](https://ollama.com)
- WordPress instance with REST API enabled
- Pixabay API key - [Get free key](https://pixabay.com/api/)

### Setup Steps

1. **Copy environment template:**
```bash
cp .env.example .env
```

2. **Edit .env with your credentials:**
```bash
WORDPRESS_URL=https://your-wordpress-site.com
WORDPRESS_USERNAME=your_username
WORDPRESS_PASSWORD=your_password
PIXABAY_API_KEY=your_key
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=mistral
```

3. **Run setup or installation:**
```bash
go mod tidy
go run main.go
```

## Running the Application

### Option 1: Direct with Go
```bash
go run main.go
```
Server runs at `http://localhost:8080`

### Option 2: Docker Compose (Development)
```bash
docker-compose up
```

### Option 3: Docker Production Build
```bash
docker build -t dc-handler:latest --target production .
docker run -p 8080:8080 --env-file .env dc-handler:latest
```

## Project Structure

```
.
├── main.go                    # Entry point
├── go.mod / go.sum           # Dependencies
├── docker-compose.yml        # Development container
├── Dockerfile                # Multi-stage build
├── .env.example              # Environment template
│
├── config/
│   └── config.go             # Configuration loading
├── handlers/
│   ├── wp_posts.go           # GET /api/wp-posts
│   ├── optimize_content.go   # POST /api/optimize-content
│   └── publish_post.go       # POST /api/publish-post
├── services/
│   ├── wordpress.go          # WordPress API
│   ├── pixabay.go            # Image search
│   └── ollama.go             # Local LLM
└── utils/
    └── auth.go               # Auth helpers
```

## API Endpoints

### GET /api/wp-posts
Fetch WordPress posts

**Query:** `?status=draft&per_page=100`

**Response:**
```json
{
  "success": true,
  "posts": [{"id": 123, "title": "Post", "content": "..."}]
}
```

### POST /api/optimize-content
Optimize content using Ollama

**Body:**
```json
{
  "title": "Title",
  "content": "Content",
  "excerpt": "Excerpt"
}
```

**Response:**
```json
{
  "optimizedTitle": "Optimized Title",
  "optimizedContent": "Optimized content from Ollama",
  "suggestedImage": "https://pixabay.com/image.jpg",
  "imageSource": "Pixabay"
}
```

### POST /api/publish-post
Publish post to WordPress

**Body:**
```json
{
  "postId": 123,
  "title": "Title",
  "content": "Content",
  "featuredImageUrl": "https://example.com/image.jpg"
}
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `ENV` | development | Environment mode |
| `WORDPRESS_URL` | - | WordPress URL |
| `WORDPRESS_USERNAME` | - | WordPress username |
| `WORDPRESS_PASSWORD` | - | WordPress password |
| `OLLAMA_URL` | http://localhost:11434 | Ollama server |
| `OLLAMA_MODEL` | mistral | Model to use |
| `PIXABAY_API_KEY` | - | Pixabay API key |

## Development Workflow

```bash
# 1. Setup
cp .env.example .env
# Edit .env with credentials

# 2. Download dependencies
go mod download

# 3. Start Ollama (in another terminal)
ollama serve

# 4. Run application
go run main.go

# 5. Test endpoints
curl http://localhost:8080/api/wp-posts
```

## Health Check

```bash
curl http://localhost:8080/health
```

## License

MIT