## Testing

The backend includes comprehensive unit, integration, and database tests. Tests are organized by type and can be run separately or together. ***Please manually test with postman in addition to written tests as they can be created incorrectly***

### Test Structure

```
internal/
├── service/
│   ├── server_test.go              # Integration tests for API endpoints
│   └── handler/
│       └── session/
│           └── handler_test.go   # Unit tests for handlers
├── storage/
│   ├── mocks/                   # Mock implementations for testing
│   └── postgres/
│       ├── schema/
│       │   └── session_test.go  # Database repository tests
│       └── testutil/            # Test utilities and helpers
```

### Running Tests

We use a Makefile to simplify running different types of tests:

```bash
# Run all tests
make test

# Run only unit tests (fast, no database required)
make test-unit

# Run only integration tests
make test-integration

# Run only database tests (requires Docker)
make test-db

# Run tests with coverage report
make test-coverage

# Clean test cache
make test-clean
```

### Test Types

#### Unit Tests
- Test individual functions and handlers in isolation
- Use mocked dependencies
- Run quickly without external services
- Example: `make test-unit`

#### Integration Tests
- Test complete API endpoints
- Verify HTTP request/response flow
- Use mocked repositories
- Example: `make test-integration`

#### Database Tests
- Test actual database operations
- Use Docker containers with PostgreSQL
- Verify SQL queries and data persistence
- Example: `make test-db`

### Test Dependencies

Install test dependencies:
```bash
go get github.com/stretchr/testify
go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/modules/postgres
go get github.com/google/uuid
```

### Coverage Reports

Generate and view test coverage:
```bash
# Generate coverage report
make test-coverage

# This creates:
# - coverage.out (raw coverage data)
# - coverage.html (HTML report)

# Open the HTML report in your browser
open coverage.html
```

### Writing Tests

#### Example Unit Test
```go
func TestHandler_GetSessions(t *testing.T) {
    // Setup mock
    mockRepo := new(mocks.MockSessionRepository)
    mockRepo.On("GetSessions", mock.Anything).Return(sessions, nil)
    
    // Create handler
    handler := session.NewHandler(mockRepo)
    
    // Test
    req := httptest.NewRequest("GET", "/sessions", nil)
    resp, _ := app.Test(req)
    
    // Assert
    assert.Equal(t, 200, resp.StatusCode)
}
```

#### Example Database Test
```go
func TestSessionRepository_GetSessions(t *testing.T) {
    // Skip if running short tests
    if testing.Short() {
        t.Skip("Skipping database test in short mode")
    }
    
    // Setup test database
    testDB := testutil.SetupTestDB(t)
    defer testDB.Cleanup()
    
    // Test repository methods
    sessions, err := repo.GetSessions(ctx)
    assert.NoError(t, err)
}
```

### Docker Required for Database Tests

Database tests use [Testcontainers](https://testcontainers.com/) to automatically spin up PostgreSQL containers. Ensure Docker is running before running database tests:

```bash
# Start Docker Desktop or Docker Engine
# Then run database tests
make test-db
```

### Postman

For further development and testing, install Postman, which simplifies making network requests. ***Please still manually test your code, as sometimes unit tests can be wrong***
