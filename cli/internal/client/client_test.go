package client

import (
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
