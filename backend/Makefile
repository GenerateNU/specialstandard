.PHONY: test test-unit test-integration test-db test-coverage test-clean lint lint-fix \
        test-quiet test-json test-summary-detailed test-summary test-verbose test-failures \
        test-one test-pretty test-assertions test-watch test-ci test-coverage-ci \
        coverage-percentage test-summary-full test-pretty test-assertions clean

# ------------------------
# OS detection (cross-platform helpers)
# ------------------------
OS := $(shell uname 2>/dev/null || echo Windows_NT)
ifeq ($(OS),Windows_NT)
PRINT_FILE = type
DELETE_FILE = del
FIND = findstr
COUNT_CMD = find /C /V ""
TRUE_CMD = ver >nul
SHELL_CMD = cmd /C
else
PRINT_FILE = cat
DELETE_FILE = rm -f
FIND = grep
COUNT_CMD = grep -c
TRUE_CMD = true
SHELL_CMD = sh -c
endif

# ------------------------
# Color codes
# ------------------------
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
MAGENTA := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m
BOLD := \033[1m

# ------------------------
# Test flags
# ------------------------
TEST_FLAGS := -v
TEST_FLAGS_JSON := -json

# Ensure we always run tests fresh
NO_CACHE := -count=1

# ------------------------
# Helper: safely print file (cat on unix, type on windows)
# Use: @$(PRINT_FILE) test_output.log
# ------------------------

# ------------------------
# Run all tests
# ------------------------
test:
	@echo "$(BOLD)Running all tests...$(NC)"
	@echo "========================================"
	@go test ./... $(TEST_FLAGS) $(NO_CACHE) 2>&1 | tee test_output.log
	@$(MAKE) test-summary-detailed

# ------------------------
# Quiet test (filtered)
# ------------------------
test-quiet:
	@echo "$(BOLD)Running all tests (quiet mode)...$(NC)"
	@echo "========================================"
	@go test ./... -v $(NO_CACHE) 2>&1 | tee test_output.log
	@$(PRINT_FILE) test_output.log | $(FIND) -E "PASS|FAIL|RUN|===" || $(TRUE_CMD)
	@$(MAKE) test-summary

# ------------------------
# JSON test output
# ------------------------
test-json:
	@echo "$(BOLD)Running all tests (JSON output)...$(NC)"
	@go test ./... $(TEST_FLAGS_JSON) > test_results.json
	@go run scripts/parse_test_results.go test_results.json
	@$(DELETE_FILE) test_results.json 2>/dev/null || $(TRUE_CMD)

# ------------------------
# Unit tests
# ------------------------
test-unit:
	@echo "$(BOLD)Running unit tests...$(NC)"
	@echo "========================================"
	@go test ./... -v $(NO_CACHE) -short 2>&1 | tee test_output.log
	@$(MAKE) test-summary-detailed

# ------------------------
# Integration tests
# ------------------------
test-integration:
	@echo "$(BOLD)Running integration tests...$(NC)"
	@echo "========================================"
	@go test ./internal/service/... -v $(NO_CACHE) 2>&1 | tee test_output.log
	@$(MAKE) test-summary-detailed

# ------------------------
# DB tests
# ------------------------
test-db:
	@echo "$(BOLD)Running database tests...$(NC)"
	@echo "========================================"
	@go test ./internal/storage/postgres/... -v $(NO_CACHE) 2>&1 | tee test_output.log
	# show only relevant lines (RUN/PASS/FAIL/===), ignore emoji lines
	@$(PRINT_FILE) test_output.log | $(FIND) -v "🐳" | $(FIND) -v "✅" | $(FIND) -v "🚫" | $(FIND) -E "RUN|PASS|FAIL|===" || $(TRUE_CMD)
	@$(MAKE) test-summary-detailed

# ------------------------
# Coverage run
# ------------------------
test-coverage:
	@echo "$(BOLD)Running tests with coverage...$(NC)"
	@echo "========================================"
	@go test ./... -v -coverprofile=coverage.out $(NO_CACHE) 2>&1 | tee test_output.log
	@$(MAKE) test-summary-detailed
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(YELLOW)Coverage report generated: coverage.html$(NC)"

# ------------------------
# test-summary-detailed (full awk parsing restored from original)
# ------------------------
test-summary-detailed:
	@echo ""
	@echo "$(BOLD)===== Test Summary =====$(NC)"
	@if $(FIND) -q "^--- FAIL" ./test_output.log 2>/dev/null || ( $(FIND) -q "^--- FAIL" ./test_output.log >/dev/null 2>&1 ); then \
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
				printf "$(RED)$(BOLD)❌ FAILED: %s$(NC)\n", $$3; \
			} \
			in_test && /^\s+.*_test\.go:[0-9]+:/ { \
				printf "$(YELLOW)   %s$(NC)\n", $$0; \
			} \
			in_test && /Error Trace:/ { \
				printf "$(MAGENTA)   Error Trace:$(NC)\n"; \
				getline; printf "$(MAGENTA)   %s$(NC)\n", $$0; \
			} \
			in_test && /Error:/ && !/Error Trace:/ { \
				printf "$(RED)   Error:$(NC)\n"; \
				getline; \
				while ($$0 ~ /^\s+/ && $$0 !~ /Test:/) { \
					printf "$(RED)   %s$(NC)\n", $$0; \
					getline; \
				} \
			} \
			in_test && /Expected/ { \
				printf "$(GREEN)   Expected:$(NC)\n"; \
				getline; \
				while ($$0 ~ /^\s+/ && $$0 !~ /(Actual|Test:|Messages:)/) { \
					printf "$(GREEN)   %s$(NC)\n", $$0; \
					getline; \
				} \
			} \
			in_test && /Actual/ { \
				printf "$(RED)   Actual:$(NC)\n"; \
				getline; \
				while ($$0 ~ /^\s+/ && $$0 !~ /(Test:|Messages:|Expected)/) { \
					printf "$(RED)   %s$(NC)\n", $$0; \
					getline; \
				} \
			} \
			in_test && /Messages:/ { \
				printf "$(BLUE)   Messages:$(NC)\n"; \
				getline; \
				while ($$0 ~ /^\s+/) { \
					printf "$(BLUE)   %s$(NC)\n", $$0; \
					getline; \
				} \
			} \
			in_test && /Test:/ { \
				printf "$(CYAN)   Test: %s$(NC)\n", $$2; \
			} \
			in_test && /^\s+.*: / && !/Error:|Expected:|Actual:|Test:|Messages:/ { \
				printf "$(YELLOW)   %s$(NC)\n", $$0; \
			} \
			/^(===|---)/ { \
				if (in_test) { \
					printf "$(BOLD)───────────────────────────────────────$(NC)\n\n"; \
				} \
				in_test = 0; \
			} \
		' ./test_output.log; \
		echo ""; \
	fi
	@echo "$(GREEN)Passed: $$( ($(COUNT_CMD) '^--- PASS' ./test_output.log 2>/dev/null) || echo $$( $(SHELL_CMD) '$(COUNT_CMD) ^--- PASS ./test_output.log 2>/dev/null || echo 0' ) )$(NC)"
	@echo "$(RED)Failed: $$( ($(COUNT_CMD) '^--- FAIL' ./test_output.log 2>/dev/null) || echo $$( $(SHELL_CMD) '$(COUNT_CMD) ^--- FAIL ./test_output.log 2>/dev/null || echo 0' ) )$(NC)"
	@echo "$(YELLOW)Skipped: $$( ($(COUNT_CMD) '^--- SKIP' ./test_output.log 2>/dev/null) || echo $$( $(SHELL_CMD) '$(COUNT_CMD) ^--- SKIP ./test_output.log 2>/dev/null || echo 0' ) )$(NC)"
	@if $(FIND) -q "^--- FAIL" ./test_output.log 2>/dev/null || ( $(FIND) -q "^--- FAIL" ./test_output.log >/dev/null 2>&1 ); then \
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
		$(DELETE_FILE) ./test_output.log 2>/dev/null || $(TRUE_CMD); \
	fi

