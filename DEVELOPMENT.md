# Development Guide

This guide provides comprehensive information for developers working on S3 Smart Browser.

## Table of Contents

- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Development Environment](#development-environment)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Building](#building)
- [Debugging](#debugging)
- [Contributing](#contributing)

## Getting Started

### Prerequisites

- Go 1.25 or later
- Docker (optional)
- Git
- AWS S3 access (for testing)

### Quick Start

1. **Clone the repository:**
   ```bash
   git clone https://github.com/BlackoJack/s3-smart-browser.git
   cd s3-smart-browser
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment:**
   ```bash
   export AWS_ACCESS_KEY_ID=your_access_key
   export AWS_SECRET_ACCESS_KEY=your_secret_key
   export AWS_BUCKET=your-bucket-name
   export AWS_REGION=us-east-1
   ```

4. **Run the application:**
   ```bash
   go run cmd/server/main.go
   ```

## Project Structure

```
s3-smart-browser/
├── cmd/server/           # Application entry point
│   └── main.go          # Main application file
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   │   └── config.go    # Configuration structs and loading
│   ├── handlers/        # HTTP request handlers
│   │   └── handlers.go  # HTTP handlers implementation
│   ├── s3/             # S3 client implementation
│   │   └── client.go    # S3 operations and client
│   └── types/           # Data types and structures
│       └── types.go     # Type definitions
├── web/                 # Web assets
│   ├── static/         # Static files (CSS, JS)
│   │   ├── css/        # Stylesheets
│   │   └── js/         # JavaScript files
│   └── templates/      # HTML templates
│       └── index.html  # Main HTML template
├── charts/             # Helm charts
│   └── s3-smart-browser/
├── docs/               # Documentation
├── docker-compose.yaml # Docker Compose configuration
├── Dockerfile          # Docker image definition
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── README.md           # Project documentation
```

## Development Environment

### Local Development

1. **Set up Go workspace:**
   ```bash
   export GOPATH=$HOME/go
   export PATH=$PATH:$GOPATH/bin
   ```

2. **Install development tools:**
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   go install github.com/air-verse/air@latest
   ```

3. **Use Air for hot reloading:**
   ```bash
   air
   ```

### Docker Development

1. **Build development image:**
   ```bash
   docker build -t s3-smart-browser-dev .
   ```

2. **Run with volume mounting:**
   ```bash
   docker run -p 8080:8080 \
     -v $(pwd):/app \
     -w /app \
     -e AWS_ACCESS_KEY_ID=your_access_key \
     -e AWS_SECRET_ACCESS_KEY=your_secret_key \
     -e AWS_BUCKET=your-bucket-name \
     s3-smart-browser-dev
   ```

### IDE Setup

#### VS Code

Install recommended extensions:
- Go extension
- Docker extension
- Kubernetes extension

#### GoLand/IntelliJ

Configure Go SDK and enable Go modules support.

## Coding Standards

### Go Code Style

1. **Use `gofmt` for formatting:**
   ```bash
   gofmt -w .
   ```

2. **Follow Go naming conventions:**
   - Use camelCase for variables and functions
   - Use PascalCase for exported types and functions
   - Use descriptive names

3. **Add comments for exported functions:**
   ```go
   // ListDirectory retrieves the contents of a directory from S3
   func (c *Client) ListDirectory(ctx context.Context, path string) (*types.DirectoryListing, error) {
       // implementation
   }
   ```

4. **Handle errors explicitly:**
   ```go
   if err != nil {
       return nil, fmt.Errorf("failed to list directory: %w", err)
   }
   ```

### Frontend Code Style

1. **Use modern JavaScript (ES6+):**
   ```javascript
   // Use const/let instead of var
   const response = await fetch(url);
   
   // Use arrow functions
   const handleClick = () => {
       // implementation
   };
   ```

2. **Follow consistent naming:**
   - Use camelCase for variables and functions
   - Use PascalCase for classes
   - Use descriptive names

3. **Add comments for complex logic:**
   ```javascript
   // Sort files: directories first, then files alphabetically
   const sortedFiles = files.sort((a, b) => {
       if (a.is_directory && !b.is_directory) return -1;
       if (!a.is_directory && b.is_directory) return 1;
       return a.name.localeCompare(b.name);
   });
   ```

## Testing

### Unit Tests

1. **Run all tests:**
   ```bash
   go test ./...
   ```

2. **Run tests with coverage:**
   ```bash
   go test -cover ./...
   ```

3. **Run specific package tests:**
   ```bash
   go test ./internal/s3
   ```

### Integration Tests

1. **Set up test environment:**
   ```bash
   export AWS_ACCESS_KEY_ID=test_access_key
   export AWS_SECRET_ACCESS_KEY=test_secret_key
   export AWS_BUCKET=test-bucket
   ```

2. **Run integration tests:**
   ```bash
   go test -tags=integration ./...
   ```

### Frontend Testing

1. **Test JavaScript functionality:**
   ```bash
   # Open browser developer tools
   # Test file navigation, downloads, etc.
   ```

2. **Test responsive design:**
   - Test on different screen sizes
   - Test on mobile devices
   - Test on different browsers

## Building

### Local Build

1. **Build for current platform:**
   ```bash
   go build -o s3-smart-browser ./cmd/server
   ```

2. **Build for specific platform:**
   ```bash
   GOOS=linux GOARCH=amd64 go build -o s3-smart-browser-linux ./cmd/server
   ```

### Docker Build

1. **Build Docker image:**
   ```bash
   docker build -t s3-smart-browser .
   ```

2. **Build with specific tag:**
   ```bash
   docker build -t s3-smart-browser:v1.0.8 .
   ```

### Release Build

1. **Build release binaries:**
   ```bash
   # Linux
   GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o s3-smart-browser-linux ./cmd/server
   
   # macOS
   GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o s3-smart-browser-macos ./cmd/server
   
   # Windows
   GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o s3-smart-browser.exe ./cmd/server
   ```

## Debugging

### Go Debugging

1. **Use Delve debugger:**
   ```bash
   go install github.com/go-delve/delve/cmd/dlv@latest
   dlv debug ./cmd/server
   ```

2. **Add debug logging:**
   ```go
   log.Printf("Debug: processing file %s", filePath)
   ```

3. **Use pprof for profiling:**
   ```go
   import _ "net/http/pprof"
   
   go func() {
       log.Println(http.ListenAndServe("localhost:6060", nil))
   }()
   ```

### Frontend Debugging

1. **Use browser developer tools:**
   - Console for JavaScript errors
   - Network tab for API calls
   - Elements tab for DOM inspection

2. **Add console logging:**
   ```javascript
   console.log('Debug: loading directory', path);
   ```

### Docker Debugging

1. **Run container with debug mode:**
   ```bash
   docker run -it --rm \
     -p 8080:8080 \
     -e DEBUG=true \
     s3-smart-browser
   ```

2. **Inspect container logs:**
   ```bash
   docker logs container-name
   ```

## Contributing

### Workflow

1. **Create feature branch:**
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Make changes and commit:**
   ```bash
   git add .
   git commit -m "Add new feature"
   ```

3. **Push and create PR:**
   ```bash
   git push origin feature/new-feature
   ```

### Code Review

1. **Self-review checklist:**
   - [ ] Code follows style guidelines
   - [ ] Tests pass
   - [ ] Documentation updated
   - [ ] No breaking changes

2. **Review process:**
   - Address feedback promptly
   - Keep PRs focused
   - Update documentation

### Release Process

1. **Update version:**
   - Update `go.mod` version
   - Update Helm chart version
   - Update Docker image tags

2. **Create release:**
   - Create GitHub release
   - Update documentation
   - Announce changes

## Troubleshooting

### Common Issues

1. **Module not found:**
   ```bash
   go mod tidy
   go mod download
   ```

2. **Docker build fails:**
   ```bash
   docker system prune
   docker build --no-cache -t s3-smart-browser .
   ```

3. **Tests fail:**
   ```bash
   go clean -testcache
   go test ./...
   ```

### Performance Issues

1. **Memory leaks:**
   - Use `go tool pprof` for memory profiling
   - Check for goroutine leaks
   - Monitor memory usage

2. **Slow performance:**
   - Profile CPU usage
   - Check database queries
   - Optimize algorithms

## Resources

### Documentation

- [Go Documentation](https://golang.org/doc/)
- [AWS SDK Go v2](https://aws.github.io/aws-sdk-go-v2/)
- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)

### Tools

- [GoLand](https://www.jetbrains.com/go/)
- [VS Code](https://code.visualstudio.com/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)

---

⭐ **Happy coding!** ⭐
