// package nodes

// import (
// 	"context"

// 	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
// )

// type InitializeSessionNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type ResolveResumeContextNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type SelectContextNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type SelectScenarioNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type GenerateQuestionNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type DetectCandidateIntentNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type GuardrailCheckNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type GenerateClarificationNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type GenerateHintNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type AnalyzeTechnicalAnswerNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type UpdateCandidateMemoryNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type DecideDifficultyNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type TransitionContextNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type BehavioralDiscussionNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }

// type GenerateFinalEvaluationNode interface {
// 	Execute(context.Context, *engine.InterviewState) error
// }


package nodes

import (
	"context"

	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
)

// ----------------------------------------------------------------------
// Shared Workflow Types
// ----------------------------------------------------------------------

type CandidateIntent string

const (
	IntentAnswer        CandidateIntent = "ANSWER"
	IntentHint          CandidateIntent = "REQUEST_HINT"
	IntentClarification CandidateIntent = "REQUEST_CLARIFICATION"
	IntentRepeat        CandidateIntent = "REQUEST_REPEAT"
)

type GuardrailResult struct {
	Triggered bool
	Reason    string
}

// ----------------------------------------------------------------------
// Session Nodes
// ----------------------------------------------------------------------

type InitializeSessionNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

type ResolveResumeContextNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

// ----------------------------------------------------------------------
// Context Nodes
// ----------------------------------------------------------------------

type SelectContextNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

type TransitionContextNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

// ----------------------------------------------------------------------
// Question Nodes
// ----------------------------------------------------------------------

type SelectScenarioNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

type GenerateQuestionNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

// ----------------------------------------------------------------------
// Intent & Guardrail Nodes
// ----------------------------------------------------------------------

type DetectCandidateIntentNode interface {
	Execute(context.Context, *engine.InterviewState) (CandidateIntent, error)
}

type GuardrailCheckNode interface {
	Execute(context.Context, *engine.InterviewState) (*GuardrailResult, error)
}

// ----------------------------------------------------------------------
// Assistance Nodes
// ----------------------------------------------------------------------

type GenerateClarificationNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

type GenerateHintNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

// ----------------------------------------------------------------------
// Evaluation Nodes
// ----------------------------------------------------------------------

type AnalyzeTechnicalAnswerNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

type UpdateCandidateMemoryNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

type DecideDifficultyNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

// ----------------------------------------------------------------------
// Completion Nodes
// ----------------------------------------------------------------------

type BehavioralDiscussionNode interface {
	Execute(context.Context, *engine.InterviewState) error
}

type GenerateFinalEvaluationNode interface {
	Execute(context.Context, *engine.InterviewState) error
}