# ------------------------
# alias full summary to detailed
# ------------------------
test-summary-full: test-summary-detailed

# ------------------------
# Minimal summary (original)
# ------------------------
test-summary:
	@echo ""
	@echo "$(BOLD)===== Test Summary =====$(NC)"
	@if $(FIND) -q "^--- FAIL" ./test_output.log 2>/dev/null || ( $(FIND) -q "^--- FAIL" ./test_output.log >/dev/null 2>&1 ); then \
		echo "$(RED)FAILURES DETECTED:$(NC)"; \
		echo ""; \
		$(PRINT_FILE) ./test_output.log | $(FIND) -B 2 -A 3 "^--- FAIL" | sed 's/^/  /' 2>/dev/null || $(TRUE_CMD); \
		echo ""; \
	fi
	@echo "$(GREEN)Passed: $$( $(COUNT_CMD) '^--- PASS' ./test_output.log 2>/dev/null || echo 0 )$(NC)"
	@echo "$(RED)Failed: $$( $(COUNT_CMD) '^--- FAIL' ./test_output.log 2>/dev/null || echo 0 )$(NC)"
	@echo "$(YELLOW)Skipped: $$( $(COUNT_CMD) '^--- SKIP' ./test_output.log 2>/dev/null || echo 0 )$(NC)"
	@if $(FIND) -q "^--- FAIL" ./test_output.log 2>/dev/null || ( $(FIND) -q "^--- FAIL" ./test_output.log >/dev/null 2>&1 ); then \
		echo ""; \
		echo "$(RED)$(BOLD)⚠️  TEST SUITE FAILED ⚠️$(NC)"; \
		echo "$(YELLOW)Run 'make test-verbose' to see full output$(NC)"; \
		exit 1; \
	else \
		echo ""; \
		echo "$(GREEN)$(BOLD)✅ ALL TESTS PASSED ✅$(NC)"; \
	fi
	@$(DELETE_FILE) ./test_output.log 2>/dev/null || $(TRUE_CMD)

