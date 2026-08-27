# AI Interview Service — Product Scope (V1)

**Version:** 1.0 (Locked V1)

**Status:** Product Requirements Document (PRD)

---

# 1. Vision

AI Interview Service is a realistic AI-powered interview coach designed specifically for **Software Engineer (SDE)** candidates.

The goal is to simulate the experience of interviewing with an experienced software engineer rather than chatting with a generic AI assistant.

The interviewer should conduct dynamic, natural conversations, deeply probe the candidate's projects and technical decisions, adapt the interview based on performance, and provide actionable feedback after the interview.

**Core Principle**

> If a real software engineer from a top tech company would ask it during an interview, our AI interviewer should be capable of asking it naturally.

---

# 2. Problem Statement

Most AI interview tools today suffer from one or more of these issues:

* Ask a predefined list of questions.
* Move to the next question regardless of the candidate's answer.
* Ask generic technical questions unrelated to the candidate's resume.
* Provide shallow or unreliable feedback.
* Feel like chatting with an LLM instead of interviewing with a human.

Candidates preparing for software engineering interviews need realistic practice that adapts to their experience and projects.

---

# 3. Product Goal

Build the most realistic AI mock interviewer for Software Engineers.

The interviewer should:

* Read and understand a candidate's resume.
* Ask resume-specific technical questions.
* Continuously adapt the conversation.
* Probe deeper when answers are incomplete.
* Increase or decrease difficulty naturally.
* Evaluate answers throughout the interview.
* Produce a detailed interview report at the end.

---

# 4. Target Audience

### Primary Users

* College students preparing for internships.
* Final-year students preparing for SDE placements.
* Software engineers with 0–3 years of experience preparing for interviews.

### Supported Roles (V1)

* Backend Engineer
* Frontend Engineer
* Full Stack Engineer
* Software Development Engineer (General)

---

# 5. V1 Scope (Locked)

## Input

Candidate uploads a resume (PDF).

## Output

A complete AI mock interview lasting approximately **30–45 minutes**, followed by a detailed interview report.

## Supported Interview Areas

### Projects

The interviewer deeply explores projects mentioned in the resume.

Examples:

* Architecture
* Design decisions
* Tradeoffs
* Scaling
* Failure scenarios
* Technologies used
* Candidate ownership

### Technical Fundamentals

Questions related to technologies mentioned in the resume.

Examples:

* Go
* Node.js
* React
* Redis
* Kafka
* SQL
* MongoDB
* Docker
* REST APIs

### Computer Science Fundamentals

Questions generated based on resume relevance.

Examples:

* DBMS
* Operating Systems
* Networking
* OOP
* Concurrency

### DSA Discussion

Discussion around algorithmic thinking and complexity (not live coding in V1).

---

# 6. Out of Scope (V1)

The following features are intentionally excluded.

## Resume + Job Description Personalization

Resume only.

## HR / Behavioral Interview

No behavioral rounds in V1.

## Live Coding Editor

No coding IDE.

## Video Interview

Audio and chat only.

## Multi-company Personalities

Only one interviewer personality.

## Payments / Subscriptions

No monetization in V1.

## Enterprise Dashboard

Single-user product only.

---

# 7. Core User Journey

1. User signs in.
2. User uploads resume.
3. Resume is analyzed.
4. Interview begins.
5. AI interviewer conducts conversation.
6. AI adapts questions throughout the interview.
7. Interview ends.
8. User receives interview report.

The experience should feel continuous and conversational.

---

# 8. Interview Experience Principles

## Dynamic Conversation

The interviewer never follows a fixed list of questions.

Each next question depends on:

* previous answers,
* unanswered technical areas,
* confidence level,
* interview progress.

## Resume-Aware

Every project and technology on the resume is considered potential interview material.

The interviewer should ask questions that verify genuine understanding rather than definitions.

## Deep Technical Probing

The interviewer continues exploring a topic until sufficient evidence is collected.

Example progression:

* Explain your project.
* Why did you choose this technology?
* What alternatives did you consider?
* What happens when this component fails?
* How would you scale this design?

## Natural Human Conversation

The interviewer behaves like a human interviewer.

Possible behaviors include:

* asking follow-up questions,
* interrupting long answers,
* asking clarification questions,
* revisiting earlier topics,
* increasing difficulty,
* changing topics naturally.

---

# 9. Interview Categories

| Category                    | Approximate Focus |
| --------------------------- | ----------------- |
| Work Experience             | Highest priority
| Resume Projects             | Higher priority  |
| Technical Stack from Resume | High priority     |
| CS Fundamentals             | Medium priority   |
| DSA Discussion              | Medium priority   |

The interviewer dynamically chooses categories during the interview.

---

# 10. What Makes This Interview Realistic

The interviewer should be capable of:

* remembering earlier answers,
* referencing earlier projects,
* identifying contradictions,
* asking increasingly difficult follow-ups,
* changing topic only after sufficient discussion,
* asking architecture and failure questions,
* validating ownership of resume claims.

The interview should never feel scripted.

---

# 11. Feedback Goals

At the end of every interview, users receive structured feedback.

### Overall Evaluation

* Overall interview score.
* Technical communication.
* Problem-solving quality.
* Depth of understanding.
* Confidence.

### Strengths

Areas where the candidate performed well.

### Weaknesses

Topics that require improvement.

### Missed Opportunities

Important concepts or tradeoffs the candidate failed to mention.

### Suggested Practice

Personalized technical topics for future practice.

---

# 12. Product Quality Bar

The product is considered successful only if interviews satisfy these principles:

* Questions are technically accurate.
* Questions are relevant to the candidate's resume.
* Follow-up questions feel natural.
* The interviewer remembers previous context.
* Difficulty adapts throughout the interview.
* Feedback is actionable and evidence-based.

The AI should behave like an experienced software engineer conducting a real interview rather than a chatbot generating random questions.

---

# 13. Success Metrics (V1)

The initial product will be evaluated using:

## User Experience

* Interview completion rate.
* Average interview duration.
* Return rate for a second interview.

## Interview Quality

* Users feel questions were relevant to their resume.
* Users feel follow-up questions were realistic.
* Users find the feedback useful for interview preparation.

## Technical Quality

* Stable real-time interview sessions.
* Low response latency.
* Reliable resume understanding.
* Consistent interview flow.

---

# 14. V1 Product Philosophy

Build one exceptional interview experience for Software Engineers before expanding the scope.

Every new feature added after V1 should improve interview realism rather than increase feature count.

**Focus over breadth.**
