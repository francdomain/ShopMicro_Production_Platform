# Multi-stage build for Python ML service
FROM python:3.12-slim AS base

# Install system dependencies and security updates
RUN apt-get update && apt-get upgrade -y && \
  apt-get install -y --no-install-recommends \
  build-essential \
  curl \
  && rm -rf /var/lib/apt/lists/* && \
  pip install --upgrade pip && \
  pip install --force-reinstall setuptools wheel

# Create app directory and user
WORKDIR /usr/src/app
RUN useradd --create-home --shell /bin/bash app && \
  chown -R app:app /usr/src/app

# Development stage
FROM base AS development
ENV PATH="/home/app/.local/bin:$PATH"
ENV PYTHONPATH="/usr/src/app:${PYTHONPATH}"
COPY --chown=app:app requirements.txt .
# Force setuptools reinstall to ensure pkg_resources is available
RUN pip install --force-reinstall setuptools wheel && \
  pip install --no-cache-dir -r requirements.txt
COPY --chown=app:app . .
USER app
EXPOSE 3002
CMD ["python", "app.py"]

# Production stage
FROM base AS production
ENV FLASK_ENV=production
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1
ENV PATH="/home/app/.local/bin:$PATH"
ENV PYTHONPATH="/usr/src/app:${PYTHONPATH}"

# Force setuptools reinstall to ensure pkg_resources is available
COPY requirements.txt .
RUN pip install --force-reinstall setuptools wheel && \
  pip install --no-cache-dir -r requirements.txt gunicorn

# Copy full application code so production image has all modules
COPY . .

# Switch to app user
USER app

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3002/health || exit 1

EXPOSE 3002

# Use gunicorn for production
CMD ["python", "-m", "gunicorn", "--bind", "0.0.0.0:3002", "--workers", "2", "--threads", "4", "--worker-class", "gthread", "app:app"]