# ------------------------
# Verbose run (keeps the verbose file as 'test_output_verbose.log')
# ------------------------
test-verbose:
	@echo "$(BOLD)Running all tests (verbose)...$(NC)"
	@echo "========================================"
	@go test ./... -v $(NO_CACHE) 2>&1 | tee test_output_verbose.log
	@$(MAKE) test-summary-detailed
	@mv test_output_verbose.log test_output.log 2>/dev/null || $(TRUE_CMD)

# ------------------------
# Show only failures with full details
# ------------------------
test-failures:
	@echo "$(BOLD)Running tests (showing only failures)...$(NC)"
	@echo "========================================"
	@go test ./... -v $(NO_CACHE) 2>&1 | tee test_output.log
	@echo ""
	@if $(FIND) -q "^--- FAIL" ./test_output.log 2>/dev/null || ( $(FIND) -q "^--- FAIL" ./test_output.log >/dev/null 2>&1 ); then \
		echo "$(RED)$(BOLD)Test Failures:$(NC)"; \
		echo ""; \
		awk ' \
			BEGIN { capture = 0; buffer = "" } \
			/^=== RUN/ { test_name = $$0; buffer = ""; capture = 0; } \
			/^--- FAIL:/ { capture = 1; print test_name; print $$0; } \
			capture && !/^(===|---)/ { print; } \
			/^(===|---)/ && capture { print "───────────────────────────────────────"; print ""; capture = 0; } \
		' ./test_output.log; \
	else \
		echo "$(GREEN)$(BOLD)All tests passed!$(NC)"; \
	fi
	@$(DELETE_FILE) ./test_output.log 2>/dev/null || $(TRUE_CMD)

