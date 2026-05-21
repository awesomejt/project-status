def validate_json(request, required_fields=None, custom_validators=None):
    """Validate JSON request body.

    Args:
        request: Flask request object
        required_fields: List of required field names
        custom_validators: Dict of field_name -> validator function

    Returns:
        tuple: (is_valid, data_or_error, status_code)
        - If valid: (True, parsed_json, None)
        - If invalid: (False, error_message, status_code)
    """
    try:
        data = request.get_json()
        if data is None:
            return (False, "Request body required", 400)
        if not isinstance(data, dict):
            return (False, "Request body must be a JSON object", 400)
    except Exception:
        return (False, "Invalid JSON in request body", 400)

    if required_fields:
        missing = [f for f in required_fields if f not in data]
        if missing:
            return (False, f"Missing required fields: {', '.join(missing)}", 400)

    if custom_validators:
        for field, validator in custom_validators.items():
            if field in data:
                is_valid, error = validator(data[field])
                if not is_valid:
                    return (False, f"Invalid value for '{field}': {error}", 400)

    return (True, data, None)


def validate_status(value):
    """Validator for status field."""
    valid_statuses = ["active", "paused", "blocked", "working", "error", "stopped", "completed"]
    if value not in valid_statuses:
        return (False, f"must be one of: {', '.join(valid_statuses)}")
    return (True, None)


def validate_string(value, name, max_length=None, allow_empty=True):
    """Validator for string fields."""
    if not isinstance(value, str):
        return (False, "must be a string")
    if not allow_empty and len(value) == 0:
        return (False, f"{name} cannot be empty")
    if max_length and len(value) > max_length:
        return (False, f"{name} exceeds maximum length of {max_length}")
    return (True, None)


def validate_optional_string(value, max_length=None):
    """Validator for optional string fields."""
    if value is None:
        return (True, None)
    if not isinstance(value, str):
        return (False, "must be a string or null")
    if max_length and len(value) > max_length:
        return (False, f"exceeds maximum length of {max_length}")
    return (True, None)


def validate_tags(value):
    """Validator for tags array."""
    if value is None:
        return (True, None)
    if not isinstance(value, list):
        return (False, "must be an array")
    for tag in value:
        if not isinstance(tag, str):
            return (False, "all tags must be strings")
        if len(tag) > 100:
            return (False, "tags must be at most 100 characters")
    return (True, None)


def make_error_response(message, code, details=None):
    """Create a consistent error response.

    Args:
        message: Human-readable error message
        code: HTTP status code
        details: Optional additional details

    Returns:
        tuple: (response_body, status_code)
    """
    response = {
        "error": {
            "code": code,
            "message": message,
        }
    }
    if details:
        response["error"]["details"] = details
    return (response, code)
