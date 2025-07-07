# Claude Development Notes

## Development Workflow

From this point forward, development should follow a pull request-based workflow:

1. Create a new feature branch for each task
2. Make changes on the feature branch
3. Create a pull request to merge changes into main
4. Review and merge the pull request

## Project Structure

- `main.go` - Main application with CLI interface
- `rancast` - Compiled binary
- `go.mod` - Go module definition
- `README.md` - User documentation

## Current Features

- Current weather information retrieval
- 7-day weather forecast
- JMA (Japan Meteorological Agency) data via Open-Meteo API
- Support for 10 major Japanese cities
- No API key required

## Testing

Tests should be added to ensure code quality and reliability. Test files should follow Go conventions (`*_test.go`).

### Running Tests

```bash
# Run all tests
go test

# Run tests with verbose output
go test -v

# Run only unit tests (skip integration tests)
go test -short

# Run tests with coverage
go test -cover

# Run benchmarks
go test -bench=.
```

### Test Structure

- `main_test.go` - Unit tests for core functions
- `integration_test.go` - Integration tests that make actual API calls

Integration tests require internet connection and can be skipped with `-short` flag.

## API Information

- Data Source: Open-Meteo JMA API
- Endpoint: https://api.open-meteo.com/v1/jma
- No authentication required
- Timezone: Asia/Tokyo