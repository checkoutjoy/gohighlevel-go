# Contributing to GoHighLevel Go SDK

Thank you for your interest in contributing to the GoHighLevel Go SDK! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful and constructive in all interactions.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/gohighlevel-go.git
   cd gohighlevel-go
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/checkoutjoy/gohighlevel-go.git
   ```

## Development Setup

### Prerequisites

- Go 1.24 or higher
- Make (optional but recommended)
- golangci-lint for linting

### Install Development Tools

```bash
make install-tools
```

Or manually:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Making Changes

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the coding standards below

3. Run tests to ensure everything works:
   ```bash
   make test-unit
   ```

4. Run linting:
   ```bash
   make lint
   ```

5. Format your code:
   ```bash
   make fmt
   ```

## Coding Standards

### Go Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Keep functions focused and small
- Write clear, descriptive variable names
- Add comments for exported functions and types

### Naming Conventions

- Use camelCase for unexported functions and variables
- Use PascalCase for exported functions and types
- Avoid abbreviations unless widely understood
- Use descriptive names that explain purpose

### Error Handling

- Always check and handle errors
- Return errors rather than panicking
- Wrap errors with context using `fmt.Errorf`
- Use descriptive error messages

Example:
```go
if err != nil {
    return fmt.Errorf("failed to create contact: %w", err)
}
```

### Testing

- Write tests for all new features
- Aim for high test coverage
- Use table-driven tests where appropriate
- Integration tests should skip when credentials are not available

Example:
```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Testing

### Unit Tests

Run unit tests without requiring API credentials:

```bash
make test-unit
```

### Integration Tests

Integration tests require valid GoHighLevel credentials:

```bash
export GHL_CLIENT_ID="your-client-id"
export GHL_CLIENT_SECRET="your-client-secret"
export GHL_ACCESS_TOKEN="your-access-token"
export GHL_LOCATION_ID="your-location-id"

make test-integration
```

### Coverage

Generate a coverage report:

```bash
make test
```

This creates `coverage.html` which you can open in a browser.

## Commit Messages

Write clear, descriptive commit messages:

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit first line to 72 characters
- Reference issues and pull requests when applicable

Examples:
```
Add support for contact custom fields

Fix error handling in OAuth token refresh

Update dependencies to latest versions
Fixes #123
```

## Pull Request Process

1. Update documentation if needed (README, code comments, etc.)
2. Add tests for new functionality
3. Ensure all tests pass
4. Update CHANGELOG.md if applicable
5. Push your changes to your fork
6. Create a pull request against the `main` branch
7. Fill out the pull request template completely
8. Wait for review and address any feedback

### Pull Request Checklist

- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Code follows style guidelines
- [ ] All tests pass
- [ ] No linting errors
- [ ] Commits are clean and descriptive

## Adding New Resources

When adding support for a new GoHighLevel resource:

1. Create a new file (e.g., `resource_name.go`)
2. Define the resource service struct:
   ```go
   type ResourceService struct {
       client *Client
   }
   ```
3. Define request/response types
4. Implement CRUD methods with proper error handling
5. Add integration tests in `resource_name_test.go`
6. Document required OAuth scopes in comments
7. Update README.md with usage examples
8. Add to client initialization in `client.go`

## Documentation

- Document all exported functions, types, and constants
- Use complete sentences in comments
- Include examples where helpful
- Update README.md for user-facing changes
- Keep documentation up to date with code changes

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about contributing
- General discussion

Thank you for contributing!
