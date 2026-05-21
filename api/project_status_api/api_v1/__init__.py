from flask import Blueprint, request, jsonify
from uuid import UUID
from ..models import StatusRecord
from .. import db

bp = Blueprint("api_v1", __name__)


@bp.route("/status-records", methods=["POST"])
def create_status_record():
    """Create a new status record."""
    data = request.get_json()
    
    if not data:
        return jsonify({"error": "Request body required"}), 400
    
    required_fields = ["project_name", "short_name", "status"]
    for field in required_fields:
        if field not in data:
            return jsonify({"error": f"Missing required field: {field}"}), 400
    
    valid_statuses = ["active", "paused", "blocked", "working", "error", "stopped", "completed"]
    if data["status"] not in valid_statuses:
        return jsonify({"error": f"Invalid status. Must be one of: {valid_statuses}"}), 400
    
    record = StatusRecord(
        project_name=data["project_name"],
        short_name=data["short_name"],
        status=data["status"],
        phase=data.get("phase"),
        summary=data.get("summary"),
        reason=data.get("reason"),
        details=data.get("details"),
        tags=data.get("tags", []),
        source=data.get("source", "api"),
    )
    
    db.add(record)
    db.commit()
    
    response = {
        "id": str(record.id),
        "project_name": record.project_name,
        "short_name": record.short_name,
        "status": record.status,
        "phase": record.phase,
        "summary": record.summary,
        "reason": record.reason,
        "details": record.details,
        "tags": record.tags,
        "source": record.source,
        "created_at": record.created_at.isoformat(),
        "updated_at": record.updated_at.isoformat(),
    }
    return jsonify(response), 201


@bp.route("/status-records", methods=["GET"])
def list_status_records():
    """List status records with pagination and filters."""
    page = request.args.get("page", 1, type=int)
    per_page = request.args.get("per_page", 20, type=int)
    status_filter = request.args.get("status")
    
    query = db.query(StatusRecord)
    
    if status_filter:
        query = query.filter(StatusRecord.status == status_filter)
    
    query = query.order_by(StatusRecord.created_at.desc())
    pagination = query.paginate(page=page, per_page=per_page, error_out=False)
    
    records = [
        {
            "id": str(r.id),
            "project_name": r.project_name,
            "short_name": r.short_name,
            "status": r.status,
            "phase": r.phase,
            "summary": r.summary,
            "created_at": r.created_at.isoformat(),
            "updated_at": r.updated_at.isoformat(),
        }
        for r in pagination.items
    ]
    
    return jsonify({
        "records": records,
        "page": page,
        "per_page": per_page,
        "total": pagination.total,
        "pages": pagination.pages,
    })


@bp.route("/status-records/<uuid:record_id>", methods=["GET"])
def get_status_record(record_id):
    """Get a specific status record."""
    record = db.get(StatusRecord, record_id)
    
    if not record:
        return jsonify({"error": "Record not found"}), 404
    
    response = {
        "id": str(record.id),
        "project_name": record.project_name,
        "short_name": record.short_name,
        "status": record.status,
        "phase": record.phase,
        "summary": record.summary,
        "reason": record.reason,
        "details": record.details,
        "tags": record.tags,
        "source": record.source,
        "created_at": record.created_at.isoformat(),
        "updated_at": record.updated_at.isoformat(),
    }
    return jsonify(response)


@bp.route("/status-records/<uuid:record_id>", methods=["PATCH"])
def update_status_record(record_id):
    """Update a status record (partial update)."""
    record = db.get(StatusRecord, record_id)
    
    if not record:
        return jsonify({"error": "Record not found"}), 404
    
    data = request.get_json()
    if not data:
        return jsonify({"error": "Request body required"}), 400
    
    updatable_fields = ["project_name", "short_name", "status", "phase", "summary", "reason", "details", "tags"]
    valid_statuses = ["active", "paused", "blocked", "working", "error", "stopped", "completed"]
    
    for field in updatable_fields:
        if field in data:
            if field == "status" and data[field] not in valid_statuses:
                return jsonify({"error": f"Invalid status. Must be one of: {valid_statuses}"}), 400
            setattr(record, field, data[field])
    
    db.commit()
    
    response = {
        "id": str(record.id),
        "project_name": record.project_name,
        "short_name": record.short_name,
        "status": record.status,
        "phase": record.phase,
        "summary": record.summary,
        "reason": record.reason,
        "details": record.details,
        "tags": record.tags,
        "source": record.source,
        "created_at": record.created_at.isoformat(),
        "updated_at": record.updated_at.isoformat(),
    }
    return jsonify(response)


@bp.route("/status-records/<uuid:record_id>", methods=["DELETE"])
def delete_status_record(record_id):
    """Delete a status record."""
    record = db.get(StatusRecord, record_id)
    
    if not record:
        return jsonify({"error": "Record not found"}), 404
    
    db.delete(record)
    db.commit()
    
    return jsonify({"message": "Record deleted"})

