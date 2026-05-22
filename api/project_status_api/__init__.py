import os

from flask import Flask, jsonify
from sqlalchemy import create_engine, text
from sqlalchemy.orm import scoped_session, sessionmaker

from .config import config

engine = None
DBSession = None


def get_engine():
    """Get or create database engine."""
    global engine
    if engine is None:
        engine = create_engine(config[os.environ.get("APP_ENV", "development")].SQLALCHEMY_DATABASE_URI)
    return engine


def get_db_session_maker():
    """Get or create database session maker."""
    global DBSession
    if DBSession is None:
        DBSession = scoped_session(sessionmaker(bind=get_engine(), expire_on_commit=False))
    return DBSession


# Initialize db session maker (leads to engine creation on first use)
db = get_db_session_maker()


def init_db():
    """Initialize database tables."""
    from .models import Base

    engine = get_engine()
    Base.metadata.create_all(bind=engine)


def create_app(config_name=None):
    """Application factory."""
    if config_name is None:
        config_name = os.environ.get("FLASK_CONFIG", "development")

    app = Flask(__name__)
    app.config.from_object(config[config_name])

    # Initialize database
    init_db()

    # Register blueprints
    from . import api

    app.register_blueprint(api.bp, url_prefix="/api/project/status")
    # Temporary compatibility prefix while clients migrate to /api/project/status.
    app.register_blueprint(api.bp_legacy, url_prefix="/api")

    # Health endpoints
    @app.route("/health")
    def health():
        return jsonify({"status": "healthy", "service": "project-status-api"})

    @app.route("/ready")
    def ready():
        try:
            with get_engine().connect() as conn:
                conn.execute(text("SELECT 1"))
            return jsonify({"status": "ready", "database": "connected"})
        except Exception as e:
            return jsonify({"status": "not-ready", "database": "disconnected", "error": str(e)}), 503

    @app.teardown_appcontext
    def shutdown_session(exception=None):
        db.remove()

    return app
