# API Contract

Canonical HTTP contract for `Project Status` MVP.

## Base

- Primary namespace: `/api/project/status`
- Compatibility namespace: `/api` (temporary while clients migrate)

## Error Shape

All non-2xx API errors should use:

```json
{
  "error": {
    "code": 400,
    "message": "Validation error",
    "details": "optional"
  }
}
```

## Endpoints

### `GET /api/project/status`

List records with pagination and optional filters.

Query parameters:

- `page` integer, default `1`, min `1`, max `10000`
- `per_page` integer, default `20`, min `1`, max `100`
- `status` one of: `active`, `paused`, `blocked`, `working`, `error`, `stopped`, `completed`
- `phase` one of: `planning`, `implementation`, `validation`, `release`

Response body:

```json
{
  "records": [
    {
      "id": "uuid-string",
      "project_name": "Project Status",
      "short_name": "status",
      "status": "active",
      "phase": "implementation",
      "summary": "...",
      "created_at": "2026-05-22T16:00:00",
      "updated_at": "2026-05-22T16:00:00"
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 1,
  "pages": 1
}
```

### `POST /api/project/status`

Create record.

Required fields:

- `project_name` string max 255
- `short_name` string max 50
- `status` allowed status value

Optional fields:

- `phase` string
- `summary` string max 500
- `reason` string
- `details` string
- `tags` array of strings
- `source` string max 50

Response:

- `201` with full record body.

### `GET /api/project/status/{id}`

Read one record by ID string.

Response:

- `200` with full record body.
- `404` error shape when not found.

### `PATCH /api/project/status/{id}`

Partial update. Any updatable fields from create payload may be included.

Response:

- `200` with updated full record body.
- `404` error shape when not found.

### `DELETE /api/project/status/{id}`

Delete one record by ID string.

Response:

- `200` with `{ "message": "Record deleted" }`
- `404` error shape when not found.

## Health Endpoints

- `GET /health` returns app health.
- `GET /ready` returns app + database readiness.

## Notes

- API route handlers accept both trailing slash and non-trailing slash forms for collection paths.
- IDs are currently stored as string values in DB and treated as string path parameters.
