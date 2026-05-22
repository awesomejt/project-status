#!/usr/bin/env bash
# Smoke test script for Project Status API
# Run against a local Docker Compose stack to validate basic API functionality.
#
# Environment variables:
#   API_BASE_URL     - API endpoint URL (default: http://localhost:5000)
#   TEST_RECORD_PREFIX - Prefix for test record short_name (default: smoke-test)
#   SMOKE_CLEANUP    - Enable cleanup on exit (default: true)
#
# Usage:
#   API_BASE_URL=http://localhost:5000 ./scripts/smoke-curl.sh
#
# Exit codes:
#   0 - All checks passed
#   1 - One or more checks failed
#   2 - Usage error or missing prerequisites

set -euo pipefail

# Configuration
API_URL="${API_BASE_URL:-http://localhost:5000}"
TEST_RECORD_PREFIX="${TEST_RECORD_PREFIX:-smoke-test}"
SMOKE_CLEANUP="${SMOKE_CLEANUP:-true}"
SMOKE_RECORD_ID=""
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

# Check if curl is available
if ! command -v curl &>/dev/null; then
    echo "Error: curl is required but not installed."
    exit 2
fi

# Function to make HTTP request and check response
http_request() {
    local method=$1
    local endpoint=$2
    local body=$3
    
    local url="${API_URL}${endpoint}"
    local args=(-s -w "\n%{http_code}" --max-time 10)
    
    if [[ "$method" == "POST" || "$method" == "PATCH" ]]; then
        args+=(--request "$method" --header "Content-Type: application/json")
        if [[ -n "$body" ]]; then
            args+=(--data "$body")
        fi
    else
        args+=(--request "$method")
    fi
    
    args+=("$url")
    
    local response
    response=$(curl "${args[@]}")
    local http_code
    http_code=$(echo "$response" | tail -n1)
    local body_content
    body_content=$(echo "$response" | sed '$d')
    
    # Return format: http_code<TAB>body_content (use TAB as delimiter to avoid conflicts)
    printf '%s\t%s' "$http_code" "$body_content"
}

# Cleanup function
cleanup() {
    if [[ "$SMOKE_CLEANUP" != "true" ]]; then
        log "Skipping cleanup (SMOKE_CLEANUP=$SMOKE_CLEANUP)"
        return
    fi
    if [[ -n "$SMOKE_RECORD_ID" ]]; then
        log "Cleaning up test record..."
        local response
        response=$(http_request "DELETE" "/api/project/status/${SMOKE_RECORD_ID}" "")
        local http_code
        http_code=$(echo "$response" | cut -d$'\t' -f1)
        if [[ "$http_code" == "200" ]]; then
            log_pass "Test record deleted"
        else
            log_fail "Failed to delete test record (HTTP $http_code)"
        fi
    fi
}

# Set trap for cleanup
trap cleanup EXIT

# Main script
echo "========================================"
echo "Project Status - Smoke Test Suite"
echo "API URL: $API_URL"
echo "========================================"
echo ""

# Check 1: Health endpoint
log "Testing health endpoint..."
response=$(http_request "GET" "/health" "")
http_code=$(echo "$response" | cut -d$'\t' -f1)
if [[ "$http_code" == "200" ]]; then
    log_pass "Health endpoint (HTTP $http_code)"
else
    log_fail "Health endpoint failed (HTTP $http_code)"
    exit 1
fi

# Check 2: Readiness endpoint
log "Testing readiness endpoint..."
response=$(http_request "GET" "/ready" "")
http_code=$(echo "$response" | cut -d$'\t' -f1)
if [[ "$http_code" == "200" ]]; then
    log_pass "Readiness endpoint (HTTP $http_code)"
else
    log_fail "Readiness endpoint failed (HTTP $http_code)"
    exit 1
fi

echo ""
log "Testing status record CRUD operations..."
echo ""

