package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListRecordsUsesProjectStatusPathAndParsesRecords(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/api/project/status"; got != want {
			t.Fatalf("unexpected path: got %s want %s", got, want)
		}
		if got := r.URL.Query().Get("status"); got != "active" {
			t.Fatalf("unexpected status query: %s", got)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"records":[{"id":"3ec273da-4d8a-4ec3-bf0c-e5104f054299","project_name":"proj","short_name":"p","status":"active","summary":"ok","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}],"total":1,"page":1,"per_page":20,"pages":1}`)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	resp, err := c.ListRecords(1, 20, "active", "")
	if err != nil {
		t.Fatalf("ListRecords error: %v", err)
	}

	if len(resp.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(resp.Records))
	}
	if resp.Records[0].ID != "3ec273da-4d8a-4ec3-bf0c-e5104f054299" {
		t.Fatalf("unexpected id: %s", resp.Records[0].ID)
	}
}

func TestGetRecordUsesStringIDPath(t *testing.T) {
	t.Parallel()

	wantID := "3ec273da-4d8a-4ec3-bf0c-e5104f054299"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/api/project/status/"+wantID; got != want {
			t.Fatalf("unexpected path: got %s want %s", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id":"%s","project_name":"proj","short_name":"p","status":"active","summary":"ok","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`, wantID)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	record, err := c.GetRecord(wantID)
	if err != nil {
		t.Fatalf("GetRecord error: %v", err)
	}
	if record.ID != wantID {
		t.Fatalf("unexpected id: %s", record.ID)
	}
}

func TestCreateRecordUsesCorrectPathAndPayload(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "POST"; got != want {
			t.Fatalf("unexpected method: got %s want %s", got, want)
		}
		if got, want := r.URL.Path, "/api/project/status"; got != want {
			t.Fatalf("unexpected path: got %s want %s", got, want)
		}

		var payload StatusRecordCreate
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if payload.ProjectName != "Test Project" || payload.ShortName != "test-proj" || payload.Status != "active" {
			t.Fatalf("unexpected payload: %+v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","project_name":"Test Project","short_name":"test-proj","status":"active","phase":"implementation","summary":"Test record","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	record := StatusRecordCreate{
		ProjectName: "Test Project",
		ShortName:   "test-proj",
		Status:      "active",
		Summary:     "Test record",
	}

	created, err := c.CreateRecord(record)
	if err != nil {
		t.Fatalf("CreateRecord error: %v", err)
	}
	if created.ID != "a1b2c3d4-e5f6-7890-abcd-ef1234567890" {
		t.Fatalf("unexpected id: %s", created.ID)
	}
	if created.Phase == nil || *created.Phase != "implementation" {
		t.Fatalf("unexpected phase: %+v", created.Phase)
	}
}

func TestCreateRecordReturnsErrorOnValidationFailure(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":400,"message":"Validation error","details":"Invalid status value"}}`)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	record := StatusRecordCreate{
		ProjectName: "Test Project",
		ShortName:   "test-proj",
		Status:      "invalid-status",
		Summary:     "Test record",
	}

	_, err := c.CreateRecord(record)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "server error (400): Validation error" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateRecordUsesCorrectPathAndPayload(t *testing.T) {
	t.Parallel()

	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "PATCH"; got != want {
			t.Fatalf("unexpected method: got %s want %s", got, want)
		}
		if got, want := r.URL.Path, "/api/project/status/"+wantID; got != want {
			t.Fatalf("unexpected path: got %s want %s", got, want)
		}

		var payload StatusRecordUpdate
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if payload.Status == nil || *payload.Status != "blocked" {
			t.Fatalf("unexpected payload: %+v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id":"%s","project_name":"Test Project","short_name":"test-proj","status":"blocked","summary":"Test record","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:01Z"}`, wantID)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	status := "blocked"
	record := StatusRecordUpdate{
		Status: &status,
	}

	updated, err := c.UpdateRecord(wantID, record)
	if err != nil {
		t.Fatalf("UpdateRecord error: %v", err)
	}
	if updated.ID != wantID {
		t.Fatalf("unexpected id: %s", updated.ID)
	}
	if updated.Status != "blocked" {
		t.Fatalf("unexpected status: %s", updated.Status)
	}
}

func TestUpdateRecordReturnsErrorOnNotFound(t *testing.T) {
	t.Parallel()

	wantID := "non-existent-id"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"code":404,"message":"Record not found"}}`)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	statusVal := "blocked"
	record := StatusRecordUpdate{
		Status: &statusVal,
	}

	_, err := c.UpdateRecord(wantID, record)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != "server error (404): Record not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteRecordUsesCorrectPath(t *testing.T) {
	t.Parallel()

	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "DELETE"; got != want {
			t.Fatalf("unexpected method: got %s want %s", got, want)
		}
		if got, want := r.URL.Path, "/api/project/status/"+wantID; got != want {
			t.Fatalf("unexpected path: got %s want %s", got, want)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	err := c.DeleteRecord(wantID)
	if err != nil {
		t.Fatalf("DeleteRecord error: %v", err)
	}
}

func TestDeleteRecordReturnsErrorOnFailure(t *testing.T) {
	t.Parallel()

	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"code":404,"message":"Record not found"}}`)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	err := c.DeleteRecord(wantID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidateURLAcceptsValidURLs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		url  string
	}{
		{"http URL", "http://localhost:5000"},
		{"https URL", "https://api.example.com"},
		{"URL with path", "http://localhost:5000/api"},
		{"URL with trailing slash", "http://localhost:5000/"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateURL(tc.url)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.url, err)
			}
		})
	}
}

func TestValidateURLRejectsInvalidURLs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		url  string
	}{
		{"empty string", ""},
		{"no scheme", "localhost:5000"},
		{"no host", "http://"},
		{"invalid format", "not-a-url"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateURL(tc.url)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", tc.url)
			}
		})
	}
}
