# ADR-002 — AI Orchestration with LangGraph and LangChain

**Status:** Accepted

## Context

The interview requires dynamic follow-up questions, conversation memory, topic transitions, structured outputs, and the flexibility to support multiple LLM providers in the future.

## Decision

* Use **LangGraph** for interview orchestration.
* Use **LangChain** only as the LLM provider abstraction layer.
* Use **Gemini** as the initial provider.

## Alternatives Considered

* **Direct Gemini SDK only** — Simpler but tightly coupled to one provider.
* **LangChain Agents** — Too flexible for a deterministic interview workflow.
* **LangChain Memory** — Backend owns interview state instead.

## Consequences

### Benefits

* Provider independence.
* Structured JSON outputs.
* Streaming support.
* Deterministic interview workflow through LangGraph.

### Trade-offs

* Adds one abstraction layer.
* Requires maintaining provider implementations.
