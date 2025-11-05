.PHONY: test test-unit test-integration test-db test-coverage test-clean lint lint-fix

# Color codes for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m # No Color
BOLD := \033[1m

# Test output formatting
TEST_FLAGS := -v
# Add JSON output for better parsing
TEST_FLAGS_JSON := -json

# Run all tests with improved output
test:
	@echo "$(BOLD)Running all tests...$(NC)"
	@echo "========================================"
	@go test ./... $(TEST_FLAGS) 2>&1 | tee test_output.log
	@make test-summary-full

# Run tests with cleaner output (suppress logs)
test-quiet:
	@echo "$(BOLD)Running all tests (quiet mode)...$(NC)"
	@echo "========================================"
	@go test ./... -v 2>&1 | tee test_output.log | grep -E "(PASS|FAIL|RUN|===)"
	@make test-summary

# Run tests with JSON output for better parsing
test-json:
	@echo "$(BOLD)Running all tests (JSON output)...$(NC)"
	@go test ./... -json > test_results.json
	@go run scripts/parse_test_results.go test_results.json
	@rm -f test_results.json

# Run only unit tests with clean output
test-unit:
	@echo "$(BOLD)Running unit tests...$(NC)"
	@echo "========================================"
	@go test ./... -v -short 2>&1 | tee test_output.log
	@make test-summary-full

# Run only integration tests
test-integration:
	@echo "$(BOLD)Running integration tests...$(NC)"
	@echo "========================================"
	@go test ./internal/service/... -v 2>&1 | tee test_output.log
	@make test-summary-full

# Run only database tests
test-db:
	@echo "$(BOLD)Running database tests...$(NC)"
	@echo "========================================"
	@go test ./internal/storage/postgres/... -v 2>&1 | tee test_output.log | grep -v "🐳\|✅\|🚫" | grep -E "(RUN|PASS|FAIL|===)" || true
	@make test-summary-full

# Run tests with coverage
test-coverage:
	@echo "$(BOLD)Running tests with coverage...$(NC)"
	@echo "========================================"
	@go test ./... -v -coverprofile=coverage.out 2>&1 | tee test_output.log
	@make test-summary-full
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(YELLOW)Coverage report generated: coverage.html$(NC)"

# New helper target for full test summary with failure details
test-summary-full:
	@echo ""
	@echo "$(BOLD)===== Test Summary =====$(NC)"
	@if grep -q "^--- FAIL" test_output.log; then \
		echo "$(RED)FAILURES DETECTED:$(NC)"; \
		echo ""; \
		awk '/^--- FAIL/{found=1} found{print; if (/^(===|---)/ && !/^--- FAIL/) found=0}' test_output.log | grep -v "^===" | sed 's/^/  /'; \
		echo ""; \
	fi
	@echo "$(GREEN)Passed: $$(grep -c '^--- PASS' test_output.log || echo 0)$(NC)"
	@echo "$(RED)Failed: $$(grep -c '^--- FAIL' test_output.log || echo 0)$(NC)"
	@echo "$(YELLOW)Skipped: $$(grep -c '^--- SKIP' test_output.log || echo 0)$(NC)"
	@if grep -q "^--- FAIL" test_output.log; then \
		echo ""; \
		echo "$(RED)$(BOLD)⚠️  TEST SUITE FAILED ⚠️$(NC)"; \
		echo "$(YELLOW)Run 'make test-verbose' to see full output$(NC)"; \
		exit 1; \
	else \
		echo ""; \
		echo "$(GREEN)$(BOLD)✅ ALL TESTS PASSED ✅$(NC)"; \
	fi
	# Only remove test_output.log if NOT in CI
	@if [ -z "$$CI" ]; then \
		rm -f test_output.log; \
	fi

# Original summary target (minimal output)
test-summary:
	@echo ""
	@echo "$(BOLD)===== Test Summary =====$(NC)"
	@if grep -q "^--- FAIL" test_output.log; then \
		echo "$(RED)FAILURES DETECTED:$(NC)"; \
		echo ""; \
		grep -B 2 -A 3 "^--- FAIL" test_output.log | sed 's/^/  /'; \
		echo ""; \
	fi
	@echo "$(GREEN)Passed: $$(grep -c '^--- PASS' test_output.log || echo 0)$(NC)"
	@echo "$(RED)Failed: $$(grep -c '^--- FAIL' test_output.log || echo 0)$(NC)"
	@echo "$(YELLOW)Skipped: $$(grep -c '^--- SKIP' test_output.log || echo 0)$(NC)"
	@if grep -q "^--- FAIL" test_output.log; then \
		echo ""; \
		echo "$(RED)$(BOLD)⚠️  TEST SUITE FAILED ⚠️$(NC)"; \
		echo "$(YELLOW)Run 'make test-verbose' to see full output$(NC)"; \
		exit 1; \
	else \
		echo ""; \
		echo "$(GREEN)$(BOLD)✅ ALL TESTS PASSED ✅$(NC)"; \
	fi
	@rm -f test_output.log

# Run tests with full verbose output (for debugging)
test-verbose:
	@echo "$(BOLD)Running all tests (verbose)...$(NC)"
	@echo "========================================"
	@go test ./... -v 2>&1 | tee test_output_verbose.log
	@make test-summary-full
	@mv test_output_verbose.log test_output.log

# New target: Run tests and show only failures with full details
test-failures:
	@echo "$(BOLD)Running tests (showing only failures)...$(NC)"
	@echo "========================================"
	@go test ./... -v 2>&1 | tee test_output.log
	@echo ""
	@if grep -q "^--- FAIL" test_output.log; then \
		echo "$(RED)$(BOLD)Test Failures:$(NC)"; \
		echo ""; \
		awk '/^=== RUN/{test=$$0} /^---/{if ($$2 == "FAIL:") {print test; found=1}} found{print} /^(===|---)/{found=0}' test_output.log; \
	else \
		echo "$(GREEN)$(BOLD)All tests passed!$(NC)"; \
	fi
	@rm -f test_output.log

# New target: Run specific test with full output
test-one:
	@if [ -z "$(TEST)" ]; then \
		echo "$(RED)Please specify a test with TEST=TestName$(NC)"; \
		echo "Example: make test-one TEST=TestThemeRepository_GetThemes_WithFilters"; \
		exit 1; \
	fi
	@echo "$(BOLD)Running test: $(TEST)$(NC)"
	@echo "========================================"
	@go test ./... -v -run "$(TEST)" 2>&1 | tee test_output.log
	@if grep -q "FAIL" test_output.log; then \
		echo ""; \
		echo "$(RED)$(BOLD)Test failed - showing full output:$(NC)"; \
		cat test_output.log; \
	fi
	@rm -f test_output.log

# Watch mode - rerun tests on file changes
test-watch:
	@echo "$(BOLD)Running tests in watch mode...$(NC)"
	@echo "Press Ctrl+C to stop"
	@while true; do \
		clear; \
		make test-quiet; \
		echo ""; \
		echo "$(YELLOW)Watching for changes...$(NC)"; \
		fswatch -1 -e ".*" -i "\\.go$$" .; \
	done

# Clean test cache
test-clean:
	go clean -testcache
	@echo "$(GREEN)Test cache cleaned$(NC)"

# CI-specific target without colors and without removing the log file
test-ci:
	@echo "Running all tests (CI mode)..."
	@echo "========================================"
	@go test ./... -v 2>&1 | tee test_output.log
	@echo ""
	@echo "===== Test Summary ====="
	@if grep -q "^--- FAIL" test_output.log; then \
		echo "FAILURES DETECTED:"; \
		echo ""; \
		grep -A 5 "^--- FAIL" test_output.log; \
		echo ""; \
	fi
	@echo "Passed: $$(grep -c '^--- PASS' test_output.log || echo 0)"
	@echo "Failed: $$(grep -c '^--- FAIL' test_output.log || echo 0)"
	@echo "Skipped: $$(grep -c '^--- SKIP' test_output.log || echo 0)"
	@if grep -q "^--- FAIL" test_output.log; then \
		echo ""; \
		echo "TEST SUITE FAILED"; \
		exit 1; \
	else \
		echo ""; \
		echo "ALL TESTS PASSED"; \
	fi
	# DO NOT remove test_output.log here - GitHub Actions needs it!

# CI coverage target that preserves the log file
test-coverage-ci:
	@echo "Running tests with coverage (CI mode)..."
	@go test ./... -v -coverprofile=coverage.out 2>&1 | tee test_output.log
	@echo ""
	@echo "===== Coverage Summary ====="
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	# DO NOT remove test_output.log here!

# Extract coverage percentage for CI
coverage-percentage:
	@go tool cover -func=coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}'

# Lint all code
lint:
	@echo "$(BOLD)Running linter...$(NC)"
	@if golangci-lint run; then \
		echo "$(GREEN)✅ No linting issues found$(NC)"; \
	else \
		echo "$(RED)❌ Linting issues detected$(NC)"; \
		exit 1; \
	fi

# Fix auto-fixable linting issues
lint-fix:
	@echo "$(BOLD)Fixing linting issues...$(NC)"
	golangci-lint run --fix
	@echo "$(GREEN)✅ Linting issues fixed$(NC)"