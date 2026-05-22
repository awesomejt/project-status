# Integration Test Dockerfile
# Uses a minimal Python image for black-box HTTP testing

FROM python:3.14-slim

WORKDIR /app

COPY tests/integration/test_runner.py tests/integration/test_runner.py

CMD ["python", "tests/integration/test_runner.py"]
