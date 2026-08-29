# Redis Strategy

**Project:** AI Interview Service

**Document:** `docs/tech-spec/08_redis_strategy.md`

**Version:** 2.0 (LOCKED)

---

# Ownership

**This document owns:**

- Runtime interview session state.
- Candidate memory stored in Redis.
- Redis key design.
- Runtime session recovery after reconnect.
- TTL strategy and cleanup rules.

**References**

- `04_interview_engine.md` → Runtime interview workflow.
- `07_database_schema.md` → Persistent storage ownership.
- `09_websocket_protocol.md` → WebSocket connection lifecycle.

---

# 1. Purpose

Redis stores **temporary runtime interview state** while an interview is active.

The Interview Engine reads and updates Redis throughout the interview.

PostgreSQL remains the **persistent source of truth** for resumes, conversation transcripts, topic evaluations, and completed interview sessions.

Redis exists only to make interview progression fast, recoverable, and stateless across backend instances.

---

# 2. What Redis Stores

| Runtime State | Stored in Redis |
|---------------|-----------------|
| Current interview context | ✅ |
| Current topic | ✅ |
| Current scenario | ✅ |
| Current difficulty | ✅ |
| Pending interviewer question | ✅ |
| Candidate memory summary | ✅ |
| Covered topics | ✅ |
| Interview status | ✅ |
| Last activity timestamp | ✅ |
| Active WebSocket session metadata | ✅ |

### Design Rule

Redis stores only the runtime state required to continue an active interview.

Redis **never stores**:

- Resume Intelligence.
- Conversation transcripts.
- Topic evaluations.
- Final interview evaluation.
- Candidate profile information.

Those remain in PostgreSQL.

---

# 3. Redis Key Design

Each interview session owns an isolated runtime namespace.

| Key | Purpose |
|-----|---------|
| `session:{sessionId}` | Complete runtime interview state. |
| `memory:{sessionId}` | Candidate memory summary used during prompt generation. |
| `ws:{sessionId}` | Active WebSocket connection metadata. |

### Design Rules

- Keys are isolated by `sessionId`.
- Runtime keys expire automatically using TTL.
- Runtime keys are deleted immediately after successful interview persistence.

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
  "current_topic_id": "redis-topic-id",

  "current_scenario": "FAILURE",
  "current_difficulty": "MEDIUM",

  "pending_question": {
    "question_id": "question-uuid",
    "question_text": "What happens if Redis becomes unavailable while interview sessions are active?",
    "expected_focus": [
      "session recovery",
      "fallback storage",
      "availability"
    ]
  },

  "covered_topic_ids": [
    "websocket-topic-id"
  ],

  "question_count": 5,

  "status": "ACTIVE",

  "last_activity_at": "2026-08-28T12:30:00Z"
}
```

### Field Responsibilities

| Field | Purpose |
|-------|---------|
| `current_context_id` | Active project, experience, or skill context. |
| `current_topic_id` | Technology currently being discussed. |
| `current_scenario` | Active interview scenario (IMPLEMENTATION, FAILURE, DEBUGGING, SCALING, etc.). |
| `current_difficulty` | Runtime interview difficulty selected by the Interview Engine. |
| `pending_question` | Current interviewer question awaiting the candidate's answer. |
| `covered_topic_ids` | Topics completed during the interview. |
| `question_count` | Number of interviewer questions asked so far. |
| `status` | ACTIVE, COMPLETED, or ABANDONED. |
| `last_activity_at` | Timestamp of the most recent interview activity. |

### Design Rules

- Runtime session state is updated after every successful Interview Engine transition.
- The pending interviewer question remains in Redis until it is evaluated.
- Only the Interview Engine updates runtime session state.

---

# 5. Candidate Memory

**Key**

```text
memory:{sessionId}
```

**Value**

```json
{
  "conversation_summary": "Candidate explained Redis session recovery correctly.",

  "strengths": [
    "Redis persistence",
    "WebSocket lifecycle"
  ],

  "weaknesses": [
    "Cache eviction strategy"
  ],

  "topics_discussed": [
    "Redis",
    "WebSocket"
  ]
}
```

### Purpose

Candidate Memory is a compact runtime summary injected into prompts.

Instead of sending the entire interview transcript to Gemini, Prompt Builder sends this summary along with the current interview context.

### Design Rules

- Memory stores summaries only.
- Full conversation history remains in PostgreSQL.
- Memory is updated after every evaluated answer.
- Memory exists only for the lifetime of an interview session.

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
  "connected_at": "2026-08-28T12:20:10Z",
  "last_ping_at": "2026-08-28T12:31:05Z"
}
```

