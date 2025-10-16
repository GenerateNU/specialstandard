.PHONY: test test-unit test-integration test-db test-coverage test-clean lint lint-fix

# Run all tests and summarize pass/fail counts
test:
	@echo "Running all tests..."
	@go test ./... -v | tee test_output.log
	@echo ""
	@echo "===== Test Summary ====="
	@echo "Passed: $$(grep -c '^--- PASS' test_output.log || true)"
	@echo "Failed: $$(grep -c '^--- FAIL' test_output.log || true)"
	@rm -f test_output.log

# Run only unit tests (skip database tests)
test-unit:
	@echo "Running unit tests..."
	@go test ./... -v -short | tee test_output.log
	@echo ""
	@echo "===== Unit Test Summary ====="
	@echo "Passed: $$(grep -c '^--- PASS' test_output.log || true)"
	@echo "Failed: $$(grep -c '^--- FAIL' test_output.log || true)"
	@rm -f test_output.log

# Run only integration tests
test-integration:
	@echo "Running integration tests..."
	@go test ./internal/service/... -v | tee test_output.log
	@echo ""
	@echo "===== Integration Test Summary ====="
	@echo "Passed: $$(grep -c '^--- PASS' test_output.log || true)"
	@echo "Failed: $$(grep -c '^--- FAIL' test_output.log || true)"
	@rm -f test_output.log

# Run only database tests
test-db:
	@echo "Running database tests..."
	@go test ./internal/storage/postgres/... -v | tee test_output.log
	@echo ""
	@echo "===== Database Test Summary ====="
	@echo "Passed: $$(grep -c '^--- PASS' test_output.log || true)"
	@echo "Failed: $$(grep -c '^--- FAIL' test_output.log || true)"
	@rm -f test_output.log

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -v -coverprofile=coverage.out | tee test_output.log
	@echo ""
	@echo "===== Coverage Test Summary ====="
	@echo "Passed: $$(grep -c '^--- PASS' test_output.log || true)"
	@echo "Failed: $$(grep -c '^--- FAIL' test_output.log || true)"
	@go tool cover -html=coverage.out -o coverage.html
	@rm -f test_output.log

# Clean test cache
test-clean:
	go clean -testcache

# Lint all code
lint:
	golangci-lint run

# Fix auto-fixable linting issues
lint-fix:
	golangci-lint run --fix
