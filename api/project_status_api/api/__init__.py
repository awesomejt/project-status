from flask import Blueprint, jsonify, request

from .. import db
from ..models import StatusRecord
from ..utils import (
    make_error_response,
    validate_json,
    validate_optional_string,
    validate_status,
    validate_string,
    validate_tags,
)

bp = Blueprint("api", __name__)


@bp.route("/", methods=["POST"])
def create_status_record():
    """Create a new status record."""
    is_valid, result, status_code = validate_json(
        request,
        required_fields=["project_name", "short_name", "status"],
        custom_validators={
            "status": validate_status,
            "project_name": lambda v: validate_string(v, "project_name", max_length=255),
            "short_name": lambda v: validate_string(v, "short_name", max_length=50),
            "phase": validate_optional_string,
            "summary": lambda v: validate_optional_string(v, max_length=500),
            "reason": validate_optional_string,
            "details": validate_optional_string,
            "tags": validate_tags,
            "source": lambda v: validate_optional_string(v, max_length=50),
        },
    )

    if not is_valid:
        response, code = make_error_response(result, status_code)
        return jsonify(response), code

    data = result

    tags_value = data.get("tags") if "tags" in data else []

    record = StatusRecord(
        project_name=data["project_name"],
        short_name=data["short_name"],
        status=data["status"],
        phase=data.get("phase"),
        summary=data.get("summary"),
        reason=data.get("reason"),
        details=data.get("details"),
        tags=tags_value,
        source=data.get("source"),
    )

    db.add(record)
    db.commit()

    return jsonify(record.to_dict()), 201


@bp.route("/", methods=["GET"])
def list_status_records():
    """List status records with pagination and filters."""
    page = request.args.get("page", 1, type=int)
    per_page = request.args.get("per_page", 20, type=int)
    status_filter = request.args.get("status")

    # Build query
    query = db.query(StatusRecord)

    if status_filter:
        query = query.filter(StatusRecord.status == status_filter)

    # Get total count
    total = query.count()
    pages = (total + per_page - 1) // per_page if total > 0 else 1

    # Apply pagination
    offset = (page - 1) * per_page
    query = query.order_by(StatusRecord.created_at.desc()).offset(offset).limit(per_page)
    records_list = query.all()

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
        for r in records_list
    ]

    return jsonify(
        {
            "records": records,
            "page": page,
            "per_page": per_page,
            "total": total,
            "pages": pages,
        }
    )


@bp.route("/<uuid:record_id>", methods=["GET"])
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


@bp.route("/<uuid:record_id>", methods=["PATCH"])
def update_status_record(record_id):
    """Update a status record (partial update)."""
    record = db.get(StatusRecord, record_id)

    if not record:
        response, code = make_error_response("Record not found", 404)
        return jsonify(response), code

    is_valid, result, status_code = validate_json(
        request,
        custom_validators={
            "status": validate_status,
            "project_name": lambda v: validate_string(v, "project_name", max_length=255),
            "short_name": lambda v: validate_string(v, "short_name", max_length=50),
            "phase": validate_optional_string,
            "summary": lambda v: validate_optional_string(v, max_length=500),
            "reason": validate_optional_string,
            "details": validate_optional_string,
            "tags": validate_tags,
        },
    )

    if not is_valid:
        response, code = make_error_response(result, status_code)
        return jsonify(response), code

    data = result
    updatable_fields = ["project_name", "short_name", "status", "phase", "summary", "reason", "details", "tags"]

    for field in updatable_fields:
        if field in data:
            setattr(record, field, data[field])

    db.commit()

    return jsonify(record.to_dict())


@bp.route("/<uuid:record_id>", methods=["DELETE"])
def delete_status_record(record_id):
    """Delete a status record."""
    record = db.get(StatusRecord, record_id)

    if not record:
        return jsonify({"error": "Record not found"}), 404

    db.delete(record)
    db.commit()

    return jsonify({"message": "Record deleted"})
