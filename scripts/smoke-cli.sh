#!/usr/bin/env bash
# CLI smoke test script for Project Status
# Run against a local Docker Compose stack to validate CLI functionality.
#
# Prerequisites:
#   - CLI binary built at build/project-status
#   - API running at API_BASE_URL
#
# Environment variables:
#   API_BASE_URL     - API endpoint URL (default: http://localhost:5000)
#   CLI_BINARY       - Path to CLI binary (default: ./build/project-status)
#   TEST_RECORD_PREFIX - Prefix for test record short_name (default: cli-smoke)
#   TEST_PROJECT_NAME  - Project name for test record (default: CLI Smoke Test Project)
#   CLI_SMOKE_CLEANUP  - Enable cleanup on exit (default: true)
#
# Usage:
#   API_BASE_URL=http://localhost:5000 ./scripts/smoke-cli.sh
#
# Exit codes:
#   0 - All checks passed
#   1 - One or more checks failed
#   2 - Usage error or missing prerequisites

set -euo pipefail

# Configuration
API_URL="${API_BASE_URL:-http://localhost:5000}"
CLI_BINARY="${CLI_BINARY:-./build/project-status}"
TEST_RECORD_PREFIX="${TEST_RECORD_PREFIX:-cli-smoke}"
TEST_PROJECT_NAME="${TEST_PROJECT_NAME:-CLI Smoke Test Project}"
CLI_SMOKE_CLEANUP="${CLI_SMOKE_CLEANUP:-true}"
CLI_SMOKE_RECORD_ID=""
FAILED=0

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log() {
    echo -e "${YELLOW}→${NC} $1"
}

log_pass() {
    echo -e "${GREEN}✓${NC} $1"
}

log_fail() {
    echo -e "${RED}✗${NC} $1"
    FAILED=1
}

# Check if CLI binary exists
if [[ ! -f "$CLI_BINARY" ]]; then
    echo "Error: CLI binary not found at $CLI_BINARY"
    echo "Build first: make build-cli"
    exit 2
fi

# Cleanup function
cleanup() {
    if [[ -n "$CLI_SMOKE_RECORD_ID" ]]; then
        log "Cleaning up test record..."
        set +e
        "$CLI_BINARY" --api-url "$API_URL" delete "$CLI_SMOKE_RECORD_ID" > /dev/null 2>&1
        set -e
        log_pass "Test record deleted"
    fi
}

# Set trap for cleanup
trap cleanup EXIT

# Main script
echo "========================================"
echo "Project Status - CLI Smoke Test Suite"
echo "API URL: $API_URL"
echo "CLI Binary: $CLI_BINARY"
echo "========================================"
echo ""

# Check 1: CLI help
log "Testing CLI help..."
output=$("$CLI_BINARY" --help 2>&1)
if echo "$output" | grep -q "Usage:"; then
    log_pass "CLI --help works"
else
    log_fail "CLI --help failed"
    exit 1
fi

# Check 2: CLI config show
log "Testing CLI config show..."
output=$("$CLI_BINARY" config show 2>&1)
if echo "$output" | grep -q "API URL"; then
    log_pass "CLI config show works"
else
    log_fail "CLI config show failed"
fi

# Check 3: Create a status record
log "Creating status record via CLI..."
set +e
output=$("$CLI_BINARY" \
    --api-url "$API_URL" \
    add \
    --project-name "$TEST_PROJECT_NAME" \
    --short-name "${TEST_RECORD_PREFIX}-1" \
    --status "active" \
    --phase "implementation" \
    --summary "Test record created by CLI smoke script" 2>&1)
exit_code=$?
set -e

if [[ $exit_code -eq 0 ]]; then
    log_pass "Status record created via CLI"
    # Extract record ID from output (format: "Created status record <ID>")
    if echo "$output" | grep -q "Created status record"; then
        CLI_SMOKE_RECORD_ID=$(echo "$output" | grep "Created status record" | awk '{print $4}')
        log "Record ID: $CLI_SMOKE_RECORD_ID"
    fi
else
    log_fail "Failed to create status record (exit code: $exit_code)"
    echo "Output: $output"
    exit 1
fi

# Check 4: List status records
log "Listing status records via CLI..."
set +e
output=$("$CLI_BINARY" --api-url "$API_URL" list 2>&1)
exit_code=$?
set -e