# ------------------------
# Run a single test (by TEST=Name)
# ------------------------
test-one:
	@if [ -z "$(TEST)" ]; then \
		echo "$(RED)Please specify a test with TEST=TestName$(NC)"; \
		echo "Example: make test-one TEST=TestThemeRepository_GetThemes_WithFilters"; \
		exit 1; \
	fi
	@echo "$(BOLD)Running test: $(TEST)$(NC)"
	@echo "========================================"
	@go test ./... -v -run "$(TEST)" $(NO_CACHE) 2>&1 | tee test_output.log
	@if $(FIND) -q "FAIL" ./test_output.log 2>/dev/null || ( $(FIND) -q "FAIL" ./test_output.log >/dev/null 2>&1 ); then \
		echo ""; \
		echo "$(RED)$(BOLD)Test failed - showing detailed output:$(NC)"; \
		$(MAKE) test-summary-detailed; \
	else \
		echo ""; \
		echo "$(GREEN)$(BOLD)Test passed!$(NC)"; \
	fi
	@$(DELETE_FILE) ./test_output.log 2>/dev/null || $(TRUE_CMD)

# ------------------------
# Pretty printer for test output
# ------------------------
test-pretty:
	@echo "$(BOLD)Running tests with pretty output...$(NC)"
	@echo "========================================"
	@go test ./... -v $(NO_CACHE) 2>&1 | tee test_output.log | awk ' \
		/^=== RUN/ { print "$(BLUE)▶ " $$0 "$(NC)" } \
		/^--- PASS:/ { print "$(GREEN)✓ " $$0 "$(NC)" } \
		/^--- FAIL:/ { print "$(RED)✗ " $$0 "$(NC)" } \
		/^--- SKIP:/ { print "$(YELLOW)⊘ " $$0 "$(NC)" } \
		!/^(===|---)/ { print } \
	'
	@$(MAKE) test-summary-detailed

# ------------------------
# Extract assertion failures
# ------------------------
test-assertions:
	@echo "$(BOLD)Extracting assertion failures...$(NC)"
	@echo "========================================"
	@go test ./... -v $(NO_CACHE) 2>&1 | tee test_output.log
	@echo ""
	@if $(FIND) -q "^--- FAIL" ./test_output.log 2>/dev/null || ( $(FIND) -q "^--- FAIL" ./test_output.log >/dev/null 2>&1 ); then \
		echo "$(RED)$(BOLD)Assertion Failures:$(NC)"; \
		echo ""; \
		$(FIND) -E "Error:|Expected:|Actual:|assert\.|require\." ./test_output.log | sed 's/^[[:space:]]*/  /' | awk ' \
			/Error:/ { print "'"$(RED)"'" $$0 "'"$(NC)"'" } \
			/Expected:/ { print "'"$(GREEN)"'" $$0 "'"$(NC)"'" } \
			/Actual:/ { print "'"$(RED)"'" $$0 "'"$(NC)"'" } \
			/(assert\.|require\.)/ { print "'"$(YELLOW)"'" $$0 "'"$(NC)"'" } \
		'; \
	else \
		echo "$(GREEN)$(BOLD)No assertion failures found!$(NC)"; \
	fi
	@$(DELETE_FILE) ./test_output.log 2>/dev/null || $(TRUE_CMD)

# ------------------------
# Watch mode (file watcher). Requires fswatch on macOS or an alternative on Linux/Windows.
# ------------------------
test-watch:
	@echo "$(BOLD)Running tests in watch mode...$(NC)"
	@echo "Press Ctrl+C to stop"
	@while true; do \
		clear; \
		$(MAKE) test-quiet; \
		echo ""; \
		echo "$(YELLOW)Watching for changes...$(NC)"; \
		# fswatch on macOS, inotifywait on Linux, no portable built-in on Windows here \
		if [ "$(OS)" = "Windows_NT" ]; then \
			powershell -Command "Get-ChildItem -Recurse -Filter *.go | Select-Object -ExpandProperty FullName | Out-Null"; \
			sleep 1; \
		else \
			fswatch -1 -e ".*" -i "\\.go$$" .; \
		fi; \
	done

# ------------------------
# Clean test cache
# ------------------------
test-clean:
	go clean -testcache
	@echo "$(GREEN)Test cache cleaned$(NC)"

