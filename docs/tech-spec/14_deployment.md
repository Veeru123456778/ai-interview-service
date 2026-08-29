# Deployment Architecture

**Project:** AI Interview Service

**Document:** `docs/tech-spec/14_deployment.md`

**Version:** 1.0 (LOCKED)

---

# Ownership

**This document owns:**

- Local development environment.
- Docker configuration.
- Environment variables.
- Production deployment architecture.
- CI/CD overview.
- Deployment design rules.

**References**

- `01_tech_stack.md` → Technology stack.
- `02_architecture.md` → High-Level Architecture.
- `03_folder_structure.md` → Project structure.
- `11_auth_security.md` → Secrets management.
- `13_observability.md` → Health checks and monitoring.

---

# 1. Purpose

This document defines how the AI Interview Service is deployed in development and production.

The backend runs as a containerized Go application connected to PostgreSQL, Redis, Supabase Authentication, and Gemini.

---

# 2. Deployment Architecture

```text
                Internet
                    │
                    ▼
            Next.js Frontend
                    │
          HTTPS / WebSocket
                    │
                    ▼
          Go + Gin Backend Container
          ├── REST APIs
          ├── WebSocket Manager
          ├── Interview Engine
          └── Resume Pipeline
              │
      ┌───────┼────────┐
      ▼       ▼        ▼
 PostgreSQL  Redis   Gemini API
                  (Google AI Studio)

 Authentication
      │
      ▼
 Supabase Auth
```

---

# 3. Runtime Components

| Component | Responsibility |
|-----------|----------------|
| Go Backend | REST APIs, WebSockets, Interview Engine, Resume Pipeline. |
| PostgreSQL | Persistent application data. |
| Redis | Runtime interview state and candidate memory. |
| Supabase Auth | User authentication and JWT issuance. |
| Gemini API | LLM prompt execution through LangChain. |

The frontend is deployed independently.

---

# 4. Local Development Setup

## Required Services

| Service | Runs Locally |
|----------|--------------|
| Go Backend | ✅ |
| PostgreSQL | ✅ Docker |
| Redis | ✅ Docker |
| Supabase | Cloud |
| Gemini | Cloud |

## Local Startup Flow

1. Start PostgreSQL container.
2. Start Redis container.
3. Configure `.env`.
4. Run database migrations.
5. Start Go server.
6. Connect frontend.

---

# 5. Docker Compose

Docker Compose is used only for local development.

## Services

| Service | Port |
|----------|------|
| Backend | 8080 |
| PostgreSQL | 5432 |
| Redis | 6379 |

Supabase and Gemini remain external services.

### Startup Order

1. PostgreSQL.
2. Redis.
3. Backend.

Backend waits until PostgreSQL and Redis become healthy.

---

# 6. Dockerfile

The backend is packaged as a single Go container.

## Build Steps

1. Download Go dependencies.
2. Compile backend binary.
3. Copy binary into lightweight runtime image.
4. Expose backend port.
5. Start application.

### Design Rules

- Multi-stage Docker build.
- Small runtime image.
- No development tooling in runtime image.

---

# 7. Environment Variables

The backend reads configuration from environment variables.

## Required Variables

| Variable | Purpose |
|----------|---------|
| `APP_ENV` | Environment name (`development` or `production`). |
| `PORT` | Backend HTTP port. |
| `DATABASE_URL` | PostgreSQL connection string. |
| `REDIS_URL` | Redis connection string. |
| `SUPABASE_URL` | Supabase project URL. |
| `SUPABASE_ANON_KEY` | Public Supabase key. |
| `SUPABASE_JWT_SECRET` | JWT verification secret. |
| `GEMINI_API_KEY` | Gemini API key. |
| `GEMINI_MODEL` | Gemini model name. |

## Example `.env.example`