### Connection Rules

- Only one active WebSocket connection is allowed per interview session.
- A new connection replaces the previous active connection for the same `sessionId`.
- Heartbeats update `last_ping_at`.
- Disconnecting the WebSocket never deletes interview runtime state.

This metadata exists only for reconnect handling.

---

# 7. Memory Update Strategy

Candidate Memory is updated after every evaluated technical answer.

### Runtime Flow

1. Evaluate the candidate's answer.
2. Update strengths.
3. Update weaknesses.
4. Refresh the conversation summary.
5. Update discussed topics if necessary.
6. Persist updated memory in Redis.

### Design Rules

- Candidate Memory is updated **after evaluation**, never before.
- Memory contains only information useful for future prompt generation.
- Previous summaries are replaced instead of growing indefinitely.

---

# 8. Session Recovery

Redis enables interviews to continue after temporary disconnects.

### Recovery Flow

1. Load `session:{sessionId}`.
2. Load `memory:{sessionId}`.
3. Restore the current context, topic, scenario, and difficulty.
4. Restore the pending interviewer question.
5. Continue waiting for the candidate's answer.

### Recovery Rules

- Gemini does **not** regenerate the pending question.
- Candidate Memory is reused immediately.
- Runtime state remains unchanged until the pending answer is evaluated.
- The interview resumes exactly where it stopped.

### Recovery Scenarios

| Event | Behavior |
|-------|----------|
| Browser refresh | Restore runtime state and pending question. |
| Temporary network disconnect | Continue interview after reconnect. |
| Backend restart with Redis available | Runtime state is restored from Redis. |
| WebSocket reconnect | Replace stale WebSocket metadata only. |

---

# 9. TTL Strategy

| Key | TTL |
|-----|-----|
| `session:{sessionId}` | 24 hours |
| `memory:{sessionId}` | 24 hours |
| `ws:{sessionId}` | 30 minutes |

### TTL Refresh Rules

TTL is refreshed whenever:

- A candidate message is received.
- A new interviewer question is generated.
- Candidate Memory is updated.
- A WebSocket heartbeat is received.

### Completion Rules

Completed interviews delete Redis keys immediately after successful persistence.

Abandoned interviews expire automatically through TTL.

---

# 10. Cleanup Strategy

### Successful Interview Completion

1. Persist conversation transcript to PostgreSQL.
2. Persist topic evaluations.
3. Persist final interview evaluation.
4. Delete Redis runtime keys.
5. Close the active WebSocket connection.

### Failed Persistence

If PostgreSQL persistence fails:

- Redis keys remain available.
- Runtime interview state remains recoverable.
- Cleanup is retried after successful persistence.

### Abandoned Interviews

If an interview remains inactive beyond TTL:

- Redis automatically expires runtime state.
- PostgreSQL remains unchanged.
- Application logic may later mark the interview session as `ABANDONED`.

---

# 11. Runtime Design Rules

| Rule | Description |
|------|-------------|
| PostgreSQL is the source of truth. | Redis stores runtime cache only. |
| Resume Intelligence is immutable. | Runtime interview state is mutable. |
| Candidate Memory stores summaries only. | Full transcripts stay in PostgreSQL. |
| One active WebSocket connection per session. | Reconnect replaces stale connections. |
| Pending questions remain in Redis until answered. | Prevent duplicate Gemini generation after reconnect. |
| Runtime state is session-scoped. | Each interview session has isolated Redis keys. |

---

# 12. Why Redis?

| Requirement | Redis Benefit |
|------------|---------------|
| Fast runtime reads and writes | In-memory access for Interview Engine state. |
| WebSocket reconnect | Restore session state instantly. |
| Candidate Memory | Avoid sending the full transcript to Gemini. |
| Runtime interview progression | Store mutable interview state independently of PostgreSQL. |
| Automatic cleanup | TTL removes abandoned runtime sessions. |

Redis allows Interview Engine instances to remain stateless while preserving interview continuity.

---

# 13. Related Documents

| Topic | Document |
|-------|----------|
| Interview Engine | `04_interview_engine.md` |
| Prompt Architecture | `06_prompt_architecture.md` |
| Database Schema | `07_database_schema.md` |
| WebSocket Protocol | `10_websocket_protocol.md` |