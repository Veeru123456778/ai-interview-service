package interview

// ----------------------------------------------------------------------
// Conversation Roles
// ----------------------------------------------------------------------

const (
	RoleInterviewer = "INTERVIEWER"
	RoleCandidate   = "CANDIDATE"
)

// ----------------------------------------------------------------------
// Redis Session Keys
// ----------------------------------------------------------------------

const (
	InterviewSessionKeyPrefix = "interview:session:"
	InterviewLockKeyPrefix    = "interview:lock:"
)

// ----------------------------------------------------------------------
// WebSocket Event Types
// ----------------------------------------------------------------------

const (
	EventInterviewStarted   = "INTERVIEW_STARTED"
	EventQuestionGenerated  = "QUESTION_GENERATED"
	EventCandidateAnswered  = "CANDIDATE_ANSWERED"
	EventHintGenerated      = "HINT_GENERATED"
	EventClarificationSent  = "CLARIFICATION_GENERATED"
	EventFollowUpGenerated  = "FOLLOWUP_GENERATED"
	EventContextTransition  = "CONTEXT_TRANSITION"
	EventBehavioralStarted  = "BEHAVIORAL_STARTED"
	EventInterviewCompleted = "INTERVIEW_COMPLETED"
	EventGuardrailTriggered = "GUARDRAIL_TRIGGERED"
)

// ----------------------------------------------------------------------
// Candidate Intent Types
// ----------------------------------------------------------------------

const (
	IntentAnswer        = "ANSWER"
	IntentHint          = "ASK_HINT"
	IntentClarification = "REQUEST_CLARIFICATION"
	IntentThinking      = "THINKING_OUT_LOUD"
	IntentRepeat        = "REQUEST_REPEAT"
	IntentOffTopic      = "OFF_TOPIC"
	IntentEndInterview  = "END_REQUEST"
)

// ----------------------------------------------------------------------
// Guardrail Categories
// ----------------------------------------------------------------------

const (
	GuardrailNormal          = "NORMAL"
	GuardrailOffTopic        = "OFF_TOPIC"
	GuardrailPromptInjection = "PROMPT_INJECTION"
	GuardrailUnsupported     = "UNSUPPORTED"
)