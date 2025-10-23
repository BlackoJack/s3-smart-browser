# Contributing to S3 Smart Browser

Thank you for your interest in contributing to S3 Smart Browser! We welcome contributions from the community and appreciate your help in making this project better.

## How to Contribute

### Reporting Issues

Before creating an issue, please:
1. Check if the issue already exists
2. Search through closed issues
3. Try the latest version

When creating an issue, please include:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Docker version, etc.)
- Logs if applicable

### Suggesting Features

We welcome feature suggestions! Please:
1. Check existing issues and discussions
2. Describe the feature clearly
3. Explain the use case
4. Consider implementation complexity

### Code Contributions

#### Getting Started

1. **Fork the repository**
2. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/s3-smart-browser.git
   cd s3-smart-browser
   ```

3. **Create a feature branch:**
   ```bash
   git checkout -b feature/amazing-feature
   ```

4. **Set up development environment:**
   ```bash
   go mod download
   ```

#### Development Guidelines

- **Code Style**: Follow Go conventions and use `gofmt`
- **Testing**: Add tests for new functionality
- **Documentation**: Update documentation for new features
- **Commits**: Use clear, descriptive commit messages

#### Testing

Run tests before submitting:
```bash
go test ./...
go vet ./...
```

#### Building

Test your changes:
```bash
go build ./cmd/server
```

#### Docker Testing

Test with Docker:
```bash
docker build -t s3-smart-browser-test .
docker run -p 8080:8080 s3-smart-browser-test
```

### Pull Request Process

1. **Update your branch** with the latest changes from main
2. **Ensure tests pass** and code is properly formatted
3. **Update documentation** if needed
4. **Create a pull request** with:
   - Clear description of changes
   - Reference to related issues
   - Screenshots for UI changes
   - Testing instructions

### Review Process

- All PRs require review
- Maintainers will review within 48 hours
- Address feedback promptly
- Keep PRs focused and atomic

## Development Setup

### Prerequisites

- Go 1.25+
- Docker (optional)
- Git

### Local Development

1. **Clone repository:**
   ```bash
   git clone https://github.com/BlackoJack/s3-smart-browser.git
   cd s3-smart-browser
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set environment variables:**
   ```bash
   export AWS_ACCESS_KEY_ID=your_access_key
   export AWS_SECRET_ACCESS_KEY=your_secret_key
   export AWS_BUCKET=your-bucket-name
   export AWS_REGION=us-east-1
   ```

4. **Run application:**
   ```bash
   go run cmd/server/main.go
   ```

### Docker Development

1. **Build image:**
   ```bash
   docker build -t s3-smart-browser-dev .
   ```

2. **Run container:**
   ```bash
   docker run -p 8080:8080 \
     -e AWS_ACCESS_KEY_ID=your_access_key \
     -e AWS_SECRET_ACCESS_KEY=your_secret_key \
     -e AWS_BUCKET=your-bucket-name \
     s3-smart-browser-dev
   ```

## Project Structure

```
s3-smart-browser/
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   ├── s3/             # S3 client implementation
│   └── types/          # Data types
├── web/
│   ├── static/         # Static assets (CSS, JS)
│   └── templates/      # HTML templates
├── charts/             # Helm charts
├── docs/              # Documentation
└── tests/             # Test files
```

## Code Style

### Go Code

- Use `gofmt` for formatting
- Follow Go naming conventions
- Add comments for exported functions
- Use meaningful variable names
- Handle errors explicitly

### Frontend Code

- Use modern JavaScript (ES6+)
- Follow consistent naming conventions
- Add comments for complex logic
- Ensure mobile responsiveness

### Documentation

- Update README.md for significant changes
- Add API documentation for new endpoints
- Include examples for new features
- Update deployment examples if needed

## Areas for Contribution

### High Priority

- [ ] User authentication and authorization
- [ ] File upload functionality
- [ ] File preview capabilities
- [ ] Search functionality
- [ ] Performance optimizations

### Medium Priority

- [ ] Additional S3-compatible storage support
- [ ] Custom themes and branding
- [ ] Bulk operations
- [ ] API rate limiting
- [ ] Metrics and monitoring

### Low Priority

- [ ] Multi-language support
- [ ] Advanced security features
- [ ] WebSocket support
- [ ] Progressive Web App features

## Community Guidelines

### Be Respectful

- Be kind and respectful to all contributors
- Welcome newcomers and help them learn
- Provide constructive feedback
- Focus on the code, not the person

### Communication

- Use clear, descriptive language
- Ask questions when unsure
- Share knowledge and help others
- Be patient with responses

### Quality

- Write clean, maintainable code
- Test your changes thoroughly
- Document your work
- Follow project conventions

## Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For general questions
- **Email**: devopsnskru@gmail.com
- **Telegram**: [@devopsnsk](https://t.me/devopsnsk)

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

## License

By contributing to S3 Smart Browser, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to S3 Smart Browser! 🚀