```env
APP_ENV=development

PORT=8080

DATABASE_URL=postgres://user:password@localhost:5432/ai_interview

REDIS_URL=redis://localhost:6379

SUPABASE_URL=https://your-project.supabase.co

SUPABASE_ANON_KEY=your_supabase_anon_key

SUPABASE_JWT_SECRET=your_supabase_jwt_secret

GEMINI_API_KEY=your_gemini_api_key

GEMINI_MODEL=gemini-2.5-flash-lite
```

`.env.example` contains placeholders only.

---

# 8. Configuration Files

## `configs/application.yaml`

Stores non-secret application configuration.

### Examples

- API timeout.
- Redis TTL.
- Upload size limit.
- Heartbeat interval.
- Retry counts.

## `configs/prompt_versions.yaml`

Maps logical prompt names to prompt versions.

Example:

```yaml
technical_question: technical_question_v1
technical_evaluator: technical_evaluator_v1
resume_parser: resume_parser_v1
```

Prompt version changes do not require code changes.

---

# 9. Database Migration Strategy

Database schema changes use versioned SQL migrations.

## Migration Rules

- One migration per schema change.
- Forward-only migrations.
- Applied before backend startup in production.

## Migration Order

1. Users.
2. Resumes.
3. Interview Sessions.
4. Topic Evaluations.
5. Final Evaluations.
6. Conversation Messages.

---

# 10. Production Deployment

Production deployment consists of independent services.

| Service | Deployment |
|----------|------------|
| Backend | Docker container. |
| PostgreSQL | Managed PostgreSQL instance. |
| Redis | Managed Redis instance. |
| Supabase Auth | Hosted Supabase. |
| Gemini | Google AI Studio API. |

The backend remains stateless.

Redis and PostgreSQL store application state.

---

# 11. Health Checks

Health endpoints are exposed for deployment platforms.

| Endpoint | Purpose |
|----------|---------|
| `/health/live` | Liveness check. |
| `/health/ready` | Readiness check. |
| `/health` | Overall application status. |

Deployment platform uses readiness before routing traffic.

---

# 12. Deployment Pipeline (CI/CD)

Deployment pipeline runs automatically after successful builds.

## Pipeline Stages

```text
Git Push

   │

Run Tests

   │

Build Docker Image

   │

Run Migrations

   │

Deploy Backend

   │

Health Check

   ▼

Production Ready
```

### Deployment Rules

- Deployment stops if tests fail.
- Deployment stops if migrations fail.
- Deployment stops if readiness check fails.

---

# 13. Runtime Scaling

The backend is designed to be horizontally scalable.

## Stateless Components

- REST handlers.
- WebSocket manager.
- Interview Engine.

## Shared State

| Storage | Shared State |
|---------|--------------|
| PostgreSQL | Persistent data. |
| Redis | Runtime interview state. |

Any backend instance can recover an interview session from Redis.

---

# 14. Backup Strategy

## PostgreSQL

- Regular automated backups.
- Point-in-time recovery handled by managed database provider.

## Redis

- Runtime state only.
- No backup required.
- Expired or completed sessions are recreated from PostgreSQL when needed.

---

# 15. Deployment Design Rules

| Rule | Description |
|------|-------------|
| Backend is stateless. | Runtime state lives in Redis. |
| PostgreSQL is the source of truth. | Permanent interview data only. |
| Redis stores temporary runtime state. | Runtime recovery only. |
| Supabase manages authentication. | Backend verifies JWTs only. |
| Secrets come from environment variables. | Never commit secrets to Git. |
| Docker Compose is for local development only. | Production services are managed independently. |

---

# 16. Related Documents

| Topic | Document |
|-------|----------|
| Tech Stack | `01_tech_stack.md` |
| Architecture | `02_architecture.md` |
| Folder Structure | `03_folder_structure.md` |
| Redis Strategy | `08_redis_strategy.md` |
| Auth & Security | `11_auth_security.md` |
| Observability | `13_observability.md` |