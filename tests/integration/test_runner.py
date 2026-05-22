import http.client
import json
import os
import sys
import uuid
from urllib.parse import urlparse

BASE_URL = os.environ.get("API_BASE_URL", "http://localhost:5000")
TEST_PROJECT_NAME = os.environ.get("TEST_PROJECT_NAME", "Integration Test Project")
TEST_RECORD_PREFIX = os.environ.get("TEST_RECORD_PREFIX", "int-test")
INTEGRATION_CLEANUP = os.environ.get("INTEGRATION_CLEANUP", "true").lower() == "true"
TEST_RECORD_ID = None


def make_request(method: str, path: str, body=None) -> tuple[int, dict | None]:
    """Make HTTP request and return (status_code, response_json)."""
    parsed = urlparse(BASE_URL)
    host = parsed.hostname or "localhost"
    port = 443 if parsed.scheme == "https" else 80
    
    if parsed.port:
        port = parsed.port
    
    conn = http.client.HTTPSConnection(host, port, timeout=30) if parsed.scheme == "https" else http.client.HTTPConnection(host, port, timeout=30)
    
    headers = {"Content-Type": "application/json"} if body else {}
    body_str = json.dumps(body) if body else None
    
    try:
        conn.request(method, path, body=body_str, headers=headers)
        resp = conn.getresponse()
        data = resp.read().decode()
        return resp.status, json.loads(data) if data else None
    finally:
        conn.close()


def test_health():
    """Test health endpoint."""
    status, _ = make_request("GET", "/health")
    assert status == 200, f"Health check failed: {status}"
    return True


def test_readiness():
    """Test readiness endpoint."""
    status, _ = make_request("GET", "/ready")
    assert status == 200, f"Readiness check failed: {status}"
    return True


def test_create_record():
    """Test creating a status record."""
    global TEST_RECORD_ID
    body = {
        "project_name": TEST_PROJECT_NAME,
        "short_name": f"{TEST_RECORD_PREFIX}-{uuid.uuid4().hex[:8]}",
        "summary": "Integration test record",
        "status": "active",
        "phase": "planning",
        "source": "integration-test",
        "tags": ["test", "integration"],
    }
    status, resp = make_request("POST", "/api/project/status", body)
    assert status == 201, f"Create record failed: {status} - {resp}"
    assert "id" in resp, "Response missing 'id' field"
    TEST_RECORD_ID = resp["id"]
    return True


def test_list_records():
    """Test listing status records."""
    status, resp = make_request("GET", "/api/project/status")
    assert status == 200, f"List records failed: {status} - {resp}"
    assert "records" in resp, "Response missing 'records' field"
    assert "total" in resp, "Response missing 'total' field"
    return True


def test_read_record():
    """Test reading a specific status record."""
    status, resp = make_request("GET", f"/api/project/status/{TEST_RECORD_ID}")
    assert status == 200, f"Read record failed: {status} - {resp}"
    assert resp["id"] == TEST_RECORD_ID, "Record ID mismatch"
    assert resp["short_name"].startswith("int-test-"), "Record short_name mismatch"
    return True


def test_update_record():
    """Test updating a status record."""
    body = {
        "status": "working",
        "summary": "Updated integration test summary",
    }
    status, resp = make_request("PATCH", f"/api/project/status/{TEST_RECORD_ID}", body)
    assert status == 200, f"Update record failed: {status} - {resp}"
    assert resp["status"] == "working", "Status not updated"
    return True


def test_validation_errors():
    """Test validation error responses."""
    invalid_bodies = [
        ({"project_name": TEST_PROJECT_NAME, "short_name": "invalid-status", "status": "invalid_status"}, 400),
        ({"project_name": TEST_PROJECT_NAME, "short_name": "x" * 101, "status": "active"}, 400),
    ]
    for body, expected_status in invalid_bodies:
        status, resp = make_request("POST", "/api/project/status", body)
        assert status == expected_status, f"Validation error test failed for {body}: got {status}, expected {expected_status}"
    return True


