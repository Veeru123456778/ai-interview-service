# AI Interview Service — Technical Overview

**Project:** AI Interview Service

**Document:** `docs/tech-spec/00_overview.md`

**Version:** 1.0 (LOCKED)

---

# 1. Purpose

This document is the entry point for the technical documentation of the AI Interview Service.

It provides the documentation structure, reading order, and ownership of each architecture document. Detailed implementation and architecture are defined in their respective tech-spec files.

---

# 2. Documentation Reading Order

Read the documents in the following order before implementation.

| Order | Document | Purpose |
|-------|----------|---------|
| 0 | `product_scope.md` | Product vision, goals, user flow, and project scope. |
| 1 | `01_tech_stack.md` | Technology choices and architectural decisions. |
| 2 | `02_architecture.md` | High-Level Architecture (HLD) and service communication. |
| 3 | `03_folder_structure.md` | Repository structure and package ownership. |
| 4 | `04_interview_engine.md` | LangGraph workflow and interview orchestration. |
| 5 | `05_resume_pipeline.md` | Resume Intelligence Pipeline. |
| 6 | `06_prompt_architecture.md` | Prompt contracts and LLM interaction design. |
| 7 | `07_database_schema.md` | PostgreSQL schema and persistence model. |
| 8 | `08_redis_strategy.md` | Runtime state model in Redis. |
| 9 | `09_api_spec.md` | REST API contracts. |
| 10 | `10_websocket_protocol.md` | Real-time interview communication protocol. |
| 11 | `11_auth_security.md` | Authentication and security architecture. |
| 12 | `12_evaluation_engine.md` | Technical evaluation model and scoring. |
| 13 | `13_observability.md` | Logging, metrics, tracing, and monitoring. |
| 14 | `14_deployment.md` | Docker, deployment, CI/CD, and production infrastructure. |

---

# 3. Documentation Ownership

Each concept has a single owner document. Other documents reference it instead of duplicating it.

| Concept | Owner Document |
|---------|----------------|
| Product Vision & Scope | `product_scope.md` |
| Technology Decisions | `01_tech_stack.md` |
| High-Level Architecture | `02_architecture.md` |
| Repository Structure | `03_folder_structure.md` |
| Interview Workflow | `04_interview_engine.md` |
| Resume Processing | `05_resume_pipeline.md` |
| Prompt Contracts | `06_prompt_architecture.md` |
| Database Design | `07_database_schema.md` |
| Redis Runtime State | `08_redis_strategy.md` |
| API Specification | `09_api_spec.md` |
| WebSocket Protocol | `10_websocket_protocol.md` |
| Authentication & Security | `11_auth_security.md` |
| Evaluation & Scoring | `12_evaluation_engine.md` |
| Observability | `13_observability.md` |
| Deployment Architecture | `14_deployment.md` |

---

# 4. Architecture Principles

The backend follows a domain-first, production-oriented architecture.

### Core Principles

- Go + Gin backend.
- Domain-first package organization.
- LangGraph orchestrates interview workflow.
- LangChain provides LLM abstraction.
- PostgreSQL is the source of truth.
- Redis stores runtime interview state only.
- WebSocket powers real-time interview communication.
- Prompt templates are versioned and managed independently.

Detailed architecture is defined in **`02_architecture.md`**.

---

# 5. Implementation Order

The project should be implemented in the following order.

| Phase | Documents |
|-------|-----------|
| Foundation | Tech Stack, Architecture, Folder Structure |
| Core Backend | Resume Pipeline, Interview Engine, Prompt Architecture |
| Persistence | Database Schema, Redis Strategy |
| Communication | WebSocket Protocol, API Specification |
| Security | Auth & Security |
| Intelligence | Evaluation Engine |
| Production | Observability, Deployment |

---

# 6. ADR Index

Architectural Decision Records (ADRs) document permanent design decisions.

| ADR | Decision |
|-----|----------|
| ADR-001 | Go + Gin as backend framework. |
| ADR-002 | LangGraph + LangChain architecture. |
| ADR-003 | Resume Intelligence Pipeline design. |
| ADR-004 | PostgreSQL + Redis persistence strategy. |
| ADR-005 | Dependency direction and package boundaries. |

All ADRs are located in `docs/adr/`.

---

# 7. Versioning Rules

- Every tech-spec is independently versioned.
- Version `1.0` indicates the document is **locked** for V1.
- Architectural changes require a new version (`v1.1`) and a corresponding ADR if the change affects system design.