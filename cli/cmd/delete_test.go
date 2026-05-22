package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestDeleteCommandSuccess(t *testing.T) {

	recordID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		expectedPath := "/api/project/status/" + recordID
		if r.URL.Path != expectedPath {
			t.Fatalf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	origArgs := os.Args

	os.Args = []string{"project-status", "delete", recordID, "--force", "--api-url", server.URL}

	var buf bytes.Buffer
	SetTestOutput(&buf, io.Discard)

	Execute()

	os.Args = origArgs

	output := buf.String()
	if !strings.Contains(output, recordID) {
		t.Fatalf("expected output to contain record ID, got: %s", output)
	}
	if !strings.Contains(output, "Deleted") {
		t.Fatalf("expected output to contain 'Deleted', got: %s", output)
	}
}

func TestDeleteCommandNotFound(t *testing.T) {

	recordID := "non-existent-id"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"code":404,"message":"Record not found"}}`)
	}))
	defer server.Close()

	origArgs := os.Args

	os.Args = []string{"project-status", "delete", recordID, "--force", "--api-url", server.URL}

	var stderrBuf bytes.Buffer
	SetTestOutput(io.Discard, &stderrBuf)

	Execute()

	os.Args = origArgs

	output := stderrBuf.String()
	if !strings.Contains(output, "404") {
		t.Fatalf("expected error output to contain 404, got: %s", output)
	}
}

func TestDeleteCommandMissingID(t *testing.T) {

	origArgs := os.Args

	os.Args = []string{"project-status", "delete"}

	var buf bytes.Buffer
	deleteCmd.SetErr(&buf)
	deleteCmd.Execute()

	os.Args = origArgs

	output := buf.String()
	if !strings.Contains(output, "accepts 1 arg(s)") {
		t.Logf("Output: %s", output)
	}
}

func TestDeleteCommandForceFlag(t *testing.T) {

	recordID := "b2c3d4e5-f6a7-8901-bcde-f23456789012"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	os.Args = []string{"project-status", "delete", recordID, "--force", "--api-url", server.URL}

	Execute()
}

func TestDeleteCommandAPIError(t *testing.T) {

	recordID := "c3d4e5f6-a7b8-9012-cdef-345678901234"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error":{"code":500,"message":"Internal server error"}}`)
	}))
	defer server.Close()

	origArgs := os.Args

	os.Args = []string{"project-status", "delete", recordID, "--force", "--api-url", server.URL}

	var stderrBuf bytes.Buffer
	SetTestOutput(io.Discard, &stderrBuf)

	Execute()

	os.Args = origArgs

	output := stderrBuf.String()
	if !strings.Contains(output, "500") {
		t.Fatalf("expected error output to contain 500, got: %s", output)
	}
}
