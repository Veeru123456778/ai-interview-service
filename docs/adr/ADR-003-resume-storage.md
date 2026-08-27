# ADR-003 — Resume Storage Strategy

**Status:** Accepted

## Context

The backend needs a reliable representation of the candidate's resume for interviews while preserving the original upload for future reprocessing and debugging.

## Decision

Store:

* **Original resume PDF** in private object storage (S3 / Cloudflare R2).
* **Resume metadata** in PostgreSQL.
* **Parsed structured resume JSON** in PostgreSQL.

The Interview Engine consumes only the structured JSON.

## Alternatives Considered

* Store only JSON and discard the PDF.
* Store the PDF inside PostgreSQL as binary data.

## Consequences

### Benefits

* Re-parse resumes if parsing improves.
* Keep database size small.
* Preserve original document for debugging and future features.

### Trade-offs

* Additional object storage dependency.
* Need to manage storage lifecycle.
