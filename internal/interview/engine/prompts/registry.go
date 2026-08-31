package prompts

const (
	ResumeParserPrompt         = "resume_parser_v1"
	IntentDetectorPrompt       = "intent_detector_v1"
	GuardrailDetectorPrompt    = "guardrail_detector_v1"
	TechnicalQuestionPrompt    = "technical_question_v1"
	TechnicalEvaluatorPrompt   = "technical_evaluator_v1"
	FollowUpGeneratorPrompt    = "followup_generator_v1"
	ClarificationPrompt        = "clarification_generator_v1"
	HintPrompt                 = "hint_generator_v1"
	ThinkingPrompt             = "thinking_prompt_v1"
	BehavioralPrompt           = "behavioral_generator_v1"
	FinalEvaluationPrompt      = "final_evaluation_v1"
)

var registry = map[string]string{
	ResumeParserPrompt:       "resume_parser_v1.txt",
	IntentDetectorPrompt:     "intent_detector_v1.txt",
	GuardrailDetectorPrompt:  "guardrail_detector_v1.txt",
	TechnicalQuestionPrompt:  "technical_question_v1.txt",
	TechnicalEvaluatorPrompt: "technical_evaluator_v1.txt",
	FollowUpGeneratorPrompt:  "followup_generator_v1.txt",
	ClarificationPrompt:      "clarification_generator_v1.txt",
	HintPrompt:               "hint_generator_v1.txt",
	ThinkingPrompt:           "thinking_prompt_v1.txt",
	BehavioralPrompt:         "behavioral_generator_v1.txt",
	FinalEvaluationPrompt:    "final_evaluation_v1.txt",
}

func File(promptName string) (string, bool) {
	file, ok := registry[promptName]
	return file, ok
}