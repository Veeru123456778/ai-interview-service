# API Specification

**Project:** AI Interview Service

**Document:** `docs/tech-spec/09_api_spec.md`

**Version:** 1.0 (LOCKED)

---

# Ownership

**This document owns:**

- REST API contracts.
- Request and response schemas.
- Authentication requirements.
- Error response format.
- API versioning conventions.

**References**

- `04_interview_engine.md` → Interview session lifecycle.
- `05_resume_pipeline.md` → Resume upload pipeline.
- `07_database_schema.md` → Persistent storage model.
- `08_redis_strategy.md` → Runtime interview state.
- `10_websocket_protocol.md` → Real-time interview communication.

---

# 1. Purpose

This document defines every REST endpoint exposed by the AI Interview Service.

REST APIs are responsible for:

- Authentication-protected CRUD operations.
- Resume upload and retrieval.
- Interview session lifecycle.
- Evaluation retrieval.
- User dashboard and interview history.

Real-time interview communication happens only through WebSocket and is defined in `10_websocket_protocol.md`.

---

# 2. API Conventions

## Base URL

```text
/api/v1
```

## Authentication

Every endpoint except `/health` requires a valid Supabase JWT.

```http
Authorization: Bearer <jwt-token>
```

The backend extracts the authenticated user from the JWT and never accepts `user_id` in request bodies.

---

## Standard Success Response

Every successful response follows the same structure.

```json
{
  "success": true,
  "data": {}
}
```

---

## Standard Error Response

Every failed response follows the same structure.

```json
{
  "success": false,
  "error": {
    "code": "RESUME_VALIDATION_FAILED",
    "message": "Resume must contain at least one skill."
  }
}
```

---

## HTTP Status Codes

| Status | Usage |
|--------|-------|
| 200 | Successful GET request. |
| 201 | Resource created successfully. |
| 204 | Resource deleted successfully. |
| 400 | Validation error. |
| 401 | Authentication failed. |
| 403 | Unauthorized resource access. |
| 404 | Resource not found. |
| 409 | Resource conflict. |
| 422 | Resume parsing or interview validation failed. |
| 500 | Internal server error. |

---

# 3. Authentication Rules

| Rule | Description |
|------|-------------|
| JWT Required | All protected APIs require Supabase JWT. |
| User Ownership | Every resource belongs to authenticated user. |
| Resume Ownership | Users can access only their own resumes. |
| Interview Ownership | Users can access only their own interview sessions. |

The backend derives `user_id` from the JWT.

---

# 4. Health API

## GET `/health`

Returns service health information.

### Authentication

Not required.

### Response — 200

```json
{
  "success": true,
  "data": {
    "status": "UP",
    "service": "ai-interview-service",
    "version": "1.0.0"
  }
}
```

Used by Docker, Kubernetes, and load balancers.

---

# 5. Resume APIs

Resume APIs generate and retrieve Resume Intelligence.

---

## POST `/api/v1/resumes`

Upload a resume and generate Resume Intelligence.

### Authentication

Required.

### Request

**Content-Type**

```text
multipart/form-data
```

| Field | Type | Required |
|------|------|----------|
| resume | PDF File | Yes |

### Success Response — 201

```json
{
  "success": true,
  "data": {
    "resume_id": "resume-uuid",
    "resume_name": "Backend_Resume.pdf",
    "processing_status": "COMPLETED",
    "uploaded_at": "2026-08-28T10:30:00Z"
  }
}
```

### Validation Errors

| Condition | Error Code |
|-----------|------------|
| Missing file | `RESUME_REQUIRED` |
| Invalid format | `INVALID_RESUME_FORMAT` |
| Parser validation failed | `RESUME_VALIDATION_FAILED` |

### Runtime Flow

1. Upload PDF.
2. Extract text.
3. Normalize resume.
4. Generate Resume Intelligence.
5. Store Resume Intelligence in PostgreSQL.
6. Return `resume_id`.

Redis is not used during resume upload.

---

## GET `/api/v1/resumes`

Returns all resumes uploaded by the authenticated user.

### Authentication

Required.

### Response — 200

```json
{
  "success": true,
  "data": {
    "resumes": [
      {
        "resume_id": "resume-uuid",
        "resume_name": "Backend_Resume.pdf",
        "uploaded_at": "2026-08-28T10:30:00Z",
        "processing_status": "COMPLETED"
      }
    ]
  }
}
```

---

## GET `/api/v1/resumes/{resumeId}`

Returns Resume Intelligence metadata.

### Authentication

Required.

### Response — 200

```json
{
  "success": true,
  "data": {
    "resume_id": "resume-uuid",
    "resume_name": "Backend_Resume.pdf",

    "technology_graph": [
      {
        "topic_id": "redis-topic-id",
        "name": "Redis",
        "category": "Database"
      },
      {
        "topic_id": "websocket-topic-id",
        "name": "WebSocket",
        "category": "Backend"
      }
    ],

    "interview_contexts": [
      {
        "context_id": "ctx-project-1",
        "context_name": "AI Interview Platform",
        "context_type": "PROJECT",
        "priority": 1,
        "topic_ids": [
          "redis-topic-id",
          "websocket-topic-id"
        ]
      }
    ]
  }
}
```

