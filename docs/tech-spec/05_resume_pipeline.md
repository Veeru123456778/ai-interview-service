# Resume Intelligence Pipeline

**Project:** AI Interview Service

**Document:** `docs/tech-spec/05_resume_pipeline.md`

**Version:** 1.1 (LOCKED)

---

# Ownership

**This document owns:**

- Resume processing pipeline.
- Resume validation.
- Resume Intelligence generation.
- Technology Graph generation.
- Interview Context generation.

**References**

- `02_architecture.md` → System architecture.
- `04_interview_engine.md` → Consumer of Resume Intelligence.
- `06_prompt_architecture.md` → Resume Parser prompt contract.
- `07_database_schema.md` → Resume persistence model.

---

# 1. Purpose

The Resume Pipeline converts an uploaded resume into **Resume Intelligence**, a structured representation of the candidate's projects, work experience, technologies, and interview contexts.

This pipeline executes **once per uploaded resume** before an interview starts.

The Interview Engine never reads the raw resume. It consumes only the generated Resume Intelligence.

---

# 2. Pipeline Overview

```mermaid
%%{init:{
  "theme":"base",
  "themeVariables":{
    "background":"#FFFFFF",
    "primaryColor":"#1E3A8A",
    "primaryTextColor":"#FFFFFF",
    "primaryBorderColor":"#111827",
    "secondaryColor":"#2563EB",
    "secondaryTextColor":"#FFFFFF",
    "tertiaryColor":"#DBEAFE",
    "lineColor":"#111827",
    "edgeLabelBackground":"#FFFFFF",
    "fontFamily":"Inter"
  }
}}%%

flowchart LR

    PDF["Resume PDF"]

    EXTRACT["Extract Text"]

    NORMALIZE["Normalize Resume"]

    PARSER["Resume Parser (LLM)"]

    VALIDATE["Validate Resume JSON"]

    GRAPH["Technology Graph Builder"]

    CONTEXT["Interview Context Builder"]

    PDF --> EXTRACT
    EXTRACT --> NORMALIZE
    NORMALIZE --> PARSER
    PARSER --> VALIDATE
    VALIDATE --> GRAPH
    GRAPH --> CONTEXT
```

---

# 3. Pipeline Stages

| Stage | Responsibility |
|-------|----------------|
| Extract Text | Extract text from uploaded PDF. |
| Normalize Resume | Remove formatting noise and normalize sections. |
| Resume Parser | Convert resume text into structured JSON. |
| Validate Resume JSON | Validate required fields and schema. |
| Technology Graph Builder | Extract technologies and relationships. |
| Interview Context Builder | Generate interview-ready contexts and topics. |

---

# 4. Resume Validation

Validation happens immediately after Resume Parser output.

## Required Sections

| Section | Required |
|---------|----------|
| Name | ✅ |
| Skills | ✅ |
| Projects | ✅ |
| Experience | Optional |
| Education | Optional |

## Validation Rules

- Projects may be empty.
- Experience may be empty.
- Skills must exist.
- Invalid parser JSON is retried once.
- Resume processing stops if validation fails twice.

---

# 5. Resume Intelligence

Resume Intelligence is the only input used by the Interview Engine.

## Generated Fields

| Field | Purpose |
|-------|---------|
| Skills | Normalized candidate skills. |
| Projects | Parsed projects with descriptions and technologies. |
| Experience | Parsed work experience. |
| Technology Graph | Technology relationships extracted from resume. |
| Interview Contexts | Context-aware interview plan. |

Resume Intelligence is stored inside the `resumes` table.

---

# 6. Technology Graph Builder

The Technology Graph identifies technologies and where they were used.

```mermaid
%%{init:{
  "theme":"base",
  "themeVariables":{
    "background":"#FFFFFF",
    "primaryColor":"#1E3A8A",
    "primaryTextColor":"#FFFFFF",
    "primaryBorderColor":"#111827",
    "secondaryColor":"#2563EB",
    "secondaryTextColor":"#FFFFFF",
    "tertiaryColor":"#DBEAFE",
    "lineColor":"#111827",
    "edgeLabelBackground":"#FFFFFF",
    "fontFamily":"Inter"
  }
}}%%

flowchart TD

    PROJECTS["Projects"]

    EXPERIENCE["Work Experience"]

    SKILLS["Skills"]

    GRAPH["Technology Graph"]

    PROJECTS --> GRAPH
    EXPERIENCE --> GRAPH
    SKILLS --> GRAPH
```

