package nodes

import (
	"context"
	"fmt"
	"strings"

	"github.com/Veeru123456778/ai-interview-service/internal/interview/engine"
)

// ----------------------------------------------------------------------
// Detect Candidate Intent Node
// ----------------------------------------------------------------------

type detectCandidateIntentNode struct{}

func NewDetectCandidateIntentNode() DetectCandidateIntentNode {
	return &detectCandidateIntentNode{}
}

func (n *detectCandidateIntentNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) (CandidateIntent, error) {

	_ = ctx

	if len(state.ConversationHistory) == 0 {
		return IntentAnswer, fmt.Errorf("conversation history is empty")
	}

	lastMessage := state.ConversationHistory[len(state.ConversationHistory)-1]

	if lastMessage.Role != "CANDIDATE" {
		return IntentAnswer, fmt.Errorf("last message is not from candidate")
	}

	content := strings.ToLower(strings.TrimSpace(lastMessage.Content))

	switch {
	case strings.Contains(content, "hint"):
		return IntentHint, nil

	case strings.Contains(content, "repeat") ||
		strings.Contains(content, "again"):
		return IntentRepeat, nil

	case strings.Contains(content, "clarify") ||
		strings.Contains(content, "don't understand") ||
		strings.Contains(content, "didn't understand"):
		return IntentClarification, nil

	default:
		return IntentAnswer, nil
	}
}

// ----------------------------------------------------------------------
// Guardrail Check Node
// ----------------------------------------------------------------------

type guardrailCheckNode struct{}

func NewGuardrailCheckNode() GuardrailCheckNode {
	return &guardrailCheckNode{}
}

func (n *guardrailCheckNode) Execute(
	ctx context.Context,
	state *engine.InterviewState,
) (*GuardrailResult, error) {

	_ = ctx

	result := &GuardrailResult{
		Triggered: false,
		Reason:    "",
	}

	if len(state.ConversationHistory) == 0 {
		return result, nil
	}

	lastMessage := state.ConversationHistory[len(state.ConversationHistory)-1]

	if lastMessage.Role != "CANDIDATE" {
		return result, nil
	}

	content := strings.ToLower(lastMessage.Content)

	guardrailPatterns := []string{
		"give me the answer",
		"what is the answer",
		"solve this for me",
		"chatgpt",
		"google",
		"copy paste",
	}

	for _, pattern := range guardrailPatterns {
		if strings.Contains(content, pattern) {
			result.Triggered = true
			result.Reason = "ANSWER_SEEKING"
			return result, nil
		}
	}

	return result, nil
}