package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestUpdateCommandSuccess(t *testing.T) {

	recordID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Fatalf("expected PATCH, got %s", r.Method)
		}
		expectedPath := "/api/project/status/" + recordID
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		if payload["status"] != "blocked" {
			t.Fatalf("expected status 'blocked', got %v", payload["status"])
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","project_name":"Test Project","short_name":"test-proj","status":"blocked","phase":"implementation","summary":"Updated summary"}`)
	}))
	defer server.Close()

	origArgs := os.Args

	os.Args = []string{"project-status", "update", recordID,
		"--status", "blocked",
		"--api-url", server.URL,
	}

	var buf bytes.Buffer
	SetTestOutput(&buf, io.Discard)

	Execute()

	os.Args = origArgs

	output := buf.String()
	if !strings.Contains(output, recordID) {
		t.Fatalf("expected output to contain record ID, got: %s", output)
	}
	if !strings.Contains(output, "blocked") {
		t.Fatalf("expected output to contain 'blocked', got: %s", output)
	}
}

func TestUpdateCommandWithMultipleFields(t *testing.T) {

	recordID := "b2c3d4e5-f6a7-8901-bcde-f23456789012"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		if payload["status"] != "paused" || payload["summary"] != "New summary" || payload["phase"] != "validation" {
			t.Fatalf("unexpected payload: %v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"b2c3d4e5-f6a7-8901-bcde-f23456789012","project_name":"Test Project","short_name":"test-proj","status":"paused","phase":"validation","summary":"New summary"}`)
	}))
	defer server.Close()

	origArgs := os.Args

	os.Args = []string{"project-status", "update", recordID,
		"--status", "paused",
		"--summary", "New summary",
		"--phase", "validation",
		"--api-url", server.URL,
	}

	var buf bytes.Buffer
	SetTestOutput(&buf, io.Discard)

	Execute()

	os.Args = origArgs

	output := buf.String()
	if !strings.Contains(output, "paused") {
		t.Fatalf("expected output to contain 'paused', got: %s", output)
	}
}

func TestUpdateCommandJSONOutput(t *testing.T) {

	recordID := "c3d4e5f6-a7b8-9012-cdef-345678901234"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"c3d4e5f6-a7b8-9012-cdef-345678901234","project_name":"Test Project","status":"active"}`)
	}))
	defer server.Close()

	origArgs := os.Args

	os.Args = []string{"project-status", "update", recordID,
		"--status", "active",
		"--output", "json",
		"--api-url", server.URL,
	}

	var buf bytes.Buffer
	SetTestOutput(&buf, io.Discard)

	Execute()

	os.Args = origArgs

	output := buf.String()
	if !strings.Contains(output, "\"id\"") {
		t.Fatalf("expected JSON output, got: %s", output)
	}
}

func TestUpdateCommandNotFound(t *testing.T) {

	recordID := "non-existent-id"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"code":4404,"message":"Record not found"}}`)
	}))
	defer server.Close()

	origArgs := os.Args

	os.Args = []string{"project-status", "update", recordID,
		"--status", "active",
		"--api-url", server.URL,
	}

	var stderrBuf bytes.Buffer
	SetTestOutput(io.Discard, &stderrBuf)

	Execute()

	os.Args = origArgs

	output := stderrBuf.String()
	if !strings.Contains(output, "404") {
		t.Fatalf("expected error output to contain 404, got: %s", output)
	}
}

func TestUpdateCommandWithTags(t *testing.T) {

	recordID := "d4e5f6a7-b8c9-0123-defg-456789012345"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		tags, ok := payload["tags"].([]interface{})
		if !ok || len(tags) != 3 {
			t.Fatalf("expected 3 tags, got %v", payload["tags"])
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"d4e5f6a7-b8c9-0123-defg-456789012345","project_name":"Test Project","short_name":"test-proj","status":"active","tags":["tag1","tag2","tag3"]}`)
	}))
	defer server.Close()

	origArgs := os.Args

	os.Args = []string{"project-status", "update", recordID,
		"--tags", "tag1,tag2,tag3",
		"--api-url", server.URL,
	}

	var buf bytes.Buffer
	SetTestOutput(&buf, io.Discard)

	Execute()

	os.Args = origArgs

	output := buf.String()
	if !strings.Contains(output, recordID) {
		t.Fatalf("expected output to contain record ID, got: %s", output)
	}
}

func TestUpdateCommandMissingID(t *testing.T) {

	origArgs := os.Args

	os.Args = []string{"project-status", "update"}

	var buf bytes.Buffer
	updateCmd.SetErr(&buf)
	updateCmd.Execute()

	os.Args = origArgs

	output := buf.String()
	if !strings.Contains(output, "accepts 1 arg(s)") {
		t.Logf("Output: %s", output)
	}
}

func TestUpdateCommandValidationFailure(t *testing.T) {

	recordID := "e5f6a7b8-c9d0-1234-efgh-567890123456"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":400,"message":"Validation error","details":"Invalid status value"}}`)
	}))
	defer server.Close()

	origArgs := os.Args

	os.Args = []string{"project-status", "update", recordID,
		"--status", "invalid-status",
		"--api-url", server.URL,
	}

	var stderrBuf bytes.Buffer
	SetTestOutput(io.Discard, &stderrBuf)

	Execute()

	os.Args = origArgs

	output := stderrBuf.String()
	if !strings.Contains(output, "400") {
		t.Fatalf("expected error output to contain 400, got: %s", output)
	}
}
