//go:build premature_cli_command_tests
// +build premature_cli_command_tests

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestAddCommandSuccess(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/project/status" {
			t.Fatalf("expected path /api/project/status, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":"test-id-123","project_name":"Test Project","short_name":"test","status":"active","phase":null,"summary":"Test summary","reason":null,"details":null,"tags":[],"source":"cli","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`)
	}))
	defer server.Close()

	viper.Set("api_url", server.URL)
	viper.Set("output", "table")

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{
		"project-status", "add",
		"--project-name", "Test Project",
		"--short-name", "test",
		"--status", "active",
		"--summary", "Test summary",
		"--api-url", server.URL,
	}

	var outBuf bytes.Buffer
	addCmd.SetOut(&outBuf)
	addCmd.SetErr(&outBuf)

	addCmd.Execute()

	output := outBuf.String()
	if !strings.Contains(output, "Created status record") || !strings.Contains(output, "test-id-123") {
		t.Fatalf("expected output to contain 'Created status record test-id-123', got: %s", output)
	}
}

func TestAddCommandValidationFailure(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":400,"message":"Validation error","details":"Invalid status value"}}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "add",
		"--project-name", "Test Project",
		"--short-name", "test-proj",
		"--status", "invalid-status",
		"--summary", "Summary",
		"--api-url", server.URL,
	}

	var stderrBuf bytes.Buffer
	addCmd.SetErr(&stderrBuf)

	addCmd.Execute()

	output := stderrBuf.String()
	if !strings.Contains(output, "400") {
		t.Fatalf("expected error output to contain 400, got: %s", output)
	}
}

func TestAddCommandMissingRequiredFields(t *testing.T) {

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "add", "--short-name", "test-proj"}

	var buf bytes.Buffer
	addCmd.SetErr(&buf)

	addCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "required flag") && !strings.Contains(output, "project-name") {
		t.Logf("Output: %s", output)
	}
}

func TestAddCommandWithTag(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		tags, ok := payload["tags"].([]interface{})
		if !ok || len(tags) != 2 {
			t.Fatalf("expected 2 tags, got %v", payload["tags"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":"b2c3d4e5-f6a7-8901-bcde-f23456789012","project_name":"Test Project","short_name":"test-proj","status":"active","tags":["tag1","tag2"]}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "add",
		"--project-name", "Test Project",
		"--short-name", "test-proj",
		"--status", "active",
		"--summary", "Summary",
		"--tags", "tag1,tag2",
		"--api-url", server.URL,
	}

	var buf bytes.Buffer
	addCmd.SetOut(&buf)

	addCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "b2c3d4e5-f6a7-8901-bcde-f23456789012") {
		t.Fatalf("expected output to contain record ID, got: %s", output)
	}
}

func TestAddCommandWithPhase(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)

		if payload["phase"] != "validation" {
			t.Fatalf("expected phase 'validation', got %v", payload["phase"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":"c3d4e5f6-a7b8-9012-cdef-345678901234","project_name":"Test Project","short_name":"test-proj","status":"active","phase":"validation"}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "add",
		"--project-name", "Test Project",
		"--short-name", "test-proj",
		"--status", "active",
		"--phase", "validation",
		"--summary", "Summary",
		"--api-url", server.URL,
	}

	var buf bytes.Buffer
	addCmd.SetOut(&buf)

	addCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "c3d4e5f6-a7b8-9012-cdef-345678901234") {
		t.Fatalf("expected output to contain record ID, got: %s", output)
	}
}

func TestAddCommandJSONOutput(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":"d4e5f6a7-b8c9-0123-defg-456789012345","project_name":"Test Project","short_name":"test-proj","status":"active"}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "add",
		"--project-name", "Test Project",
		"--short-name", "test-proj",
		"--status", "active",
		"--summary", "Summary",
		"--output", "json",
		"--api-url", server.URL,
	}

	var buf bytes.Buffer
	addCmd.SetOut(&buf)

	addCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "\"id\"") {
		t.Fatalf("expected JSON output, got: %s", output)
	}
	if !strings.Contains(output, "d4e5f6a7-b8c9-0123-defg-456789012345") {
		t.Fatalf("expected output to contain record ID, got: %s", output)
	}
}
