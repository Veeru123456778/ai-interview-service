# Folder Structure

**Project:** AI Interview Service

**Document:** `docs/tech-spec/03_folder_structure.md`

**Version:** 1.0 (LOCKED)

---

# Ownership

**This document owns:**

- Repository structure.
- Package responsibilities.
- Dependency direction.
- Internal module boundaries.

**References**

- `02_architecture.md` → Service architecture.
- `04_interview_engine.md` → Engine workflow.
- `05_resume_pipeline.md` → Resume processing.
- `06_prompt_architecture.md` → Prompt contracts.

---

# 1. Repository Structure

```text
ai-interview-service/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── configs/
│   ├── application.yaml
│   └── prompt_versions.yaml
│
├── docs/
│   ├── tech-spec/
│   └── adr/
│
├── internal/
│   ├── app/
│   ├── middleware/
│   ├── shared/
│   ├── auth/
│   ├── user/
│   ├── resume/
│   ├── interview/
│   ├── websocket/
│   ├── storage/
│   └── health/
│
├── prompts/
│
├── migrations/
├── scripts/
├── tests/
│   ├── integration/
│   └── fixtures/
│
├── go.mod
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

---

# 2. Internal Package Structure

```text
internal/
│
├── app/
│
├── middleware/
│
├── shared/
│   ├── constants/
│   ├── errors/
│   ├── logger/
│   ├── response/
│   ├── utils/
│   └── validator/
│
├── auth/
│
├── user/
│
├── resume/
│
├── interview/
│
├── websocket/
│
├── storage/
│
└── health/
```

Each top-level package owns a single business capability.

---

# 3. Interview Module Structure

```text
internal/interview/
│
├── handler.go
├── service.go
├── repository.go
├── model.go
├── dto.go
├── session_manager.go
├── types.go
├── errors.go
│
└── engine/
    ├── graph/
    ├── nodes/
    ├── prompts/
    ├── provider/
    ├── memory/
    ├── evaluator/
    ├── schemas/
    └── utils/
```

### Responsibilities

| Package | Responsibility |
|---------|----------------|
| `graph` | LangGraph workflow definition and routing. |
| `nodes` | Individual interview nodes (intent, evaluation, follow-up, etc.). |
| `prompts` | Prompt builder, loader, registry, and prompt versions. |
| `provider` | LangChain + Gemini integration. |
| `memory` | Candidate memory management. |
| `evaluator` | Topic and final evaluation aggregation. |
| `schemas` | Typed request/response schemas for prompt outputs. |
| `utils` | Engine-only helper utilities. |

---

# 4. Resume Module Structure

```text
internal/resume/
│
├── handler.go
├── service.go
├── repository.go
├── model.go
├── dto.go
├── extractor.go
├── parser.go
├── normalizer.go
├── intelligence.go
└── validator.go
```

### Responsibilities

| File | Responsibility |
|------|----------------|
| `extractor.go` | Extract text from uploaded PDF. |
| `normalizer.go` | Clean and normalize extracted text. |
| `parser.go` | Call Resume Parser prompt and validate JSON. |
| `validator.go` | Validate parsed resume structure. |
| `intelligence.go` | Build Technology Graph and Interview Topics. |

---

# 5. Shared Package

Reusable code shared across modules.

```text
shared/
├── constants/
├── errors/
├── logger/
├── response/
├── utils/
└── validator/
```

| Package | Used For |
|---------|----------|
| `constants` | Application constants and enums. |
| `errors` | Shared error definitions. |
| `logger` | Zap logger wrapper. |
| `response` | Standard API response helpers. |
| `validator` | Custom validation helpers. |
| `utils` | Generic utility functions. |

Business logic never lives here.

---

# 6. Storage Package

```text
internal/storage/
├── postgres.go
└── redis.go
```

| File | Responsibility |
|------|----------------|
| `postgres.go` | PostgreSQL connection and configuration. |
| `redis.go` | Redis client and configuration. |

Repositories consume these clients.

---

# 7. WebSocket Package

```text
internal/websocket/
├── handler.go
├── manager.go
├── client.go
├── events.go
└── heartbeat.go
```

| File | Responsibility |
|------|----------------|
| `handler.go` | WebSocket endpoint registration. |
| `manager.go` | Connection lifecycle and routing. |
| `client.go` | Connected client abstraction. |
| `events.go` | Incoming/outgoing event definitions. |
| `heartbeat.go` | Ping/Pong keepalive handling. |

Interview logic does not live in this package.

---

# 8. Prompt Templates

Prompt templates are versioned assets.

```text
prompts/
├── resume_parser_v1.txt
├── intent_detector_v1.txt
├── guardrail_detector_v1.txt
├── technical_evaluator_v1.txt
├── followup_generator_v1.txt
├── clarification_generator_v1.txt
├── hint_generator_v1.txt
├── thinking_prompt_v1.txt
├── topic_transition_v1.txt
├── behavioral_generator_v1.txt
└── final_evaluation_v1.txt
```

The Interview Engine loads prompts through `engine/prompts`.

---

# 9. Dependency Direction

The project follows one-way dependency flow.

```text
handler
   │
   ▼
service
   │
   ▼
repository
   │
   ▼
storage
```

The Interview Engine is consumed only from the service layer.

---

# 10. Package Rules

| Rule | Description |
|------|-------------|
| Domain-first architecture | Every business domain owns its files. |
| One package, one responsibility | No mixed business logic. |
| No circular dependencies | Packages depend only downward. |
| Shared package is utility-only | No business rules in `shared`. |
| Prompt templates are not Go code | Prompt files remain under `/prompts`. |

---

# 11. Related Documents

| Topic | Document |
|-------|----------|
| High-Level Architecture | `02_architecture.md` |
| Interview Workflow | `04_interview_engine.md` |
| Resume Pipeline | `05_resume_pipeline.md` |
| Prompt Architecture | `06_prompt_architecture.md` |