# Observability

**Project:** AI Interview Service

**Document:** `docs/tech-spec/13_observability.md`

**Version:** 1.0 (LOCKED)

---

# Ownership

**This document owns:**

- Structured logging strategy.
- Request tracing and correlation IDs.
- Application metrics.
- Health check endpoints.
- Error monitoring.
- Observability design rules.

**References**

- `02_architecture.md` → Service architecture.
- `09_api_spec.md` → REST endpoints.
- `10_websocket_protocol.md` → WebSocket lifecycle.
- `11_auth_security.md` → Authentication middleware.

---

# 1. Purpose

Observability provides visibility into the runtime behavior of the AI Interview Service.

It helps monitor API requests, WebSocket sessions, Redis usage, PostgreSQL health, Gemini requests, and interview execution without affecting business logic.

Observability consists of:

- Structured logs.
- Metrics.
- Health checks.
- Request tracing.

---

# 2. Observability Architecture

```text
Client Request
      │
      ▼
Gin Middleware
      │
      ├── Structured Logger
      ├── Request ID Generator
      ├── Metrics Collector
      ▼
Application Services
      │
      ├── PostgreSQL
      ├── Redis
      └── Gemini
```

Every incoming request receives a unique request identifier.

---

# 3. Structured Logging

The backend uses structured JSON logging.

## Logging Library

- Go logging implementation uses Zap.
- Logs are emitted in JSON format.

## Required Log Fields

| Field | Purpose |
|-------|---------|
| `timestamp` | UTC timestamp. |
| `level` | INFO, WARN, ERROR. |
| `request_id` | Unique request identifier. |
| `user_id` | Authenticated user if available. |
| `session_id` | Interview session if available. |
| `operation` | Business operation name. |
| `duration_ms` | Execution duration. |
| `status` | Success or failure status. |

## Example Log

```json
{
  "timestamp": "2026-08-28T12:30:00Z",
  "level": "INFO",
  "request_id": "req-123",
  "user_id": "user-uuid",
  "session_id": "session-uuid",
  "operation": "CreateInterviewSession",
  "duration_ms": 182,
  "status": "SUCCESS"
}
```

---

# 4. Logging Levels

| Level | Used For |
|-------|----------|
| INFO | Successful requests and interview lifecycle events. |
| WARN | Recoverable failures and retries. |
| ERROR | Unexpected failures requiring investigation. |

## INFO Events

- Resume uploaded.
- Resume processed.
- Interview created.
- Interview completed.
- WebSocket connected.
- WebSocket disconnected.

## WARN Events

- Resume parser retry.
- Gemini timeout retry.
- WebSocket reconnect.
- Redis reconnect.

## ERROR Events

- PostgreSQL unavailable.
- Redis unavailable.
- Gemini request failed after retry.
- Authentication failure.
- Internal server error.

---

# 5. Request Correlation

Every REST request receives a request identifier.

## Request ID Flow

```text
Incoming Request

      │

Generate Request ID

      │

Gin Context

      │

Logger

      │

Repositories / Services
```

## Design Rules

- Generated once per request.
- Passed through service and repository layers.
- Included in all logs generated during request execution.

---

# 6. Interview Session Logging

Interview sessions include session-specific logging.

## Logged Events

| Event | Logged |
|-------|--------|
| Interview Started | ✅ |
| Question Generated | ✅ |
| Candidate Answer Received | ✅ |
| Topic Completed | ✅ |
| Interview Completed | ✅ |
| Session Recovered | ✅ |

## Example

```json
{
  "level": "INFO",
  "request_id": "req-456",
  "session_id": "session-uuid",
  "operation": "QuestionGenerated",
  "topic_id": "redis-topic-id"
}
```

Logs never include full interview answers.

---

# 7. WebSocket Observability

WebSocket lifecycle events are logged.

## Logged Events

| Event | Logged |
|-------|--------|
| Connection Established | ✅ |
| Authentication Success | ✅ |
| Authentication Failure | ✅ |
| Heartbeat Received | ✅ |
| Reconnect | ✅ |
| Connection Closed | ✅ |

## Design Rules