### Purpose

Used by frontend Resume Details page.

No runtime interview state is returned.

---

## DELETE `/api/v1/resumes/{resumeId}`

Deletes a resume owned by the authenticated user.

### Authentication

Required.

### Response — 204

Empty response body.

### Delete Rules

- Deletes Resume Intelligence.
- Deletes interview contexts linked to the resume.
- Does not delete completed interview sessions.

---

# 6. Interview Session APIs

Interview session APIs create, retrieve, and complete interview sessions.

Runtime conversation happens over WebSocket after session creation.

---

## POST `/api/v1/interviews`

Creates a new interview session.

### Authentication

Required.

### Request

```json
{
  "resume_id": "resume-uuid"
}
```

### Success Response — 201

```json
{
  "success": true,
  "data": {
    "session_id": "session-uuid",
    "resume_id": "resume-uuid",
    "status": "ACTIVE",
    "websocket_url": "/api/v1/ws/interviews/session-uuid",
    "created_at": "2026-08-28T11:00:00Z"
  }
}
```

### Runtime Flow

1. Validate resume ownership.
2. Load Resume Intelligence from PostgreSQL.
3. Create interview session.
4. Initialize Redis runtime state.
5. Return session information.

No interview question is returned through REST.

---

## GET `/api/v1/interviews/{sessionId}`

Returns interview session metadata.

### Authentication

Required.

### Response — 200

```json
{
  "success": true,
  "data": {
    "session_id": "session-uuid",
    "resume_id": "resume-uuid",
    "status": "ACTIVE",
    "started_at": "2026-08-28T11:00:00Z",
    "completed_at": null
  }
}
```

### Purpose

Used before reconnecting to an existing interview session.

Redis runtime state is not returned through REST.

---

## POST `/api/v1/interviews/{sessionId}/complete`

Marks an interview as completed.

### Authentication

Required.

### Request

Empty body.

### Success Response — 200

```json
{
  "success": true,
  "data": {
    "session_id": "session-uuid",
    "status": "COMPLETED",
    "overall_score": 84,
    "recommendation": "STRONG_HIRE"
  }
}
```

### Completion Flow

1. Finish Interview Engine workflow.
2. Generate final evaluation.
3. Persist transcript.
4. Persist topic evaluations.
5. Persist final evaluation.
6. Delete Redis runtime state.
7. Return evaluation summary.

The complete report is fetched using Evaluation APIs.


---

# 7. Evaluation APIs

Evaluation APIs return interview results after an interview has been completed.

Evaluation data is read from PostgreSQL only.

Redis is never used for completed interview evaluations.

---

## GET `/api/v1/interviews/{sessionId}/evaluation`

Returns the complete final evaluation report.

### Authentication

Required.

### Success Response — 200

```json
{
  "success": true,
  "data": {
    "session_id": "session-uuid",
    "overall_score": 84,
    "recommendation": "STRONG_HIRE",

    "strengths": [
      "Strong understanding of Redis session recovery.",
      "Good reasoning during debugging scenarios."
    ],

    "weaknesses": [
      "Limited knowledge of cache eviction strategies."
    ],

    "learning_recommendations": [
      "Redis persistence strategies.",
      "Distributed caching patterns."
    ],

    "completed_at": "2026-08-28T12:10:00Z"
  }
}
```

### Purpose

Used by the Interview Results page immediately after interview completion.

---

## GET `/api/v1/interviews/{sessionId}/topics`

Returns topic-wise interview evaluation.

### Authentication

Required.

### Success Response — 200

```json
{
  "success": true,
  "data": {
    "topics": [
      {
        "topic_id": "redis-topic-id",
        "topic_name": "Redis",
        "context_name": "AI Interview Platform",
        "score": 88,
        "difficulty_reached": "HARD",

        "strengths": [
          "Correct recovery strategy."
        ],

        "weaknesses": [
          "Missed cache invalidation edge cases."
        ]
      },
      {
        "topic_id": "websocket-topic-id",
        "topic_name": "WebSocket",
        "context_name": "AI Interview Platform",
        "score": 79,
        "difficulty_reached": "MEDIUM",

        "strengths": [
          "Correct heartbeat lifecycle."
        ],

        "weaknesses": [
          "Reconnect timeout handling."
        ]
      }
    ]
  }
}
```

### Purpose

Used by topic score cards and learning recommendations.

---

# 8. Interview History APIs

Interview History APIs allow candidates to review previous interview sessions.

History is read from PostgreSQL.

---

## GET `/api/v1/interviews/history`

Returns all completed interview sessions for the authenticated user.

### Authentication

Required.

### Success Response — 200

