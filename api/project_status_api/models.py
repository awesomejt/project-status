import json
from datetime import datetime, timezone
from uuid import uuid4
from sqlalchemy import Column, String, Text, DateTime, ARRAY
from sqlalchemy.orm import declarative_base, object_session
from sqlalchemy.ext.declarative import declared_attr

Base = declarative_base()


def get_db_session():
    """Get the current database session."""
    from contextlib import contextmanager
    
    @contextmanager
    def session_provider():
        from . import get_db_session as get_session
        return get_session()()
    
    return session_provider()


class StatusRecord(Base):
    """Status record model."""
    __tablename__ = "status_records"
    
    @declared_attr
    def id(self):
        return Column(String(36), primary_key=True, default=lambda: str(uuid4()))
    
    project_name = Column(String(255), nullable=False)
    short_name = Column(String(50), nullable=False, unique=True)
    status = Column(String(20), nullable=False)
    phase = Column(String(50), nullable=True)
    summary = Column(String(500), nullable=True)
    reason = Column(Text, nullable=True)
    details = Column(Text, nullable=True)
    tags = Column(ARRAY(String(100)), nullable=True)
    source = Column(String(50), nullable=True)
    
    @declared_attr
    def created_at(self):
        return Column(DateTime, nullable=False, default=lambda: datetime.now(timezone.utc))
    
    @declared_attr
    def updated_at(self):
        return Column(DateTime, nullable=False, 
                      default=lambda: datetime.now(timezone.utc),
                      onupdate=lambda: datetime.now(timezone.utc))
    
    def to_dict(self):
        """Convert to dictionary."""
        return {
            "id": str(self.id),
            "project_name": self.project_name,
            "short_name": self.short_name,
            "status": self.status,
            "phase": self.phase,
            "summary": self.summary,
            "reason": self.reason,
            "details": self.details,
            "tags": self.tags,
            "source": self.source,
            "created_at": self.created_at.isoformat() if self.created_at else None,
            "updated_at": self.updated_at.isoformat() if self.updated_at else None,
        }

