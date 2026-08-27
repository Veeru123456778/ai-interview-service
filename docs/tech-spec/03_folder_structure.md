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
- `04_interview_engine.md` → Interview Engine workflow.
- `05_resume_pipeline.md` → Resume processing pipeline.
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
├── migrations/
├── scripts/
├── tests/
│   ├── integration/
│   └── fixtures/
│
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .env.example
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
    │   ├── builder.go
    │   ├── loader.go
    │   ├── registry.go
    │   ├── resume_parser_v1.txt
    │   ├── intent_detector_v1.txt
    │   ├── guardrail_detector_v1.txt
    │   ├── technical_evaluator_v1.txt
    │   ├── followup_generator_v1.txt
    │   ├── clarification_generator_v1.txt
    │   ├── hint_generator_v1.txt
    │   ├── thinking_prompt_v1.txt
    │   ├── topic_transition_v1.txt
    │   ├── behavioral_generator_v1.txt
    │   └── final_evaluation_v1.txt
    ├── provider/
    ├── memory/
    ├── evaluator/
    ├── schemas/
    └── utils/
```

### Engine Package Responsibilities

| Package | Responsibility |
|---------|----------------|
| `graph` | LangGraph workflow definition and routing. |
| `nodes` | Individual LangGraph nodes (`DetectIntent`, `EvaluateAnswer`, `GenerateFollowUp`, etc.). |
| `prompts` | Prompt builder, prompt loader, prompt registry, and versioned prompt templates. |
| `provider` | LangChain client, Gemini integration, and structured LLM execution. |
| `memory` | Candidate memory management and conversation summaries. |
| `evaluator` | Topic-level and final interview evaluation aggregation. |
| `schemas` | Typed request/response schemas for prompt outputs. |
| `utils` | Engine-specific helper utilities. |

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

### Resume Package Responsibilities

| File | Responsibility |
|------|----------------|
| `extractor.go` | Extract raw text from uploaded PDF. |
| `normalizer.go` | Clean and normalize extracted text. |
| `parser.go` | Execute Resume Parser prompt and validate output. |
| `validator.go` | Validate parsed resume JSON. |
| `intelligence.go` | Build Technology Graph and Interview Topics. |

---

# 5. Shared Package Structure

```text
internal/shared/
│
├── constants/
├── errors/
├── logger/
├── response/
├── utils/
└── validator/
```

### Shared Package Responsibilities

| Package | Used For |
|---------|----------|
| `constants` | Application constants and enums. |
| `errors` | Shared error definitions. |
| `logger` | Zap logger wrapper. |
| `response` | Standard API response formatters. |
| `validator` | Custom validation helpers. |
| `utils` | Generic reusable utilities. |

Business logic never belongs in `shared`.

---

# 6. Storage Package

```text
internal/storage/
│
├── postgres.go
└── redis.go
```

### Responsibilities

| File | Responsibility |
|------|----------------|
| `postgres.go` | PostgreSQL client initialization and configuration. |
| `redis.go` | Redis client initialization and configuration. |

Repositories consume these clients.

---

# 7. WebSocket Package

```text
internal/websocket/
│
├── handler.go
├── manager.go
├── client.go
├── events.go
└── heartbeat.go
```

### Responsibilities

| File | Responsibility |
|------|----------------|
| `handler.go` | Registers WebSocket endpoint. |
| `manager.go` | Connection lifecycle and message routing. |
| `client.go` | Connected client abstraction. |
| `events.go` | Incoming and outgoing event definitions. |
| `heartbeat.go` | Ping/Pong keepalive handling. |

This package contains networking only. Interview logic stays inside `interview/engine`.

---

# 8. Prompt Organization

All prompt-related code and prompt templates live inside the Interview Engine.

```text
internal/interview/engine/prompts/
│
├── builder.go
├── loader.go
├── registry.go
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

### Design Rules

- `builder.go` builds runtime prompt context.
- `loader.go` loads embedded prompt templates using `go:embed`.
- `registry.go` maps prompt names to prompt versions.
- Each prompt template is independently versioned (`*_v1.txt`).

This is the **only prompt location** in the repository.

---

# 9. Dependency Direction

The project follows one-way dependency flow.

```text
Handler
   │
   ▼
Service
   │
   ▼
Repository
   │
   ▼
Storage
```

The Interview Engine is consumed only from the **service layer**.

---

# 10. Package Rules

| Rule | Description |
|------|-------------|
| Domain-first architecture | Every business domain owns its package. |
| One package, one responsibility | No mixed business logic. |
| No circular dependencies | Dependencies flow downward only. |
| `shared` contains reusable utilities only | No business rules. |
| Prompt templates stay inside `engine/prompts` | Single source of truth for prompts. |

---

# 11. Related Documents

| Topic | Document |
|-------|----------|
| High-Level Architecture | `02_architecture.md` |
| Interview Workflow | `04_interview_engine.md` |
| Resume Pipeline | `05_resume_pipeline.md` |
| Prompt Architecture | `06_prompt_architecture.md` |