package interview

import (
	"net/http"

	apperrors "github.com/Veeru123456778/ai-interview-service/internal/shared/errors"
)

var (

	// ----------------------------------------------------------------------
	// Interview Lifecycle Errors
	// ----------------------------------------------------------------------

	ErrInterviewNotFound = &apperrors.AppError{
		Code:       "INTERVIEW_NOT_FOUND",
		Message:    "Interview not found.",
		StatusCode: http.StatusNotFound,
	}

	ErrInterviewAlreadyCompleted = &apperrors.AppError{
		Code:       "INTERVIEW_ALREADY_COMPLETED",
		Message:    "Interview has already been completed.",
		StatusCode: http.StatusConflict,
	}

	ErrInterviewAlreadyActive = &apperrors.AppError{
		Code:       "INTERVIEW_ALREADY_ACTIVE",
		Message:    "An active interview already exists for this session.",
		StatusCode: http.StatusConflict,
	}

	// ----------------------------------------------------------------------
	// Session Errors
	// ----------------------------------------------------------------------

	ErrInterviewSessionExpired = &apperrors.AppError{
		Code:       "INTERVIEW_SESSION_EXPIRED",
		Message:    "Interview session has expired.",
		StatusCode: http.StatusUnauthorized,
	}

	ErrInterviewSessionNotFound = &apperrors.AppError{
		Code:       "INTERVIEW_SESSION_NOT_FOUND",
		Message:    "Interview session not found.",
		StatusCode: http.StatusNotFound,
	}

	// ----------------------------------------------------------------------
	// Interview Runtime Errors
	// ----------------------------------------------------------------------

	ErrInterviewEngineFailed = &apperrors.AppError{
		Code:       "INTERVIEW_ENGINE_FAILED",
		Message:    "Interview engine execution failed.",
		StatusCode: http.StatusInternalServerError,
	}

	ErrQuestionGenerationFailed = &apperrors.AppError{
		Code:       "QUESTION_GENERATION_FAILED",
		Message:    "Unable to generate interview question.",
		StatusCode: http.StatusInternalServerError,
	}

	ErrAnswerEvaluationFailed = &apperrors.AppError{
		Code:       "ANSWER_EVALUATION_FAILED",
		Message:    "Unable to evaluate candidate answer.",
		StatusCode: http.StatusInternalServerError,
	}

	ErrFinalEvaluationFailed = &apperrors.AppError{
		Code:       "FINAL_EVALUATION_FAILED",
		Message:    "Unable to generate final interview evaluation.",
		StatusCode: http.StatusInternalServerError,
	}
)