# ------------------------
# CI-specific target (no colors, but detailed)
# ------------------------
test-ci:
	@echo "Running all tests (CI mode)..."
	@echo "========================================"
	@go test ./... -v $(NO_CACHE) 2>&1 | tee test_output.log
	@echo ""
	@echo "===== Test Summary ====="
	@if $(FIND) -q "^--- FAIL" ./test_output.log 2>/dev/null || ( $(FIND) -q "^--- FAIL" ./test_output.log >/dev/null 2>&1 ); then \
		echo "FAILURES DETECTED:"; \
		echo ""; \
		awk ' \
			BEGIN { in_test = 0 } \
			/^--- FAIL:/ { in_test = 1; print "FAILED: " $$3; print "Duration: " $$4; } \
			in_test && /^\s+.*_test\.go:[0-9]+:/ { print "  Location: " $$0; } \
			in_test && /Error:/ { print "  Error:"; getline; while ($$0 ~ /^\s+/ && $$0 !~ /Test:/) { print "    " $$0; getline; } } \
			in_test && /Expected/ { print "  Expected:"; getline; while ($$0 ~ /^\s+/ && $$0 !~ /(Actual|Test:)/) { print "    " $$0; getline; } } \
			in_test && /Actual/ { print "  Actual:"; getline; while ($$0 ~ /^\s+/ && $$0 !~ /Test:/) { print "    " $$0; getline; } } \
			/^(===|---)/ { if (in_test) print "---"; in_test = 0; } \
		' ./test_output.log; \
		echo ""; \
	fi
	@echo "Passed: $$( $(COUNT_CMD) '^--- PASS' ./test_output.log 2>/dev/null || echo 0 )"
	@echo "Failed: $$( $(COUNT_CMD) '^--- FAIL' ./test_output.log 2>/dev/null || echo 0 )"
	@echo "Skipped: $$( $(COUNT_CMD) '^--- SKIP' ./test_output.log 2>/dev/null || echo 0 )"
	@if $(FIND) -q "^--- FAIL" ./test_output.log 2>/dev/null || ( $(FIND) -q "^--- FAIL" ./test_output.log >/dev/null 2>&1 ); then \
		echo ""; \
		echo "TEST SUITE FAILED"; \
		exit 1; \
	else \
		echo ""; \
		echo "ALL TESTS PASSED"; \
	fi
	# DO NOT remove test_output.log here - GitHub Actions needs it!

# ------------------------
# CI coverage target that preserves the log file
# ------------------------
test-coverage-ci:
	@echo "Running tests with coverage (CI mode)..."
	@go test ./... -v -coverprofile=coverage.out $(NO_CACHE) 2>&1 | tee test_output.log
	@echo ""
	@echo "===== Coverage Summary ====="
	@go tool cover -func=coverage.out | $(FIND) -E "total|^total" || $(TRUE_CMD)
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	# DO NOT remove test_output.log here!

# ------------------------
# Extract coverage percentage for CI
# ------------------------
coverage-percentage:
	@go tool cover -func=coverage.out | $(FIND) -E "total" | awk '{print substr($$3, 1, length($$3)-1)}'

# ------------------------
# Linting
# ------------------------
lint:
	@echo "$(BOLD)Running linter...$(NC)"
	@if golangci-lint run; then \
		echo "$(GREEN)✅ No linting issues found$(NC)"; \
	else \
		echo "$(RED)❌ Linting issues detected$(NC)"; \
		exit 1; \
	fi

lint-fix:
	@echo "$(BOLD)Fixing linting issues...$(NC)"
	golangci-lint run --fix
	@echo "$(GREEN)✅ Linting issues fixed$(NC)"

# ------------------------
# Final cleanup (log file)
# ------------------------
clean:
	@$(DELETE_FILE) ./test_output.log 2>/dev/null || $(TRUE_CMD)
	@$(DELETE_FILE) ./coverage.out 2>/dev/null || $(TRUE_CMD)
	@$(DELETE_FILE) ./coverage.html 2>/dev/null || $(TRUE_CMD)
	@echo "Cleaned test artifacts"
