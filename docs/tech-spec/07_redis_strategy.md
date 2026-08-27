# Redis Strategy

**Project:** AI Interview Service

**Document:** `docs/tech-spec/08_redis_strategy.md`

**Version:** 1.0 (LOCKED)

---

# Ownership

**This document owns:**

- Runtime interview session state.
- Candidate memory stored in Redis.
- Redis key design.
- TTL strategy.
- Session recovery after reconnect.

**References**

- `04_interview_engine.md` → Runtime workflow.
- `07_database_schema.md` → Persistent storage.

---

# 1. Purpose

Redis stores **temporary runtime state** while an interview is active.

PostgreSQL stores permanent interview history after the interview is completed.

---

# 2. What Redis Stores

| Runtime State | Stored in Redis |
|---------------|-----------------|
| Current interview context | ✅ |
| Current topic | ✅ |
| Current scenario | ✅ |
| Current difficulty | ✅ |
| Candidate memory summary | ✅ |
| Covered topics | ✅ |
| Active WebSocket session metadata | ✅ |

Redis never stores permanent transcripts or evaluations.

---

# 3. Redis Key Design

| Key | Purpose |
|-----|---------|
| `session:{sessionId}` | Complete runtime interview state. |
| `memory:{sessionId}` | Candidate memory summary. |
| `ws:{sessionId}` | Active WebSocket connection metadata. |

Each interview session has its own isolated keys.

---

# 4. Runtime Session State

**Key**

```text
session:{sessionId}
```

**Value**

```json
{
  "session_id": "session-uuid",

  "current_context_id": "ctx-project-1",
  "current_topic_id": "redis-id",

  "current_scenario": "FAILURE",
  "current_difficulty": "MEDIUM",

  "covered_topic_ids": [
    "websocket-id"
  ],

  "question_count": 5,

  "status": "ACTIVE",

  "last_activity_at": "2026-08-28T12:30:00Z"
}
```

This object is updated after every interview action.

---

# 5. Candidate Memory

**Key**

```text
memory:{sessionId}
```

**Value**

```json
{
  "conversation_summary": "Candidate explained Redis session storage correctly.",

  "strengths": [
    "Redis persistence",
    "WebSocket lifecycle"
  ],

  "weaknesses": [
    "Cache eviction strategy"
  ]
}
```

### Purpose

The Interview Engine sends this summary to Gemini instead of the entire conversation history.

Memory is updated after every evaluated answer.

---

# 6. WebSocket Session Metadata

**Key**

```text
ws:{sessionId}
```

**Value**

```json
{
  "connection_id": "socket-123",
  "connected": true,
  "last_ping_at": "2026-08-28T12:31:05Z"
}
```

Used only for reconnect handling.

---

# 7. Memory Update Strategy

After each answer:

1. Evaluate the answer.
2. Update strengths and weaknesses.
3. Refresh the conversation summary.
4. Save the updated memory in Redis.

The full transcript is **not** rewritten in Redis.

---

# 8. Session Recovery

If the WebSocket disconnects:

1. Load `session:{sessionId}`.
2. Load `memory:{sessionId}`.
3. Resume from the current context and topic.
4. Continue with the pending interview question.

The interview does not restart.

---

# 9. TTL Strategy

| Key | TTL |
|-----|-----|
| `session:{sessionId}` | 24 hours |
| `memory:{sessionId}` | 24 hours |
| `ws:{sessionId}` | 30 minutes |

TTL is refreshed whenever interview activity occurs.

Completed interviews remove Redis keys after persistence.

---

# 10. Cleanup Strategy

Interview completion:

1. Persist evaluations to PostgreSQL.
2. Delete runtime Redis keys.
3. Close WebSocket session.

Abandoned interviews expire automatically using TTL.

---

# 11. Why Redis?

| Requirement | Redis Benefit |
|------------|---------------|
| Fast runtime reads/writes | In-memory access. |
| WebSocket reconnect | Restore session state instantly. |
| Candidate memory | Avoid sending full transcript to Gemini. |
| Automatic cleanup | TTL removes abandoned sessions. |

---

# 12. Related Documents

| Topic | Document |
|-------|----------|
| Interview Engine | `04_interview_engine.md` |
| Database Schema | `07_database_schema.md` |
| Evaluation Engine | `12_evaluation_engine.md` |