# Technology Stack

**Project:** AI Interview Service

**Document:** `docs/tech-spec/01_tech_stack.md`

**Version:** 1.0 (LOCKED)

---

# Ownership

**This document owns:**

- Technology selection.
- Why each technology was chosen.
- Version strategy.
- Architectural technology decisions.

**References:**

- `02_architecture.md` → System architecture.
- `03_folder_structure.md` → Package organization.
- `04_interview_engine.md` → LangGraph workflow.
- `05_resume_pipeline.md` → Resume processing.
- `06_prompt_architecture.md` → Prompt contracts.

---

# 1. Backend Technology Stack

| Layer | Technology | Why |
|-------|------------|-----|
| Language | **Go** | Fast, simple concurrency, production backend. |
| Framework | **Gin** | Lightweight HTTP framework with middleware support. |
| ORM | **GORM** | PostgreSQL integration and migrations. |
| API Style | **REST + WebSocket** | REST for APIs, WebSocket for real-time interview. |
| Authentication | Supabase Auth | User authentication and JWT issuance. |
| Authorization | JWT Verification | Verify Supabase-issued JWTs inside Gin middleware. |

---

# 2. AI Stack

| Layer | Technology | Why |
|-------|------------|-----|
| Workflow Engine | **LangGraph** | Deterministic interview workflow orchestration. |
| LLM Abstraction | **LangChain (Go)** | Provider abstraction and structured output parsing. |
| Initial LLM | **Gemini 2.5 Flash** | Fast, low-latency LLM for interview generation, evaluation, and resume intelligence in V1. |

### AI Design Decisions

- LangGraph manages interview state transitions.
- LangChain is used only for LLM interaction.
- Prompt templates are versioned independently.
- Backend owns memory and orchestration.
- The interview engine is provider-agnostic through LangChain-Go. Gemini 2.5 Flash is the default provider for V1 and can be replaced without changing the workflow.

---

# 3. Database Stack

| Purpose | Technology |
|---------|------------|
| Primary Database | PostgreSQL |
| Runtime State | Redis |
| Cache / Session State | Redis |

### Storage Responsibilities

| Storage | Owns |
|---------|------|
| PostgreSQL | Users, resumes, interview sessions, evaluations. |
| Redis | Active interview state, current context, current topic, current scenario, candidate memory summary, and WebSocket session metadata. |
---

# 4. Frontend Integration

| Feature | Technology |
|---------|------------|
| Frontend | Next.js |
| UI Components | ShadCN UI |
| Styling | Tailwind CSS |
| Authentication | Supabase Auth + JWT |
| Speech-to-Text | Web Speech API |
| Text-to-Speech | Web Speech API |

### Speech Decision (V1)

Speech recognition and speech synthesis run entirely in the browser using the Web Speech API. The backend receives only text messages over WebSocket.

---

# 5. Communication Stack

| Communication | Technology |
|--------------|------------|
| Client → Backend | HTTPS REST APIs |
| Interview Session | WebSocket |
| Internal AI Calls | LangChain Client |

---

# 6. Development Stack

| Tool | Purpose |
|------|---------|
| Docker | Local development and deployment. |
| Docker Compose | Local infrastructure. |
| Air | Live reload during development. |
| Make | Common development commands. |

---

# 7. Testing Stack

| Test Type | Technology |
|-----------|------------|
| Unit Tests | Go Testing Package |
| Integration Tests | Go Testing + Test Containers |
| API Tests | Postman / Bruno |
| AI Prompt Tests | JSON schema validation + golden responses |

---

# 8. Observability Stack

| Purpose | Technology |
|---------|------------|
| Logging | Zap |
| Metrics | Prometheus |
| Dashboards | Grafana |
| Tracing | OpenTelemetry |

Detailed design is defined in `13_observability.md`.

---

# 9. Security Stack

| Purpose | Technology |
|---------|------------|
| Authentication | Supabase Auth |
| Authorization | JWT Verification Middleware |
| Environment Secrets | `.env` + Config Loader |
| API Validation | Gin Binding + Custom Validators |

Detailed rules are defined in `11_auth_security.md`.

---

# 10. Version Strategy

| Component | Initial Version |
|-----------|-----------------|
| Go | 1.25.x |
| Gin | Latest stable |
| PostgreSQL | 17 |
| Redis | 8 |
| LangChain Go | Latest stable |
| LangGraph (Go) | Latest stable |
| Gemini | Gemini 2.5 Flash |

Minor dependency updates are allowed during V1 without changing architecture.

---

# 11. Technology Decisions (ADR Summary)

| Decision | ADR |
|----------|-----|
| Go + Gin backend | ADR-001 |
| LangGraph + LangChain | ADR-002 |
| PostgreSQL + Redis | ADR-004 |
| Domain-first architecture | ADR-005 |
| Web Speech API for speech mode | Tech Stack v1.0 |
| No permanent resume PDF storage in V1 | ADR-003 |

---

# 12. Future Upgrades (Out of Scope for V1)

| Area | Possible Upgrade |
|------|-------------------|
| LLM Provider | OpenAI / Claude / Local LLMs via LangChain. |
| Speech | Streaming speech APIs for lower latency. |
| Resume Search | Embeddings + Vector Database. |
| Queue | Redis Streams for asynchronous processing. |
| Object Storage | S3 / Cloudflare R2 if PDF retention is introduced. |