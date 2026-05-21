.PHONY: install lint format format-check test

install:
	uv sync --all-extras

lint:
	uv run ruff check .

format:
	uv run ruff format .

format-check:
	uv run ruff format --check .

test:
	uv run pytest

test-coverage:
	uv run pytest --cov=project_status_api --cov-report=term-missing

typecheck:
	@if command -v mypy > /dev/null; then \
		uv run mypy .; \
	else \
		echo "mypy not installed, skipping typecheck"; \
	fi

clean:
	rm -rf __pycache__ .mypy_cache .ruff_cache .pytest_cache .coverage htmlcov
