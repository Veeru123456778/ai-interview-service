# ADR-004 — PostgreSQL and Redis Responsibilities

**Status:** Accepted

## Context

The application needs both durable persistence and fast temporary state for realtime interviews.

## Decision

* **PostgreSQL** is the single source of truth.
* **Redis** stores only temporary interview state and cache.

## Responsibility Split

| PostgreSQL           | Redis                          |
| -------------------- | ------------------------------ |
| Users                | Active interview session state |
| Resume metadata      | WebSocket session mapping      |
| Parsed resume JSON   | Current interview topic        |
| Conversation history | Rate limiting counters         |
| Evaluation reports   | Temporary cache                |

## Alternatives Considered

* Redis as the conversation store.
* PostgreSQL-only architecture without Redis.

## Consequences

### Benefits

* Durable conversation history.
* Fast reconnects.
* Small Redis memory footprint.

### Trade-offs

* State synchronization between PostgreSQL and Redis during active interviews.
