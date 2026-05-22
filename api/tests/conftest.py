import os

import pytest
from sqlalchemy import text

# Must set environment variables BEFORE importing the app module
# Use PostgreSQL 18 for tests. Do not use SQLite.
TEST_DATABASE_URL = "postgresql://project_status:project_status_dev@db:5432/project_status_test"

# Set testing environment
os.environ["APP_ENV"] = "testing"
os.environ["FLASK_ENV"] = "testing"
os.environ["TEST_DATABASE_URL"] = TEST_DATABASE_URL
os.environ["DATABASE_URL"] = TEST_DATABASE_URL


@pytest.fixture(scope="session")
def app():
    """Create and configure a test Flask application."""
    from project_status_api import create_app, get_engine

    app = create_app(config_name="testing")
    app.config.update(
        TESTING=True,
    )

    with app.app_context():
        with get_engine().connect() as conn:
            conn.execute(text("TRUNCATE TABLE status_records CASCADE"))
            conn.commit()
        yield app
        with get_engine().connect() as conn:
            conn.execute(text("TRUNCATE TABLE status_records CASCADE"))
            conn.commit()


@pytest.fixture
def client(app):
    """Create a test client for the Flask application."""
    return app.test_client()


@pytest.fixture
def runner(app):
    """Create a test CLI runner for the Flask application."""
    return app.test_cli_runner()


@pytest.fixture(autouse=True)
def clean_db(app):
    """Clean database before and after each test."""
    from project_status_api import get_engine

    with app.app_context():
        with get_engine().connect() as conn:
            conn.execute(text("DELETE FROM status_records"))
            conn.commit()
    yield
    with app.app_context():
        with get_engine().connect() as conn:
            conn.execute(text("DELETE FROM status_records"))
            conn.commit()


@pytest.fixture
def sample_status_record(app):
    """Create a sample status record for testing."""
    from project_status_api import db
    from project_status_api.models import StatusRecord

    with app.app_context():
        record = StatusRecord(
            project_name="Test Project",
            short_name="test-proj",
            status="active",
            phase="planning",
            summary="A test status record",
            reason=None,
            details="This is a test record for pytest fixtures",
            tags=["test", "sample"],
            source="test",
        )
        db.add(record)
        db.commit()
        yield record
        db.delete(record)
        db.commit()


@pytest.fixture
def sample_status_record_data():
    """Return sample status record data as a dictionary."""
    return {
        "project_name": "Test Project",
        "short_name": "test-proj",
        "status": "active",
        "phase": "planning",
        "summary": "A test status record",
        "reason": "Test reason",
        "details": "Test details",
        "tags": ["test", "sample"],
    }


@pytest.fixture
def sample_status_record_update_data():
    """Return sample status record update data as a dictionary."""
    return {"project_name": "Updated Test Project", "status": "working", "summary": "Updated test summary"}
