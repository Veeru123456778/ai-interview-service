# System Architecture (HLD)

**Project:** AI Interview Service

**Document:** `docs/tech-spec/02_architecture.md`

**Version:** 1.0 (LOCKED)

---

# Ownership

**This document owns:**

- High-Level Architecture (HLD).
- Service boundaries.
- Communication between backend components.
- Request lifecycle.
- Runtime component responsibilities.

**References:**

- `01_tech_stack.md` → Technology choices.
- `03_folder_structure.md` → Package structure.
- `04_interview_engine.md` → Interview workflow.
- `05_resume_pipeline.md` → Resume processing internals.
- `08_redis_strategy.md` → Redis state design.

---

# 1. High-Level Architecture

```mermaid
%%{init:{
  "theme":"base",
  "themeVariables":{
    "background":"#0F172A",
    "primaryColor":"#1E3A8A",
    "primaryTextColor":"#FFFFFF",
    "primaryBorderColor":"#FFFFFF",
    "secondaryColor":"#2563EB",
    "secondaryTextColor":"#FFFFFF",
    "tertiaryColor":"#1E293B",
    "lineColor":"#FFFFFF",
    "fontFamily":"Inter",
    "edgeLabelBackground":"#0F172A"
  }
}}%%

flowchart TD

    USER["Frontend (Next.js)"]

    WS["WebSocket Manager"]

    API["Interview Service (Go + Gin)"]

    RESUME["Resume Pipeline"]

    ENGINE["Interview Engine (LangGraph)"]

    LLM["LangChain + Gemini"]

    PG["PostgreSQL"]

    REDIS["Redis"]

    USER --> API
    USER --> WS

    WS --> API

    API --> RESUME
    API --> ENGINE

    ENGINE --> LLM

    RESUME --> PG
    API --> PG

    ENGINE --> REDIS
```

---

# 2. Component Responsibilities

| Component | Responsibility |
|-----------|----------------|
| **Frontend** | Upload resume, manage interview UI, speech input/output. |
| **Interview Service** | REST APIs, authentication, session management, WebSocket entry point. |
| **Resume Pipeline** | Parse resume and generate structured interview intelligence. |
| **Interview Engine** | Execute LangGraph workflow and manage interview state. |
| **LangChain + Gemini** | Execute prompts and return structured responses. |
| **PostgreSQL** | Persistent application data. |
| **Redis** | Runtime interview state and caching. |

---

# 3. Runtime Data Ownership

| Storage | Owns |
|---------|------|
| **PostgreSQL** | Users, resumes, interview sessions, evaluations, metadata. |
| **Redis** | Current topic, covered topics, candidate memory, interview timers, active session state. |

Redis never becomes the source of truth.

---

# 4. Backend Request Flow

```mermaid
%%{init:{
  "theme":"base",
  "themeVariables":{
    "background":"#0F172A",
    "primaryColor":"#1E3A8A",
    "primaryTextColor":"#FFFFFF",
    "primaryBorderColor":"#FFFFFF",
    "secondaryColor":"#2563EB",
    "secondaryTextColor":"#FFFFFF",
    "tertiaryColor":"#1E293B",
    "lineColor":"#FFFFFF",
    "fontFamily":"Inter",
    "edgeLabelBackground":"#0F172A"
  }
}}%%

sequenceDiagram

    participant User
    participant API as Interview Service
    participant Resume
    participant Engine
    participant LLM
    participant PG
    participant Redis

    User->>API: Upload Resume
    API->>Resume: Process Resume
    Resume->>LLM: Parse Resume
    LLM-->>Resume: Structured Resume JSON
    Resume->>PG: Store Resume Intelligence

    User->>API: Start Interview
    API->>Engine: Initialize Session
    Engine->>Redis: Create Runtime State
    Engine->>PG: Load Resume Intelligence
    Engine->>LLM: Generate First Question
    LLM-->>Engine: Question
    Engine-->>User: Send Question
```

---

# 5. Interview Communication Flow

```mermaid
%%{init:{
  "theme":"base",
  "themeVariables":{
    "background":"#0F172A",
    "primaryColor":"#1E3A8A",
    "primaryTextColor":"#FFFFFF",
    "primaryBorderColor":"#FFFFFF",
    "secondaryColor":"#2563EB",
    "secondaryTextColor":"#FFFFFF",
    "tertiaryColor":"#1E293B",
    "lineColor":"#FFFFFF",
    "fontFamily":"Inter",
    "edgeLabelBackground":"#0F172A"
  }
}}%%

sequenceDiagram

    participant Frontend
    participant WS as WebSocket Manager
    participant Engine
    participant Redis
    participant LLM

    Frontend->>WS: Candidate message
    WS->>Engine: Forward message

    Engine->>Redis: Load session state
    Engine->>LLM: Execute current prompt
    LLM-->>Engine: Structured response

    Engine->>Redis: Update runtime state
    Engine-->>WS: Interview response
    WS-->>Frontend: Render response
```

Speech is converted to text on the frontend before being sent through WebSocket.

---

# 6. Resume Processing Flow

```mermaid
%%{init:{
  "theme":"base",
  "themeVariables":{
    "background":"#0F172A",
    "primaryColor":"#1E3A8A",
    "primaryTextColor":"#FFFFFF",
    "primaryBorderColor":"#FFFFFF",
    "secondaryColor":"#2563EB",
    "secondaryTextColor":"#FFFFFF",
    "tertiaryColor":"#1E293B",
    "lineColor":"#FFFFFF",
    "fontFamily":"Inter",
    "edgeLabelBackground":"#0F172A"
  }
}}%%

flowchart LR

    PDF["Resume PDF"]

    EXTRACT["Extract Text"]

    PARSE["LLM Resume Parser"]

    BUILD["Resume Intelligence Builder"]

    STORE["PostgreSQL"]

    PDF --> EXTRACT
    EXTRACT --> PARSE
    PARSE --> BUILD
    BUILD --> STORE
```

Detailed implementation is defined in **`05_resume_pipeline.md`**.

---

# 7. Interview Engine Placement

The Interview Engine is an internal module inside the Interview Service.

```text
Interview Service
    ├── Resume Pipeline
    ├── Interview Engine
    ├── WebSocket Manager
    └── Storage Layer
```

The engine is responsible only for interview orchestration.

---

# 8. External Integrations

| Integration | Purpose |
|-------------|---------|
| Gemini | LLM reasoning. |
| LangChain | Prompt execution and structured output parsing. |
| Web Speech API | Speech-to-text and text-to-speech on frontend. |

No external service communicates directly with PostgreSQL or Redis.

---

# 9. Architecture Principles

- Stateless HTTP APIs.
- Stateful interview sessions through Redis.
- PostgreSQL as persistent storage.
- LangGraph controls interview workflow.
- LangChain abstracts the LLM provider.
- Resume parsing happens once before interview initialization.

---

# 10. Out of Scope

The following are intentionally documented elsewhere.

| Topic | Owner Document |
|-------|----------------|
| Folder Structure | `03_folder_structure.md` |
| Interview Workflow | `04_interview_engine.md` |
| Resume Parsing Logic | `05_resume_pipeline.md` |
| Prompt Design | `06_prompt_architecture.md` |
| Database Tables | `07_database_schema.md` |
| Redis Keys & TTLs | `08_redis_strategy.md` |
| WebSocket Events | `09_websocket_protocol.md` |