# Check 3: Create a status record
log "Creating status record..."
create_body="{\"project_name\":\"Smoke Test Project\",\"short_name\":\"${TEST_RECORD_PREFIX}-$(date +%s)\",\"status\":\"active\",\"phase\":\"implementation\",\"summary\":\"Test record created by smoke script\",\"source\":\"smoke-test\"}"
response=$(http_request "POST" "/api/project/status" "$create_body")
http_code=$(echo "$response" | cut -d$'\t' -f1)
body_content=$(echo "$response" | cut -d$'\t' -f2-)
if [[ "$http_code" == "201" ]]; then
    log_pass "Status record created (HTTP $http_code)"
    # Extract the record ID from response
    SMOKE_RECORD_ID=$(echo "$body_content" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    if [[ -z "$SMOKE_RECORD_ID" ]]; then
        log_fail "Could not extract record ID from response: $body_content"
        exit 1
    fi
    log "Record ID: $SMOKE_RECORD_ID"
else
    log_fail "Failed to create status record (HTTP $http_code)"
    exit 1
fi

# Check 4: Read the status record
log "Reading status record..."
response=$(http_request "GET" "/api/project/status/${SMOKE_RECORD_ID}" "")
http_code=$(echo "$response" | cut -d$'\t' -f1)
if [[ "$http_code" == "200" ]]; then
    log_pass "Status record read (HTTP $http_code)"
else
    log_fail "Failed to read status record (HTTP $http_code)"
    exit 1
fi

# Check 5: List status records
log "Listing status records..."
response=$(http_request "GET" "/api/project/status" "")
http_code=$(echo "$response" | cut -d$'\t' -f1)
body_content=$(echo "$response" | cut -d$'\t' -f2-)
if [[ "$http_code" == "200" ]]; then
    log_pass "Status records listed (HTTP $http_code)"
    # Verify our record is in the list
    if echo "$body_content" | grep -q "$SMOKE_RECORD_ID"; then
        log_pass "Test record found in list"
    else
        log_fail "Test record NOT found in list"
    fi
else
    log_fail "Failed to list status records (HTTP $http_code)"
    exit 1
fi

# Check 6: Update the status record
log "Updating status record..."
update_body='{"status":"completed","summary":"Updated by smoke script"}'
response=$(http_request "PATCH" "/api/project/status/${SMOKE_RECORD_ID}" "$update_body")
http_code=$(echo "$response" | cut -d$'\t' -f1)
if [[ "$http_code" == "200" ]]; then
    log_pass "Status record updated (HTTP $http_code)"
else
    log_fail "Failed to update status record (HTTP $http_code)"
    exit 1
fi

# Check 7: Verify update
log "Verifying update..."
response=$(http_request "GET" "/api/project/status/${SMOKE_RECORD_ID}" "")
http_code=$(echo "$response" | cut -d$'\t' -f1)
body_content=$(echo "$response" | cut -d$'\t' -f2-)
if [[ "$http_code" == "200" ]]; then
    if echo "$body_content" | grep -q '"status":"completed"'; then
        log_pass "Update verified - status is 'completed'"
    else
        log_fail "Update NOT reflected in database"
    fi
else
    log_fail "Failed to read updated record (HTTP $http_code)"
fi

# Check 8: Test validation error
log "Testing validation (invalid status)..."
invalid_body='{"project_name":"Invalid Test","short_name":"invalid","status":"invalid_status"}'
response=$(http_request "POST" "/api/project/status" "$invalid_body")
http_code=$(echo "$response" | cut -d$'\t' -f1)
if [[ "$http_code" == "400" ]]; then
    log_pass "Validation error returned (HTTP $http_code)"
else
    log_fail "Expected 400 for invalid status, got $http_code"
fi

# Check 9: Test not-found error
log "Testing not-found error..."
response=$(http_request "GET" "/api/project/status/00000000-0000-0000-0000-000000000000" "")
http_code=$(echo "$response" | cut -d$'\t' -f1)
if [[ "$http_code" == "404" ]]; then
    log_pass "Not-found error returned (HTTP $http_code)"
else
    log_fail "Expected 404 for non-existent record, got $http_code"
fi

echo ""
log "Cleaning up..."

# Disable trap and run cleanup manually once for clearer output.
trap - EXIT
cleanup

echo ""
echo "========================================"
if [[ $FAILED -eq 0 ]]; then
    echo -e "${GREEN}All smoke tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some smoke tests failed.${NC}"
    exit 1
fi
