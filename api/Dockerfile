FROM python:3.14-slim

WORKDIR /app

# Install uv for dependency management
COPY --from=ghcr.io/astral-sh/uv:latest /uv /bin/uv

# Install dependencies
COPY uv.lock pyproject.toml /app/
RUN uv sync --frozen --no-dev

# Set Python path
ENV PYTHONPATH=/app

EXPOSE 5000

CMD ["flask", "--app", "project_status_api", "run", "--host=0.0.0.0"]

