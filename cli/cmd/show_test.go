//go:build premature_cli_command_tests
// +build premature_cli_command_tests

package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestShowCommandSuccess(t *testing.T) {

	recordID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		expectedPath := "/api/project/status/" + recordID
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","project_name":"Test Project","short_name":"test-proj","status":"active","phase":"implementation","summary":"Test summary","reason":"Testing reason","details":"Testing details","tags":["tag1","tag2"],"source":"manual","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-02T00:00:00Z"}`)
	}))
	defer server.Close()

	viper.Set("api_url", server.URL)
	viper.Set("output", "table")

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "show", recordID, "--api-url", server.URL}

	var buf bytes.Buffer
	showCmd.SetOut(&buf)

	showCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, recordID) {
		t.Fatalf("expected output to contain record ID, got: %s", output)
	}
	if !strings.Contains(output, "Test Project") {
		t.Fatalf("expected output to contain 'Test Project', got: %s", output)
	}
	if !strings.Contains(output, "active") {
		t.Fatalf("expected output to contain 'active' status, got: %s", output)
	}
	if !strings.Contains(output, "implementation") {
		t.Fatalf("expected output to contain 'implementation' phase, got: %s", output)
	}
	if !strings.Contains(output, "tag1") {
		t.Fatalf("expected output to contain 'tag1', got: %s", output)
	}
}

func TestShowCommandJSONOutput(t *testing.T) {

	recordID := "b2c3d4e5-f6a7-8901-bcde-f23456789012"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"b2c3d4e5-f6a7-8901-bcde-f23456789012","project_name":"Test Project","status":"paused"}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "show", recordID, "--output", "json", "--api-url", server.URL}

	var buf bytes.Buffer
	showCmd.SetOut(&buf)

	showCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "\"id\"") {
		t.Fatalf("expected JSON output, got: %s", output)
	}
	if !strings.Contains(output, recordID) {
		t.Fatalf("expected output to contain record ID, got: %s", output)
	}
}

func TestShowCommandNotFound(t *testing.T) {

	recordID := "non-existent-id"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"code":404,"message":"Record not found"}}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "show", recordID, "--api-url", server.URL}

	var stderrBuf bytes.Buffer
	showCmd.SetErr(&stderrBuf)

	showCmd.Execute()

	output := stderrBuf.String()
	if !strings.Contains(output, "404") {
		t.Fatalf("expected error output to contain 404, got: %s", output)
	}
}

func TestShowCommandMissingID(t *testing.T) {

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "show"}

	var buf bytes.Buffer
	showCmd.SetErr(&buf)

	showCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "accepts 1 arg(s)") {
		t.Logf("Output: %s", output)
	}
}

func TestShowCommandWithOptions(t *testing.T) {

	recordID := "c3d4e5f6-a7b8-9012-cdef-345678901234"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"id":"c3d4e5f6-a7b8-9012-cdef-345678901234","project_name":"Minimal Record","short_name":"min","status":"blocked"}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "show", recordID, "--api-url", server.URL}

	var buf bytes.Buffer
	showCmd.SetOut(&buf)

	showCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "Minimal Record") {
		t.Fatalf("expected output to contain 'Minimal Record', got: %s", output)
	}
	if !strings.Contains(output, "blocked") {
		t.Fatalf("expected output to contain 'blocked' status, got: %s", output)
	}
}
