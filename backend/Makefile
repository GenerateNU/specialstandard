.PHONY: test test-unit test-integration test-db test-coverage test-clean lint lint-fix

# Color codes for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
MAGENTA := \033[0;35m
CYAN := \033[0;36m
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
	@go test ./... $(TEST_FLAGS) -count=1 2>&1 | tee test_output.log
	@make test-summary-detailed

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
	@go test ./... -v -count=1 -short 2>&1 | tee test_output.log
	@make test-summary-detailed

# Run only integration tests
test-integration:
	@echo "$(BOLD)Running integration tests...$(NC)"
	@echo "========================================"
	@go test ./internal/service/... -v -count=1 2>&1 | tee test_output.log
	@make test-summary-detailed

# Run only database tests
test-db:
	@echo "$(BOLD)Running database tests...$(NC)"
	@echo "========================================"
	@go test ./internal/storage/postgres/... -v -count=1 2>&1 | tee test_output.log | grep -v "🐳\|✅\|🚫" | grep -E "(RUN|PASS|FAIL|===)" || true
	@make test-summary-detailed

# Run tests with coverage
test-coverage:
	@echo "$(BOLD)Running tests with coverage...$(NC)"
	@echo "========================================"
	@go test ./... -v -coverprofile=coverage.out 2>&1 | tee test_output.log
	@make test-summary-detailed
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(YELLOW)Coverage report generated: coverage.html$(NC)"

# Enhanced test summary with detailed failure information
test-summary-detailed:
	@echo ""
	@echo "$(BOLD)===== Test Summary =====$(NC)"
	@if grep -q "^--- FAIL" test_output.log; then \
		echo "$(RED)$(BOLD)FAILURES DETECTED:$(NC)"; \
		echo ""; \
		awk ' \
			BEGIN { \
				in_test = 0; \
				test_name = ""; \
				buffer = ""; \
			} \
			/^=== RUN/ { \
				test_name = $$3; \
				buffer = ""; \
			} \
			/^--- FAIL:/ { \
				in_test = 1; \
				print "$(RED)$(BOLD)❌ FAILED: " $$3 "$(NC)"; \
				print "$(CYAN)   Duration: " $$4 "$(NC)"; \
				print ""; \
			} \
			in_test && /^\s+.*_test\.go:[0-9]+:/ { \
				print "$(YELLOW)   " $$0 "$(NC)"; \
			} \
			in_test && /Error Trace:/ { \
				print "$(MAGENTA)   Error Trace:$(NC)"; \
				getline; print "$(MAGENTA)   " $$0 "$(NC)"; \
			} \
			in_test && /Error:/ && !/Error Trace:/ { \
				print "$(RED)   Error:$(NC)"; \
				getline; \
				while ($$0 ~ /^\s+/ && $$0 !~ /Test:/) { \
					print "$(RED)   " $$0 "$(NC)"; \
					getline; \
				} \
			} \
			in_test && /Expected/ { \
				print "$(GREEN)   Expected:$(NC)"; \
				getline; \
				while ($$0 ~ /^\s+/ && $$0 !~ /(Actual|Test:|Messages:)/) { \
					print "$(GREEN)   " $$0 "$(NC)"; \
					getline; \
				} \
			} \
			in_test && /Actual/ { \
				print "$(RED)   Actual:$(NC)"; \
				getline; \
				while ($$0 ~ /^\s+/ && $$0 !~ /(Test:|Messages:|Expected)/) { \
					print "$(RED)   " $$0 "$(NC)"; \
					getline; \
				} \
			} \
			in_test && /Messages:/ { \
				print "$(BLUE)   Messages:$(NC)"; \
				getline; \
				while ($$0 ~ /^\s+/) { \
					print "$(BLUE)   " $$0 "$(NC)"; \
					getline; \
				} \
			} \
			in_test && /Test:/ { \
				print "$(CYAN)   Test: " $$2 "$(NC)"; \
			} \
			in_test && /^\s+.*: / && !/Error:|Expected:|Actual:|Test:|Messages:/ { \
				print "$(YELLOW)   " $$0 "$(NC)"; \
			} \
			/^(===|---)/ { \
				if (in_test) { \
					print "$(BOLD)───────────────────────────────────────$(NC)"; \
					print ""; \
				} \
				in_test = 0; \
			} \
		' test_output.log; \
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

# Original test-summary-full target (now just calls the detailed version)
test-summary-full: test-summary-detailed

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
	@make test-summary-detailed
	@mv test_output_verbose.log test_output.log

# Enhanced: Run tests and show only failures with full details
test-failures:
	@echo "$(BOLD)Running tests (showing only failures)...$(NC)"
	@echo "========================================"
	@go test ./... -v 2>&1 | tee test_output.log
	@echo ""
	@if grep -q "^--- FAIL" test_output.log; then \
		echo "$(RED)$(BOLD)Test Failures:$(NC)"; \
		echo ""; \
		awk ' \
			BEGIN { capture = 0; buffer = "" } \
			/^=== RUN/ { \
				test_name = $$0; \
				buffer = ""; \
				capture = 0; \
			} \
			/^--- FAIL:/ { \
				capture = 1; \
				print test_name; \
				print $$0; \
			} \
			capture && !/^(===|---)/ { \
				print; \
			} \
			/^(===|---)/ && capture { \
				print "───────────────────────────────────────"; \
				print ""; \
				capture = 0; \
			} \
		' test_output.log; \
	else \
		echo "$(GREEN)$(BOLD)All tests passed!$(NC)"; \
	fi
	@rm -f test_output.log

# Enhanced: Run specific test with full output and detailed parsing
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
		echo "$(RED)$(BOLD)Test failed - showing detailed output:$(NC)"; \
		make test-summary-detailed; \
	else \
		echo ""; \
		echo "$(GREEN)$(BOLD)Test passed!$(NC)"; \
	fi
	@rm -f test_output.log

# New target: Pretty print test failures with all details
test-pretty:
	@echo "$(BOLD)Running tests with pretty output...$(NC)"
	@echo "========================================"
	@go test ./... -v 2>&1 | tee test_output.log | awk ' \
		/^=== RUN/ { print "$(BLUE)▶ " $$0 "$(NC)" } \
		/^--- PASS:/ { print "$(GREEN)✓ " $$0 "$(NC)" } \
		/^--- FAIL:/ { print "$(RED)✗ " $$0 "$(NC)" } \
		/^--- SKIP:/ { print "$(YELLOW)⊘ " $$0 "$(NC)" } \
		!/^(===|---)/ { print } \
	'
	@make test-summary-detailed

# New target: Extract only assertion failures
test-assertions:
	@echo "$(BOLD)Extracting assertion failures...$(NC)"
	@echo "========================================"
	@go test ./... -v 2>&1 | tee test_output.log
	@echo ""
	@if grep -q "^--- FAIL" test_output.log; then \
		echo "$(RED)$(BOLD)Assertion Failures:$(NC)"; \
		echo ""; \
		grep -E "(Error:|Expected:|Actual:|assert\.|require\.)" test_output.log | \
		sed 's/^[[:space:]]*/  /' | \
		awk ' \
			/Error:/ { print "$(RED)" $$0 "$(NC)" } \
			/Expected:/ { print "$(GREEN)" $$0 "$(NC)" } \
			/Actual:/ { print "$(RED)" $$0 "$(NC)" } \
			/(assert\.|require\.)/ { print "$(YELLOW)" $$0 "$(NC)" } \
		'; \
	else \
		echo "$(GREEN)$(BOLD)No assertion failures found!$(NC)"; \
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

# CI-specific target without colors but with detailed output
test-ci:
	@echo "Running all tests (CI mode)..."
	@echo "========================================"
	@go test ./... -v 2>&1 | tee test_output.log
	@echo ""
	@echo "===== Test Summary ====="
	@if grep -q "^--- FAIL" test_output.log; then \
		echo "FAILURES DETECTED:"; \
		echo ""; \
		awk ' \
			BEGIN { in_test = 0 } \
			/^--- FAIL:/ { \
				in_test = 1; \
				print "FAILED: " $$3; \
				print "Duration: " $$4; \
			} \
			in_test && /^\s+.*_test\.go:[0-9]+:/ { \
				print "  Location: " $$0; \
			} \
			in_test && /Error:/ { \
				print "  Error:"; \
				getline; \
				while ($$0 ~ /^\s+/ && $$0 !~ /Test:/) { \
					print "    " $$0; \
					getline; \
				} \
			} \
			in_test && /Expected/ { \
				print "  Expected:"; \
				getline; \
				while ($$0 ~ /^\s+/ && $$0 !~ /(Actual|Test:)/) { \
					print "    " $$0; \
					getline; \
				} \
			} \
			in_test && /Actual/ { \
				print "  Actual:"; \
				getline; \
				while ($$0 ~ /^\s+/ && $$0 !~ /Test:/) { \
					print "    " $$0; \
					getline; \
				} \
			} \
			/^(===|---)/ { \
				if (in_test) print "---"; \
				in_test = 0; \
			} \
		' test_output.log; \
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