if [[ $exit_code -eq 0 ]]; then
    log_pass "Status records listed via CLI"
    # Verify our record is in the list
    if echo "$output" | grep -q "${TEST_RECORD_PREFIX}"; then
        log_pass "Test record found in list"
    else
        log_fail "Test record NOT found in list"
    fi
else
    log_fail "Failed to list status records (exit code: $exit_code)"
    echo "Output: $output"
    exit 1
fi

# Check 5: Show status record
log "Showing status record via CLI..."
set +e
output=$("$CLI_BINARY" --api-url "$API_URL" show "$CLI_SMOKE_RECORD_ID" 2>&1)
exit_code=$?
set -e

if [[ $exit_code -eq 0 ]]; then
    log_pass "Status record shown via CLI"
    if echo "$output" | grep -q "$TEST_PROJECT_NAME"; then
        log_pass "Correct project name displayed"
    else
        log_fail "Wrong project name displayed"
    fi
else
    log_fail "Failed to show status record (exit code: $exit_code)"
    echo "Output: $output"
    exit 1
fi

# Check 6: Update status record
log "Updating status record via CLI..."
set +e
output=$("$CLI_BINARY" \
    --api-url "$API_URL" \
    update "$CLI_SMOKE_RECORD_ID" \
    --status "completed" \
    --summary "Updated by CLI smoke script" 2>&1)
exit_code=$?
set -e

if [[ $exit_code -eq 0 ]]; then
    log_pass "Status record updated via CLI"
else
    log_fail "Failed to update status record (exit code: $exit_code)"
    echo "Output: $output"
    exit 1
fi

# Check 7: Verify update
log "Verifying update via CLI..."
set +e
output=$("$CLI_BINARY" --api-url "$API_URL" show "$CLI_SMOKE_RECORD_ID" 2>&1)
exit_code=$?
set -e

if [[ $exit_code -eq 0 ]]; then
    if echo "$output" | grep -q "completed"; then
        log_pass "Update verified - status is 'completed'"
    else
        log_fail "Update NOT reflected"
    fi
else
    log_fail "Failed to show updated record (exit code: $exit_code)"
fi

# Check 8: Test validation error (invalid status)
log "Testing validation (invalid status via CLI)..."
set +e
output=$("$CLI_BINARY" \
    --api-url "$API_URL" \
    add \
    --project-name "Invalid Test" \
    --short-name "invalid-test" \
    --status "invalid_status" \
    --summary "Should fail" 2>&1)
exit_code=$?
set -e

if [[ $exit_code -ne 0 ]]; then
    log_pass "Validation error returned (exit code: $exit_code)"
else
    log_fail "Expected non-zero exit for invalid status"
fi

# Check 9: Test not-found error
log "Testing not-found error via CLI..."
set +e
output=$("$CLI_BINARY" --api-url "$API_URL" show "00000000-0000-0000-0000-000000000000" 2>&1)
exit_code=$?
set -e

if [[ $exit_code -ne 0 ]]; then
    log_pass "Not-found error returned (exit code: $exit_code)"
else
    log_fail "Expected non-zero exit for non-existent record"
fi

# Check 10: JSON output format
log "Testing JSON output format..."
set +e
output=$("$CLI_BINARY" --api-url "$API_URL" list --output json 2>&1)
exit_code=$?
set -e

if [[ $exit_code -eq 0 ]]; then
    if echo "$output" | grep -q "\["; then
        log_pass "JSON output format works"
    else
        log_fail "JSON output format incorrect"
    fi
else
    log_fail "Failed to get JSON output (exit code: $exit_code)"
fi

echo ""
log "Cleaning up..."

# Disable trap and cleanup manually
trap - EXIT
if [[ -n "$CLI_SMOKE_RECORD_ID" ]]; then
    set +e
    "$CLI_BINARY" --api-url "$API_URL" delete "$CLI_SMOKE_RECORD_ID" > /dev/null 2>&1
    set -e
    log_pass "Test record deleted"
fi

echo ""
echo "========================================"
if [[ $FAILED -eq 0 ]]; then
    echo -e "${GREEN}All CLI smoke tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some CLI smoke tests failed.${NC}"
    exit 1
fi
