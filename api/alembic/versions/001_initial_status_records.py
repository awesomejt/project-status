"""Initial migration - Create status_records table

Revision ID: 001
Revises: 
Create Date: 2026-05-21

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '001'
down_revision = None
branch_labels = None
depends_on = None


def upgrade():
    """Create the status_records table with initial schema."""
    op.create_table(
        'status_records',
        sa.Column('id', sa.String(36), nullable=False, primary_key=True),
        sa.Column('project_name', sa.String(255), nullable=False),
        sa.Column('short_name', sa.String(50), nullable=False, unique=True),
        sa.Column('status', sa.String(20), nullable=False),
        sa.Column('phase', sa.String(50), nullable=True),
        sa.Column('summary', sa.String(500), nullable=True),
        sa.Column('reason', sa.Text(), nullable=True),
        sa.Column('details', sa.Text(), nullable=True),
        sa.Column('tags', sa.ARRAY(sa.String(100)), nullable=True),
        sa.Column('source', sa.String(50), nullable=True),
        sa.Column('created_at', sa.DateTime(), nullable=False),
        sa.Column('updated_at', sa.DateTime(), nullable=False),
    )
    
    # Create indexes for common query patterns
    op.create_index('idx_status_records_short_name', 'status_records', ['short_name'])
    op.create_index('idx_status_records_status', 'status_records', ['status'])
    op.create_index('idx_status_records_phase', 'status_records', ['phase'])
    op.create_index('idx_status_records_created_at', 'status_records', ['created_at'])


def downgrade():
    """Drop the status_records table and its indexes."""
    op.drop_index('idx_status_records_created_at', 'status_records')
    op.drop_index('idx_status_records_phase', 'status_records')
    op.drop_index('idx_status_records_status', 'status_records')
    op.drop_index('idx_status_records_short_name', 'status_records')
    op.drop_table('status_records')
