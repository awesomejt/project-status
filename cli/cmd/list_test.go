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

func TestListCommandSuccess(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"records":[{"id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","project_name":"Test Project","short_name":"test-proj","status":"active","phase":"implementation","summary":"Test summary"},{"id":"b2c3d4e5-f6a7-8901-bcde-f23456789012","project_name":"Another Project","short_name":"another","status":"paused","phase":"planning","summary":"Another summary"}],"total":2,"page":1,"per_page":20,"pages":1}`)
	}))
	defer server.Close()

	viper.Set("api_url", server.URL)
	viper.Set("output", "table")

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "list", "--api-url", server.URL}

	var buf bytes.Buffer
	listCmd.SetOut(&buf)

	listCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "a1b2c3d4-e5f6-7890-abcd-ef1234567890") {
		t.Fatalf("expected output to contain first record ID, got: %s", output)
	}
	if !strings.Contains(output, "Test Project") {
		t.Fatalf("expected output to contain 'Test Project', got: %s", output)
	}
	if !strings.Contains(output, "Total: 2 records") {
		t.Fatalf("expected output to contain 'Total: 2 records', got: %s", output)
	}
}

func TestListCommandWithFilters(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") != "active" {
			t.Fatalf("expected status=active, got %s", r.URL.Query().Get("status"))
		}
		if r.URL.Query().Get("phase") != "validation" {
			t.Fatalf("expected phase=validation, got %s", r.URL.Query().Get("phase"))
		}
		if r.URL.Query().Get("per_page") != "50" {
			t.Fatalf("expected per_page=50, got %s", r.URL.Query().Get("per_page"))
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"records":[],"total":0,"page":1,"per_page":50,"pages":1}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "list",
		"--status", "active",
		"--phase", "validation",
		"--per-page", "50",
		"--api-url", server.URL,
	}

	listCmd.Execute()
}

func TestListCommandJSONOutput(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"records":[{"id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","project_name":"Test Project","status":"active"}],"total":1,"page":1,"per_page":20,"pages":1}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "list", "--output", "json", "--api-url", server.URL}

	var buf bytes.Buffer
	listCmd.SetOut(&buf)

	listCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "\"records\"") {
		t.Fatalf("expected JSON output with records array, got: %s", output)
	}
	if !strings.Contains(output, "a1b2c3d4-e5f6-7890-abcd-ef1234567890") {
		t.Fatalf("expected output to contain record ID, got: %s", output)
	}
}

func TestListCommandEmptyResults(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"records":[],"total":0,"page":1,"per_page":20,"pages":1}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "list", "--api-url", server.URL}

	var buf bytes.Buffer
	listCmd.SetOut(&buf)

	listCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "Total: 0 records") {
		t.Fatalf("expected output to contain 'Total: 0 records', got: %s", output)
	}
}

func TestListCommandAPIError(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error":{"code":500,"message":"Internal server error"}}`)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "list", "--api-url", server.URL}

	var stderrBuf bytes.Buffer
	listCmd.SetErr(&stderrBuf)

	listCmd.Execute()

	output := stderrBuf.String()
	if !strings.Contains(output, "500") {
		t.Fatalf("expected error output to contain 500, got: %s", output)
	}
}
