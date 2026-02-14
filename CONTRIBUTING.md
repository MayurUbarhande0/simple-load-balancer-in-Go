# Contributing to Simple Load Balancer

Thank you for considering contributing to this project! Here are some guidelines to help you get started.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/simple-load-balancer-in-Go.git
   cd simple-load-balancer-in-Go
   ```
3. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- Go 1.25.3 or higher
- Make (optional, for using Makefile commands)
- Docker (optional, for containerized testing)

### Install Dependencies

```bash
make install
# or
go mod download
```

### Building

```bash
make build
# or
go build -o lb
```

### Running Tests

```bash
make test
# or
go test ./...
```

### Running with Coverage

```bash
make test-coverage
```

This generates `coverage.html` that you can open in a browser.

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format your code
- Run `go vet` to check for common mistakes
- Use meaningful variable and function names
- Add comments for exported functions and types

### Formatting

```bash
make fmt
# or
go fmt ./...
```

### Linting

```bash
make lint
# or
make vet
```

## Making Changes

### Adding New Features

1. Create an issue describing the feature
2. Discuss the implementation approach
3. Create a branch and implement the feature
4. Add tests for your feature
5. Update documentation (README.md, code comments)
6. Submit a pull request

### Fixing Bugs

1. Create an issue describing the bug (if one doesn't exist)
2. Create a branch with a descriptive name (e.g., `fix/health-check-timeout`)
3. Write a test that reproduces the bug
4. Fix the bug
5. Ensure all tests pass
6. Submit a pull request

### Writing Tests

- Write unit tests for new functions
- Ensure all tests pass before submitting PR
- Aim for good test coverage (>80%)
- Test edge cases and error conditions

Example test structure:

```go
func TestYourFunction(t *testing.T) {
    // Arrange
    input := setupTestData()
    
    // Act
    result := YourFunction(input)
    
    // Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

## Project Structure

```
.
├── main.go              # Entry point
├── config/              # Configuration handling
├── models/              # Data structures
├── routes/              # HTTP handlers
├── helper/              # Utility functions
├── middleware/          # HTTP middleware
├── examples/            # Example code
└── arch/                # Architecture documentation
```

## Commit Messages

Write clear, concise commit messages:

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- First line should be 50 characters or less
- Reference issues and pull requests when relevant

Example:
```
Add retry logic for failed backends (#123)

- Implement exponential backoff
- Add configuration for max retries
- Update tests

Fixes #123
```

## Pull Request Process

1. Update the README.md with details of changes if needed
2. Update the documentation for any API changes
3. Add tests for new functionality
4. Ensure all tests pass
5. Update the CHANGELOG if applicable
6. Request review from maintainers

### PR Checklist

- [ ] Tests added/updated and passing
- [ ] Documentation updated
- [ ] Code follows project style guidelines
- [ ] Commit messages are clear
- [ ] No merge conflicts
- [ ] PR description clearly explains the changes

## Testing Your Changes

### Manual Testing

1. Start backend servers:
   ```bash
   # Terminal 1
   go run examples/backend.go -port 8081
   
   # Terminal 2
   go run examples/backend.go -port 8082
   ```

2. Start the load balancer:
   ```bash
   ./lb
   ```

3. Test the endpoints:
   ```bash
   # Test load balancing
   curl http://localhost:8080/
   
   # Check health
   curl http://localhost:8080/health
   
   # Check metrics
   curl http://localhost:8080/metrics
   ```

### Docker Testing

```bash
docker-compose up
```

## Reporting Issues

When reporting issues, please include:

- Go version (`go version`)
- Operating system and version
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Relevant logs or error messages

## Questions?

Feel free to open an issue for any questions or concerns.

## License

By contributing, you agree that your contributions will be licensed under the project's license.