```json
{
  "success": true,
  "data": {
    "interviews": [
      {
        "session_id": "session-uuid",
        "resume_id": "resume-uuid",
        "resume_name": "Backend_Resume.pdf",

        "overall_score": 84,
        "recommendation": "STRONG_HIRE",

        "status": "COMPLETED",
        "duration_minutes": 42,

        "completed_at": "2026-08-28T12:10:00Z"
      },
      {
        "session_id": "session-uuid-2",
        "resume_id": "resume-uuid-2",
        "resume_name": "AI_Resume.pdf",

        "overall_score": 91,
        "recommendation": "HIRE",

        "status": "COMPLETED",
        "duration_minutes": 38,

        "completed_at": "2026-08-22T18:30:00Z"
      }
    ]
  }
}
```

### Purpose

Used by the Dashboard Interview History page.

---

## GET `/api/v1/interviews/history/{sessionId}`

Returns complete details of one completed interview.

### Authentication

Required.

### Success Response — 200

```json
{
  "success": true,
  "data": {
    "session_id": "session-uuid",

    "resume": {
      "resume_id": "resume-uuid",
      "resume_name": "Backend_Resume.pdf"
    },

    "overall_score": 84,
    "recommendation": "STRONG_HIRE",

    "duration_minutes": 42,
    "completed_at": "2026-08-28T12:10:00Z",

    "topic_summary": [
      {
        "topic_name": "Redis",
        "score": 88
      },
      {
        "topic_name": "WebSocket",
        "score": 79
      }
    ]
  }
}
```

### Purpose

Used by the Interview Details page before loading detailed evaluation APIs.

---

# 9. Profile APIs

Profile APIs power the candidate dashboard.

---

## GET `/api/v1/profile`

Returns authenticated user profile.

### Authentication

Required.

### Success Response — 200

```json
{
  "success": true,
  "data": {
    "user_id": "user-uuid",
    "name": "Varun Kumar",
    "email": "varun@example.com",
    "created_at": "2026-08-01T09:30:00Z"
  }
}
```

### Design Rule

The backend derives user identity from the JWT.

---

## GET `/api/v1/profile/stats`

Returns dashboard statistics.

### Authentication

Required.

### Success Response — 200

```json
{
  "success": true,
  "data": {
    "total_interviews": 12,
    "average_score": 82,
    "best_score": 91,

    "completed_interviews": 11,
    "active_interviews": 1
  }
}
```

### Purpose

Used by dashboard summary cards.

---

# 10. Error Codes

Every API returns a structured error object.

## Error Response Format

```json
{
  "success": false,
  "error": {
    "code": "INTERVIEW_NOT_FOUND",
    "message": "Interview session does not exist."
  }
}
```

## Error Catalog

| Error Code | HTTP Status | Description |
|------------|------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid JWT. |
| `FORBIDDEN` | 403 | Resource belongs to another user. |
| `RESUME_REQUIRED` | 400 | Resume file missing. |
| `INVALID_RESUME_FORMAT` | 400 | Uploaded file is not a PDF. |
| `RESUME_VALIDATION_FAILED` | 422 | Resume parsing validation failed. |
| `RESUME_NOT_FOUND` | 404 | Resume does not exist. |
| `INTERVIEW_NOT_FOUND` | 404 | Interview session not found. |
| `INTERVIEW_ALREADY_COMPLETED` | 409 | Interview is already completed. |
| `INTERVIEW_ALREADY_ACTIVE` | 409 | Session is already active. |
| `INTERNAL_ERROR` | 500 | Unexpected server error. |

### Design Rules

- Error codes are stable.
- Frontend uses `code` for UI behavior.
- Messages are human-readable.

---

# 11. API Design Rules

| Rule | Description |
|------|-------------|
| REST APIs never carry runtime interview messages. | Runtime communication happens through WebSocket only. |
| JWT authentication protects every user resource. | User ownership is validated server-side. |
| Resume upload generates Resume Intelligence once. | Resume Intelligence is immutable after processing. |
| Interview creation initializes Redis runtime state. | Runtime interview state is never returned through REST. |
| Interview completion persists data before Redis cleanup. | Redis stores runtime state only. |
| Evaluation APIs read from PostgreSQL only. | Completed interview reports never depend on Redis. |
| Standard success/error envelope is used everywhere. | Consistent frontend integration. |

---

# 12. API Lifecycle

| Stage | REST API | WebSocket |
|-------|----------|-----------|
| Upload Resume | `POST /resumes` | — |
| View Resume | `GET /resumes/{resumeId}` | — |
| Create Interview | `POST /interviews` | — |
| Interview Conversation | — | `/ws/interviews/{sessionId}` |
| Finish Interview | `POST /interviews/{sessionId}/complete` | End Interview event |
| View Evaluation | `GET /interviews/{sessionId}/evaluation` | — |
| View Topic Scores | `GET /interviews/{sessionId}/topics` | — |
| View Interview History | `GET /interviews/history` | — |

REST initializes and finalizes interviews.

WebSocket manages the real-time interview conversation.

---

# 13. Related Documents

| Topic | Document |
|-------|----------|
| Interview Engine | `04_interview_engine.md` |
| Resume Pipeline | `05_resume_pipeline.md` |
| Database Schema | `07_database_schema.md` |
| Redis Strategy | `08_redis_strategy.md` |
| WebSocket Protocol | `10_websocket_protocol.md` |
| Auth & Security | `11_auth_security.md` |