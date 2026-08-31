package engine

import (
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/resume"
)

// Represents the complete runtime state of an active interview.
type InterviewState struct {
	InterviewID string `json:"interview_id"`
	UserID      string `json:"user_id"`
	ResumeID    string `json:"resume_id"`

	Status string `json:"status"`

	// Resume contexts available for this interview.
	ResumeContexts []resume.InterviewContext `json:"resume_contexts"`

	// Active interview context.
	CurrentContextIndex int `json:"current_context_index"`
	CurrentContextID    string `json:"current_context_id"`
	CurrentContextType  string `json:"current_context_type"`
	CurrentContextTitle string `json:"current_context_title"`

	// Active topic and question state.
	CurrentTopicID    string `json:"current_topic_id"`
	CurrentDifficulty string `json:"current_difficulty"`
	CurrentQuestionID string `json:"current_question_id"`
	CurrentQuestionNo int    `json:"current_question_no"`
	CurrentScenario string `json:"current_scenario"`
	FollowUpCount     int    `json:"follow_up_count"`

	AskedQuestions  []string `json:"asked_questions"`
	CompletedTopics []string `json:"completed_topics"`

	ConversationHistory []ConversationMessage `json:"conversation_history"`
	TopicScores         []TopicScore          `json:"topic_scores"`

	FinalEvaluation *FinalEvaluation `json:"final_evaluation,omitempty"`

	StartedAt     time.Time `json:"started_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// Stores one message exchanged during the interview.
type ConversationMessage struct {
	Role       string    `json:"role"` // INTERVIEWER or CANDIDATE
	Content    string    `json:"content"`
	QuestionID string    `json:"question_id,omitempty"`
	TopicID    string    `json:"topic_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Running score for one technology/topic.
type TopicScore struct {
	TopicID    string `json:"topic_id"`
	TopicName  string `json:"topic_name"`
	Difficulty string `json:"difficulty"`

	QuestionsAsked  int `json:"questions_asked"`
	QuestionsPassed int `json:"questions_passed"`

	Score float64 `json:"score"`
}

// Final interview summary generated after all contexts are completed.
type FinalEvaluation struct {
	AverageScore   float64 `json:"average_score"`
	Recommendation string  `json:"recommendation"`
	TotalQuestions int     `json:"total_questions"`
}