def test_not_found_error():
    """Test not-found error response."""
    fake_id = str(uuid.uuid4())
    status, resp = make_request("GET", f"/api/project/status/{fake_id}")
    assert status == 404, f"Not-found test failed: got {status}, expected 404"
    assert "error" in resp, "Error response missing 'error' field"
    assert resp["error"]["code"] == 404, "Error code mismatch"
    return True


def test_pagination():
    """Test pagination parameters."""
    status, resp = make_request("GET", "/api/project/status?page=1&per_page=10")
    assert status == 200, f"Pagination test failed: {status} - {resp}"
    
    invalid_requests = [
        ("/api/project/status?page=0", 400),
        ("/api/project/status?page=10001", 400),
        ("/api/project/status?per_page=0", 400),
        ("/api/project/status?per_page=101", 400),
    ]
    for path, expected_status in invalid_requests:
        status, _ = make_request("GET", path)
        assert status == expected_status, f"Pagination validation failed for {path}: got {status}, expected {expected_status}"
    return True


def test_filtering():
    """Test filter parameters."""
    status, resp = make_request("GET", f"/api/project/status?status=working")
    assert status == 200, f"Filter test failed: {status} - {resp}"
    
    status, resp = make_request("GET", f"/api/project/status?phase=planning")
    assert status == 200, f"Phase filter test failed: {status} - {resp}"
    
    status, _ = make_request("GET", "/api/project/status?status=invalid")
    assert status == 400, "Invalid status filter should return 400"
    
    status, _ = make_request("GET", "/api/project/status?phase=invalid")
    assert status == 400, "Invalid phase filter should return 400"
    return True


def test_delete_record():
    """Test deleting a status record."""
    global TEST_RECORD_ID
    status, resp = make_request("DELETE", f"/api/project/status/{TEST_RECORD_ID}")
    assert status == 200, f"Delete record failed: {status} - {resp}"
    
    status, resp = make_request("GET", f"/api/project/status/{TEST_RECORD_ID}")
    assert status == 404, "Record should be deleted"
    TEST_RECORD_ID = None
    return True


def main():
    """Run all integration tests."""
    global TEST_RECORD_ID
    
    tests = [
        ("Health check", test_health),
        ("Readiness check", test_readiness),
        ("Create record", test_create_record),
        ("List records", test_list_records),
        ("Read record", test_read_record),
        ("Update record", test_update_record),
        ("Validation errors", test_validation_errors),
        ("Not-found error", test_not_found_error),
        ("Pagination", test_pagination),
        ("Filtering", test_filtering),
        ("Delete record", test_delete_record),
    ]
    
    print(f"Integration Tests: {BASE_URL}")
    print("=" * 60)
    
    failed = []
    for name, test_func in tests:
        try:
            test_func()
            print(f"  [PASS] {name}")
        except AssertionError as e:
            print(f"  [FAIL] {name}: {e}")
            failed.append(name)
        except Exception as e:
            print(f"  [ERROR] {name}: {e}")
            failed.append(name)
    
    print("=" * 60)
    if INTEGRATION_CLEANUP and TEST_RECORD_ID:
        cleanup_status, _ = make_request("DELETE", f"/api/project/status/{TEST_RECORD_ID}")
        if cleanup_status != 200:
            print(f"  [WARN] Cleanup failed for record {TEST_RECORD_ID}: HTTP {cleanup_status}")
        else:
            print(f"  [CLEANUP] Deleted test record {TEST_RECORD_ID}")

    if failed:
        print(f"FAILED: {len(failed)} test(s) failed:")
        for name in failed:
            print(f"  - {name}")
        sys.exit(1)
    else:
        print("PASSED: All tests passed")
        sys.exit(0)


if __name__ == "__main__":
    main()
