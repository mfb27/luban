# 开发指南 (Development Guide)

This guide explains the development workflow for the Luban chat application with proper frontend-backend separation.

## Frontend-Backend Separation

The application is now properly separated into frontend and backend components:

- **Backend**: Go application running on port 8080
- **Frontend**: Static HTML/CSS/JS files served by the backend

## Development Workflow

### 1. Start Backend Services

First, start the required services (MySQL, Redis, MinIO):

```bash
docker-compose up -d
```

### 2. Start Backend Development Server

```bash
go run ./cmd/server.go
```

The backend will start on http://localhost:8080 and automatically serve the frontend from the `./frontend` directory.

### 3. Frontend Development

For frontend development, you have two options:

#### Option A: Using Backend Static File Serving

The backend automatically serves the frontend from `./frontend`. Any changes to HTML, CSS, or JS files will be reflected immediately after refreshing the browser.

#### Option B: Using Frontend Dev Server (Recommended for Frontend-Only Development)

Run the frontend development server separately:

```bash
cd frontend
chmod +x dev-server.sh
./dev-server.sh
```

Access the frontend at http://localhost:3000

**Note**: When using the frontend dev server, you'll need to configure CORS or proxy API requests to the backend (port 8080).

### 4. Building for Production

Before deploying, optimize the frontend assets:

```bash
cd frontend
chmod +x build.sh
./build.sh
```

The build script will:
- Minify CSS and JavaScript (if `minify` is installed)
- Create a version.js file for cache busting
- Optimize assets for production

## API Integration

The frontend communicates with the backend via REST API:

```javascript
// Example API call
fetch('http://localhost:8080/api/chat', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    message: 'Hello',
    model: 'gpt-3.5-turbo'
  })
})
```

### Key API Endpoints

- `POST /api/chat` - Send chat messages
- `GET /api/models` - Get available models
- `GET /api/sessions` - List user sessions
- `POST /api/sessions` - Create new session
- `GET /api/sessions/:id/messages` - Get session messages
- `POST /api/upload` - Upload files
- `POST /api/auth/login` - Login
- `POST /api/auth/register` - Register

## File Structure

```
luban/
├── frontend/              # Frontend source files
│   ├── css/              # Stylesheets
│   ├── js/               # JavaScript files
│   ├── images/           # Image assets
│   ├── index.html        # Main HTML file
│   ├── build.sh          # Production build script
│   └── dev-server.sh     # Development server script
├── internal/              # Backend source code
├── cmd/                  # Application entry points
└── config.yaml          # Configuration file
```

## Hot Reloading

### Backend Hot Reloading

Use tools like `air` for automatic Go code reloading:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run air
air
```

### Frontend Hot Reloading

When using the frontend dev server, it provides hot reloading automatically. When using the backend static file serving, you'll need to manually refresh the browser.

## Testing

### Backend Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific tests
go test ./internal/handler
```

### Frontend Testing

The frontend is plain HTML/CSS/JS and can be tested by:
- Opening the HTML file directly in the browser
- Using browser dev tools for debugging
- Writing unit tests for JavaScript (optional)

## Deployment

### Backend Deployment

1. Build the Go application
```bash
go build -o luban ./cmd
```

2. Copy the frontend to the server
```bash
rsync -av frontend/ /path/to/frontend/
```

3. Run the application
```bash
./luban server --config config.yaml
```

### Frontend Deployment

The frontend doesn't need a separate server. It's served by the Go application from the `static_dir` path.

## Environment Variables

Set these environment variables for development:

```bash
# Backend configuration
export LUBAN_SERVER_HOST=0.0.0.0
export LUBAN_SERVER_PORT=8080
export LUBAN_MYSQL_DSN=root:password@tcp(localhost:3306)/luban
export LUBAN_REDIS_ADDR=localhost:6379
```

## Troubleshooting

### Frontend Not Loading

1. Check that the backend is running on port 8080
2. Verify the `static_dir` in config.yaml points to the correct frontend directory
3. Check browser console for errors

### API Not Working

1. Verify backend is running
2. Check CORS settings if using a separate frontend dev server
3. Verify the API endpoints are correct

### File Upload Issues

1. Check MinIO is running
2. Verify MinIO credentials in config.yaml
3. Check file permissions and size limits