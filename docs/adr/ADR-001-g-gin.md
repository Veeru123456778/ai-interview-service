# ADR-001 — Backend Architecture (Go + Gin)

**Status:** Accepted

## Context

The backend needs to support concurrent interview sessions, WebSocket communication, resume processing, Redis caching, PostgreSQL persistence, and AI orchestration while remaining simple enough for a single deployable service.

## Decision

Use **Go** as the backend language and **Gin** as the HTTP framework.

## Alternatives Considered

* **Node.js + Express** — Familiar ecosystem but weaker concurrency model for long-lived WebSocket connections.
* **Spring Boot** — Powerful but heavier for this project's scope.

## Consequences

### Benefits

* Lightweight concurrency with goroutines.
* Excellent performance for realtime communication.
* Small deployment footprint.
* Mature middleware ecosystem.

### Trade-offs

* Slightly steeper learning curve than Express.
* More explicit project structure than many Node frameworks.
