"""Test the status record API endpoints."""
import pytest
import json


class TestPingEndpoint:
    """Tests for the ping endpoint."""
    
    def test_ping_returns_200(self, client):
        """Test that ping endpoint returns 200 OK."""
        response = client.get("/api/ping")
        assert response.status_code == 200
        data = json.loads(response.data)
        assert data["status"] == "pong"


class TestCreateStatusRecord:
    """Tests for the create status record endpoint."""
    
    def test_create_status_record_success(self, client, sample_status_record_data):
        """Test successful creation of a status record."""
        response = client.post(
            "/api",
            json=sample_status_record_data,
            content_type="application/json"
        )
        
        assert response.status_code == 201
        data = json.loads(response.data)
        
        assert data["project_name"] == "Test Project"
        assert data["short_name"] == "test-proj"
        assert data["status"] == "active"
        assert data["id"] is not None
    
    def test_create_status_record_missing_required_fields(self, client):
        """Test that creation fails without required fields."""
        response = client.post(
            "/api",
            json={"project_name": "Test"},
            content_type="application/json"
        )
        
        assert response.status_code in [400, 422]  # Bad Request or Unprocessable Entity


class TestListStatusRecords:
    """Tests for the list status records endpoint."""
    
    def test_list_status_records_empty(self, client):
        """Test listing status records when none exist."""
        response = client.get("/api")
        
        assert response.status_code == 200
        data = json.loads(response.data)
        
        assert "records" in data
        assert data["total"] == 0
    
    def test_list_status_records_with_data(self, client, sample_status_record):
        """Test listing status records with existing data."""
        response = client.get("/api")
        
        assert response.status_code == 200
        data = json.loads(response.data)
        
        assert len(data["records"]) >= 1
        assert data["total"] >= 1


class TestGetStatusRecord:
    """Tests for the get status record endpoint."""
    
    def test_get_status_record_success(self, client, sample_status_record):
        """Test successful retrieval of a status record."""
        response = client.get(f"/api/{sample_status_record.id}")
        
        assert response.status_code == 200
        data = json.loads(response.data)
        
        assert data["id"] == sample_status_record.id
        assert data["project_name"] == "Test Project"
    
    def test_get_status_record_not_found(self, client):
        """Test retrieving a non-existent status record."""
        response = client.get("/api/99999")
        
        assert response.status_code == 404


class TestUpdateStatusRecord:
    """Tests for the update status record endpoint."""
    
    def test_update_status_record_success(self, client, sample_status_record, sample_status_record_update_data):
        """Test successful update of a status record."""
        response = client.patch(
            f"/api/{sample_status_record.id}",
            json=sample_status_record_update_data,
            content_type="application/json"
        )
        
        assert response.status_code == 200
        data = json.loads(response.data)
        
        assert data["project_name"] == "Updated Test Project"
        assert data["status"] == "working"
    
    def test_update_status_record_not_found(self, client):
        """Test updating a non-existent status record."""
        response = client.patch(
            "/api/99999",
            json={"summary": "Updated"},
            content_type="application/json"
        )
        
        assert response.status_code == 404


class TestDeleteStatusRecord:
    """Tests for the delete status record endpoint."""
    
    def test_delete_status_record_success(self, client, sample_status_record):
        """Test successful deletion of a status record."""
        response = client.delete(f"/api/{sample_status_record.id}")
        
        assert response.status_code == 200
        
        # Verify the record is deleted
        get_response = client.get(f"/api/{sample_status_record.id}")
        assert get_response.status_code == 404
    
    def test_delete_status_record_not_found(self, client):
        """Test deleting a non-existent status record."""
        response = client.delete("/api/99999")
        
        assert response.status_code == 404