## Technology Node

Each technology node contains:

| Field | Purpose |
|-------|---------|
| `topic_id` | Stable UUID. |
| `name` | Technology name. |
| `category` | Backend, Database, AI, Cloud, etc. |
| `confidence` | Extraction confidence. |
| `sources` | Projects or experience where it appeared. |

Example:

```json
{
  "topic_id": "uuid",
  "name": "Redis",
  "category": "Database",
  "sources": [
    "AI Interview Platform",
    "Backend Internship"
  ]
}
```

---

# 7. Interview Context Builder

This stage converts the Technology Graph into **Interview Contexts**.

> The interview is driven by **where a technology was used**, not by isolated technologies.

## Context Types

| Context Type | Description |
|--------------|-------------|
| `PROJECT` | Candidate-built projects. |
| `EXPERIENCE` | Professional work experience or internship. |
| `SKILL` | Standalone skills not tied to a project or experience. |

## Context Priority

1. Projects.
2. Work Experience.
3. Standalone Skills.

Projects receive the highest interview priority because they provide the richest implementation context.

---

# 8. Interview Context Structure

Each interview context contains technologies that belong to that project or experience.

```json
{
  "context_id": "uuid",
  "context_type": "PROJECT",
  "context_name": "AI Interview Platform",
  "description": "Real-time AI mock interview platform.",
  "priority": 1,

  "topics": [
    {
      "topic_id": "uuid",
      "name": "Redis",
      "difficulty": "MEDIUM"
    },
    {
      "topic_id": "uuid",
      "name": "WebSocket",
      "difficulty": "MEDIUM"
    },
    {
      "topic_id": "uuid",
      "name": "PostgreSQL",
      "difficulty": "HARD"
    }
  ]
}
```

This structure is stored in the `interview_topics` JSONB column.

---

# 9. Interview Planning Rules

The Interview Context Builder generates the interview plan before the interview starts.

## Planning Rules

- Group technologies by project or work experience.
- Preserve the relationship between technologies and their implementation context.
- Assign interview priority to each context.
- Assign initial difficulty to each technology.
- Generate stable `context_id` and `topic_id`.

The Interview Engine consumes this plan without modifying it.

---

# 10. Practical Question Philosophy

Resume Intelligence prepares the Interview Engine for **scenario-based interviews**, not resume-reading questions.

## Question Style

| Context | Expected Question Style |
|--------|--------------------------|
| Project | Implementation, debugging, scaling, failures, trade-offs. |
| Work Experience | Ownership, production incidents, architecture decisions. |
| Standalone Skill | Practical concepts and applied scenarios. |

### Example

Instead of:

> Why did you use Redis?

The Interview Engine receives enough context to ask:

> What happens if Redis experiences a large number of cache misses while thousands of interview sessions are active?

This behavior is implemented in the Interview Engine and Prompt Architecture.

---

# 11. Resume Parser Output

The Resume Parser returns structured JSON containing:

- Skills.
- Projects.
- Experience.
- Education.
- Certifications (optional).

The parser **does not** generate interview questions or scores.

Prompt details are defined in `06_prompt_architecture.md`.

---

# 12. Error Handling

| Failure | Action |
|---------|--------|
| PDF extraction failed | Reject upload. |
| Invalid parser JSON | Retry parser once. |
| Validation failed | Return validation error. |
| Technology graph generation failed | Reject resume processing. |
| Interview context generation failed | Reject resume processing. |

Interview creation is allowed only after successful Resume Intelligence generation.

---

# 13. Output Consumers

| Consumer | Uses |
|----------|------|
| Interview Engine | Interview contexts, technologies, resume metadata. |
| PostgreSQL | Persist Resume Intelligence. |
| Evaluation Engine | Resume metadata for final report. |

Resume Intelligence becomes read-only after successful processing.

---

# 14. Related Documents

| Topic | Document |
|-------|----------|
| Interview Engine | `04_interview_engine.md` |
| Prompt Architecture | `06_prompt_architecture.md` |
| Database Schema | `07_database_schema.md` |