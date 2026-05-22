# API Documentation

Project Status API - REST API for managing status records.

## Base URL

```
http://localhost:5000/api
```

## Endpoints

### Health Check

```
GET /health
```

Check if the API is running.

**Response:** 200 OK

```json
{
  "status": "healthy",
  "service": "project-status-api"
}
```

### Readiness Check

```
GET /ready
```

Check if the API and database are ready.

**Response:** 200 OK

```json
{
  "status": "ready",
  "database": "connected"
}
```

**Response:** 503 Service Unavailable (when database is down)

```json
{
  "status": "not-ready",
  "database": "disconnected",
  "error": "..."
}
```

### Create Status Record

```
POST /api
```

Create a new status record.

**Request Body:**

| Field | Type | Required | Description | Constraints |
|-------|------|----------|-------------|-------------|
| project_name | string | Yes | Human-readable project name | Max 255 characters |
| short_name | string | Yes | Short project identifier | Max 50 characters, unique |
| status | string | Yes | Status value | See allowed values below |
| phase | string | No | Workflow phase | Max 50 characters |
| summary | string | No | Short status summary | Max 500 characters |
| reason | string | No | Explanation for state |
| details | string | No | Longer notes |
| tags | array | No | String labels for filtering | String array |
| source | string | No | Origin (api, web, cli, agent) | Max 50 characters |

**Allowed status values:**

- `active` - Project is actively being worked on
- `paused` - Project is temporarily paused
- `blocked` - Project is blocked waiting on something
- `working` - Agent is actively working
- `error` - Project has an error state
- `stopped` - Project is stopped
- `completed` - Project is finished

**Example Request:**

```json
{
  "project_name": "Project Status",
  "short_name": "status",
  "status": "active",
  "phase": "implementation",
  "summary": "Building the API and clients",
  "tags": ["api", "web", "cli"],
  "source": "api"
}
```

**Response:** 201 Created

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "project_name": "Project Status",
  "short_name": "status",
  "status": "active",
  "phase": "implementation",
  "summary": "Building the API and clients",
  "reason": null,
  "details": null,
  "tags": ["api", "web", "cli"],
  "source": "api",
  "created_at": "2026-05-21T19:00:00Z",
  "updated_at": "2026-05-21T19:00:00Z"
}
```

**Error Response:** 400 Bad Request

```json
{
  "error": {
    "code": 400,
    "message": "Validation error",
    "details": "Field 'project_name' is required"
  }
}
```

### List Status Records

```
GET /api
```

List all status records with pagination and filtering.

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | integer | 1 | Page number (1-indexed) |
| per_page | integer | 20 | Records per page (max 100) |
| status | string | - | Filter by status |

**Example:** `GET /api?page=1&per_page=10&status=active`

**Response:** 200 OK

```json
{
  "records": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "project_name": "Project Status",
      "short_name": "status",
      "status": "active",
      "phase": "implementation",
      "summary": "Building the API and clients",
      "created_at": "2026-05-21T19:00:00Z",
      "updated_at": "2026-05-21T19:00:00Z"
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 1,
  "pages": 1
}
```

### Get Status Record

```
GET /api/{id}
```

Get a specific status record by ID.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | uuid | The status record ID |

**Response:** 200 OK

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "project_name": "Project Status",
  "short_name": "status",
  "status": "active",
  "phase": "implementation",
  "summary": "Building the API and clients",
  "reason": null,
  "details": null,
  "tags": ["api", "web", "cli"],
  "source": "api",
  "created_at": "2026-05-21T19:00:00Z",
  "updated_at": "2026-05-21T19:00:00Z"
}
```

**Error Response:** 404 Not Found

```json
{
  "error": "Record not found"
}
```

### Update Status Record

```
PATCH /api/{id}
```

Partially update a status record.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | uuid | The status record ID |

**Request Body (any subset of updatable fields):**

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| project_name | string | Human-readable project name | Max 255 characters |
| short_name | string | Short project identifier | Max 50 characters, unique |
| status | string | Status value | See allowed values |
| phase | string | Workflow phase | Max 50 characters |
| summary | string | Short status summary | Max 500 characters |
| reason | string | Explanation for state |
| details | string | Longer notes |
| tags | array | String labels for filtering | String array |

**Example Request:**

```json
{
  "status": "blocked",
  "reason": "Waiting on API key approval"
}
```

**Response:** 200 OK

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "project_name": "Project Status",
  "short_name": "status",
  "status": "blocked",
  "phase": "implementation",
  "summary": "Building the API and clients",
  "reason": "Waiting on API key approval",
  "details": null,
  "tags": ["api", "web", "cli"],
  "source": "api",
  "created_at": "2026-05-21T19:00:00Z",
  "updated_at": "2026-05-21T19:05:00Z"
}
```

**Error Response:** 400 Bad Request

```json
{
  "error": {
    "code": 400,
    "message": "Validation error",
    "details": "Invalid status value"
  }
}
```

### Delete Status Record

```
DELETE /api/{id}
```

Delete a status record.

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| id | uuid | The status record ID |

**Response:** 200 OK

```json
{
  "message": "Record deleted"
}
```

**Error Response:** 404 Not Found

```json
{
  "error": "Record not found"
}
```

## Error Format

All API errors follow a consistent format:

```json
{
  "error": {
    "code": 400,
    "message": "Human-readable error message",
    "details": "Optional additional details"
  }
}
```

## Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created |
| 400 | Bad Request - Invalid input |
| 404 | Not Found - Resource doesn't exist |
| 500 | Internal Server Error - Server error |
| 503 | Service Unavailable - Database not ready |
