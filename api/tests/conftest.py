import os
import pytest
from flask import Flask


# Must set environment variables BEFORE importing the app module
# Use PostgreSQL 18 for tests. Do not use SQLite.
TEST_DATABASE_URL = "postgresql://project_status:project_status_dev@db:5432/project_status_test"

# Set testing environment
# Note: config.py expects APP_ENV="testing" for TestingConfig
os.environ["APP_ENV"] = "testing"
os.environ["FLASK_ENV"] = "testing"
os.environ["TEST_DATABASE_URL"] = TEST_DATABASE_URL


@pytest.fixture(scope="session")
def app():
    """Create and configure a test Flask application."""
    from project_status_api import create_app
    
    app = create_app(
        DATABASE_URL=TEST_DATABASE_URL,
        APP_ENV="test",
        FLASK_ENV="testing"
    )
    app.config.update(
        TESTING=True,
        SQLALCHEMY_DATABASE_URI=TEST_DATABASE_URL,
        SQLALCHEMY_TRACK_MODIFICATIONS=False,
    )
    
    with app.app_context():
        app.db.create_all()
        yield app
        app.db.drop_all()


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
    with app.app_context():
        app.db.drop_all()
        app.db.create_all()
    yield
    with app.app_context():
        app.db.drop_all()
        app.db.create_all()


@pytest.fixture
def sample_status_record(app):
    """Create a sample status record for testing."""
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
            source="test"
        )
        app.db.session.add(record)
        app.db.session.commit()
        yield record
        app.db.session.delete(record)
        app.db.session.commit()


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
        "tags": ["test", "sample"]
    }


@pytest.fixture
def sample_status_record_update_data():
    """Return sample status record update data as a dictionary."""
    return {
        "project_name": "Updated Test Project",
        "status": "working",
        "summary": "Updated test summary"
    }
