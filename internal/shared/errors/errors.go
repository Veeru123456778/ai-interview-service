package errors

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Authentication required.",
		StatusCode: 401,
	}

	ErrForbidden = &AppError{
		Code:       "FORBIDDEN",
		Message:    "You do not have access to this resource.",
		StatusCode: 403,
	}

	ErrResumeRequired = &AppError{
		Code:       "RESUME_REQUIRED",
		Message:    "Resume file is required.",
		StatusCode: 400,
	}

	ErrInvalidResumeFormat = &AppError{
		Code:       "INVALID_RESUME_FORMAT",
		Message:    "Only PDF resumes are supported.",
		StatusCode: 400,
	}

	ErrResumeValidationFailed = &AppError{
		Code:       "RESUME_VALIDATION_FAILED",
		Message:    "Resume validation failed.",
		StatusCode: 422,
	}

	ErrInternalServer = &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "Something went wrong.",
		StatusCode: 500,
	}
)