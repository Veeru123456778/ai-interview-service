# Evaluation Engine

**Project:** AI Interview Service

**Document:** `docs/tech-spec/12_evaluation_engine.md`

**Version:** 1.0 (LOCKED)

---

# Ownership

**This document owns:**

- Topic-level evaluation model.
- Final interview evaluation model.
- Scoring dimensions.
- Evaluation aggregation rules.
- Learning recommendations generation.

**References**

- `04_interview_engine.md` → Runtime interview workflow.
- `05_resume_pipeline.md` → Interview contexts and Technology Graph.
- `07_database_schema.md` → Evaluation persistence.
- `08_redis_strategy.md` → Candidate memory updates.
- `09_api_spec.md` → Evaluation APIs.
- `06_prompt_architecture.md` → Technical Evaluator prompt.

---

# 1. Purpose

The Evaluation Engine measures a candidate's performance throughout the interview.

Evaluation happens at **two levels**:

1. **Topic Evaluation** — after every completed topic.
2. **Final Evaluation** — after the interview finishes.

The Interview Engine owns interview progression.

The Evaluation Engine owns scoring only.

---

# 2. Evaluation Flow

```text
QUESTION

   │

ANSWER

   │

Technical Evaluator Prompt

   │

Topic Evaluation Update

   │

Candidate Memory Update

   │

Next Topic / Interview Completion

   │

Final Evaluation Prompt

   ▼

Persist Results to PostgreSQL
```

---

# 3. Evaluation Levels

| Level | Purpose |
|--------|---------|
| Answer Evaluation | Evaluate one candidate answer. |
| Topic Evaluation | Aggregate answers for one technology topic. |
| Final Evaluation | Aggregate all topic evaluations into one interview report. |

Answer evaluations are temporary.

Topic and Final evaluations are persisted in PostgreSQL.

---

# 4. Answer Evaluation

Every candidate answer is evaluated independently.

## Evaluation Dimensions

| Dimension | Description | Score Range |
|-----------|-------------|-------------|
| Technical Correctness | Accuracy of concepts and implementation. | 1–5 |
| Depth of Understanding | Level of understanding beyond basic knowledge. | 1–5 |
| Reasoning | Problem-solving and decision making. | 1–5 |
| Communication | Clarity and structure of explanation. | 1–5 |

The evaluator assigns each dimension a score between **1 and 5**.

The weighted topic score is calculated by applying the dimension weights and then normalizing the result to a **0–100** topic score.


## Evaluator Output

```json
{
  "technical_score": 4,
  "depth_score": 3,
  "reasoning_score": 5,
  "communication_score": 4,

  "strengths": [
    "Explained Redis session recovery."
  ],

  "weaknesses": [
    "Did not discuss persistence trade-offs."
  ],

  "follow_up_required": true
}
```

This output is consumed by the Interview Engine.

---

# 5. Topic Evaluation

A Topic Evaluation represents the candidate's overall performance for one technology.

Example: Redis, WebSocket, PostgreSQL.

## Aggregation Inputs

- All answer evaluations for the topic.
- Follow-up performance.
- Candidate reasoning improvements during the topic.

## Topic Evaluation Structure

```json
{
  "topic_id": "redis-topic-id",

  "topic_name": "Redis",

  "context_name": "AI Interview Platform",

  "score": 86,

  "difficulty_reached": "HARD",

  "strengths": [
    "Good understanding of session recovery."
  ],

  "weaknesses": [
    "Cache invalidation strategy."
  ]
}
```

## Stored In

`topic_evaluations` table.

---

# 6. Topic Aggregation Rules

## Scoring Rule

The topic score is calculated from all answer evaluations within that topic.

| Dimension | Weight |
|-----------|--------|
| Technical Correctness | 40% |
| Depth of Understanding | 25% |
| Reasoning | 25% |
| Communication | 10% |

The final topic score is normalized to **0–100**.

## Difficulty Tracking

The Interview Engine tracks runtime difficulty progression.

The Evaluation Engine stores the **highest difficulty successfully reached** for the topic.

Possible values:

- EASY
- MEDIUM
- HARD

---

# 7. Candidate Memory Updates

After every evaluated answer, candidate memory is refreshed.

## Memory Update Responsibilities

| Memory Field | Updated From |
|--------------|-------------|
| Conversation Summary | Latest evaluated discussion. |
| Strengths | Correct concepts observed. |
| Weaknesses | Missing or incorrect concepts. |

Memory is stored only in Redis during the interview.

Memory is **not** persisted as evaluation history.

---

# 8. Final Interview Evaluation

The Final Evaluation aggregates every topic evaluation.

## Inputs

- Topic evaluations.
- Resume metadata.
- Candidate strengths.
- Candidate weaknesses.

## Output Structure

```json
{
  "overall_score": 84,

  "recommendation": "STRONG_HIRE",

  "strengths": [
    "Strong backend reasoning.",
    "Good debugging approach."
  ],

  "weaknesses": [
    "Distributed caching concepts need improvement."
  ],

  "learning_recommendations": [
    "Redis persistence.",
    "Cache invalidation.",
    "Horizontal scaling."
  ]
}
```

Stored in `final_evaluations`.

---

# 9. Overall Score Aggregation

The overall interview score is calculated from completed topic evaluations.

## Aggregation Rule

- Every completed topic contributes equally.
- Only completed topics participate.
- Overall score is normalized to **0–100**.

Example:

| Topic | Score |
|-------|------:|
| Redis | 88 |
| WebSocket | 80 |
| PostgreSQL | 84 |

Overall score = Average topic score.

---

# 10. Recommendation Rules

The recommendation is generated from the overall interview score.

| Overall Score | Recommendation |
|--------------|----------------|
| 85–100 | HIRE |
| 70–84 | STRONG_HIRE |
| 55–69 | HOLD |
| 0–54 | NO_HIRE |

Recommendation generation is performed by the Final Evaluation prompt.

---

# 11. Learning Recommendations

Learning recommendations are generated from repeated weaknesses across topics.

## Example

| Weakness Detected | Recommendation |
|-------------------|---------------|
| Redis persistence | Learn Redis persistence strategies. |
| Cache invalidation | Learn cache invalidation patterns. |
| WebSocket reconnects | Practice WebSocket lifecycle and recovery. |

Recommendations are concise improvement areas, not study plans.

---

# 12. Evaluation Persistence

Evaluation is persisted only after interview completion.

## PostgreSQL Tables

| Table | Stores |
|-------|--------|
| `topic_evaluations` | Topic-wise evaluation. |
| `final_evaluations` | Overall interview evaluation. |
| `conversation_messages` | Conversation transcript. |

Redis evaluation state is deleted after persistence.

---

# 13. Evaluation Design Rules

| Rule | Description |
|------|-------------|
| Evaluation never changes interview progression. | Interview Engine controls progression. |
| Topic evaluation is updated after each completed topic. | Runtime aggregation only. |
| Final evaluation is generated once. | After interview completion. |
| Redis stores temporary memory only. | PostgreSQL stores permanent evaluations. |
| Scores are normalized to 0–100. | Consistent reporting across interviews. |

---

# 14. Related Documents

| Topic | Document |
|-------|----------|
| Interview Engine | `04_interview_engine.md` |
| Resume Pipeline | `05_resume_pipeline.md` |
| Prompt Architecture | `06_prompt_architecture.md` |
| Database Schema | `07_database_schema.md` |
| Redis Strategy | `08_redis_strategy.md` |
| API Specification | `09_api_spec.md` |