- Log connection lifecycle only.
- Do not log heartbeat payload contents.
- Do not log full conversation messages.

---

# 8. Gemini Request Logging

Gemini calls are monitored without logging prompt contents.

## Logged Fields

| Field | Description |
|-------|-------------|
| Request ID | Correlation ID. |
| Prompt Name | Prompt template executed. |
| Prompt Version | Version from registry. |
| Model | Gemini model name. |
| Latency | Request duration. |
| Retry Count | Retry attempts. |
| Status | Success or failure. |

## Example

```json
{
  "operation": "GeminiRequest",
  "prompt": "technical_question_v1",
  "model": "gemini-2.5-flash-lite",
  "latency_ms": 540,
  "status": "SUCCESS"
}
```

Prompt text and candidate answers are never logged.

---

# 9. Metrics

Application metrics measure runtime health.

## API Metrics

| Metric | Description |
|--------|-------------|
| API Request Count | Total REST requests. |
| API Request Duration | Request latency. |
| API Error Count | Failed REST requests. |

## Interview Metrics

| Metric | Description |
|--------|-------------|
| Active Interviews | Interviews currently in Redis. |
| Completed Interviews | Total completed interviews. |
| Average Interview Duration | Completion duration. |
| Topic Evaluation Count | Evaluated topics. |

## WebSocket Metrics

| Metric | Description |
|--------|-------------|
| Active Connections | Current socket connections. |
| Connection Count | Connections established. |
| Reconnect Count | Successful reconnects. |
| Disconnect Count | Closed connections. |

## Gemini Metrics

| Metric | Description |
|--------|-------------|
| Gemini Request Count | Total prompt executions. |
| Gemini Error Count | Failed requests. |
| Gemini Retry Count | Retry attempts. |
| Average Gemini Latency | Mean response time. |

---

# 10. Health Check Endpoints

Health endpoints verify application dependencies.

## GET `/health`

Returns application health.

### Success Response

```json
{
  "status": "UP"
}
```

---

## GET `/health/ready`

Readiness check verifies dependencies.

### Success Response

```json
{
  "status": "READY",
  "dependencies": {
    "postgres": "UP",
    "redis": "UP"
  }
}
```

### Failure Response

```json
{
  "status": "NOT_READY",
  "dependencies": {
    "postgres": "DOWN",
    "redis": "UP"
  }
}
```

---

## GET `/health/live`

Liveness endpoint verifies the application process.

### Success Response

```json
{
  "status": "ALIVE"
}
```

Liveness does not verify external dependencies.

---

# 11. Error Monitoring

Errors are logged with structured metadata.

## Logged Fields

| Field | Purpose |
|-------|---------|
| Request ID | Correlation. |
| Operation | Failed operation. |
| Error Code | Stable application error. |
| Error Message | Internal diagnostic message. |

## Example

```json
{
  "level": "ERROR",
  "request_id": "req-789",
  "operation": "ResumeProcessing",
  "error_code": "RESUME_VALIDATION_FAILED"
}
```

Sensitive information is never logged.

---

# 12. Sensitive Data Logging Rules

The application must never log sensitive runtime data.

## Never Log

- JWT tokens.
- Resume PDF contents.
- Candidate interview answers.
- Candidate memory summaries.
- Gemini prompt text.
- Gemini API keys.
- Database credentials.
- Redis credentials.

## Safe to Log

- Resume identifier.
- Session identifier.
- Topic identifier.
- Prompt name and version.
- Execution duration.
- Error codes.

---

# 13. Monitoring Design Rules

| Rule | Description |
|------|-------------|
| Logs are structured JSON. | Machine-readable logging. |
| Every request has a request ID. | End-to-end tracing. |
| Session logs include session ID. | Interview traceability. |
| Metrics are aggregated separately from logs. | Lightweight monitoring. |
| Sensitive runtime data is never logged. | Security-first observability. |

---

# 14. Related Documents

| Topic | Document |
|-------|----------|
| Architecture | `02_architecture.md` |
| API Specification | `09_api_spec.md` |
| WebSocket Protocol | `10_websocket_protocol.md` |
| Authentication & Security | `11_auth_security.md` |