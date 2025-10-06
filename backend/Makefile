.PHONY: test test-unit test-integration test-db test-coverage

# Run all tests
test:
	go test ./... | grep -v "no test files" || true

# Run only unit tests (skip database tests)
test-unit:
	go test ./... -v -short

# Run only integration tests
test-integration:
	go test ./internal/service/... -v

# Run only database tests
test-db:
	go test ./internal/storage/postgres/... -v

# Run tests with coverage
test-coverage:
	go test ./... -v -coverprofile=coverage_temp.out
	grep -v "/mocks/" coverage_temp.out > coverage.out || true
	rm coverage_temp.out
	go tool cover -html=coverage.out -o coverage.html

# Clean test cache
test-clean:
	go clean -testcache