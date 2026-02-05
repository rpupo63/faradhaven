package errs

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Third-Party API & LLM Specific Errors
var (
	ErrRateLimitExceeded      = errors.New("rate limit exceeded")
	ErrModelOverloaded        = errors.New("model overloaded")
	ErrContextLengthExceeded  = errors.New("context length exceeded")
	ErrContentPolicyViolation = errors.New("content policy violation")
	ErrBillingQuotaExhausted  = errors.New("billing quota exhausted")
	ErrStreamingChunkDropped  = errors.New("streaming chunk dropped")
	ErrInvalidAPIKey          = errors.New("invalid API key")
	ErrServiceUnavailable     = errors.New("service unavailable")
)

// LLM Service Specific Error Constructors
func NewRateLimitError(service string, retryAfter time.Duration) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusTooManyRequests,
		err:        ErrRateLimitExceeded,
		Details:    fmt.Sprintf("Rate limit exceeded for %s service", service),
		Field:      "rate_limit",
	}
}

func NewModelOverloadedError(service string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusServiceUnavailable,
		err:        ErrModelOverloaded,
		Details:    fmt.Sprintf("Model overloaded for %s service", service),
		Field:      "model_capacity",
	}
}

func NewContextLengthError(service string, maxTokens int) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusBadRequest,
		err:        ErrContextLengthExceeded,
		Details:    fmt.Sprintf("Context length exceeded for %s service (max: %d tokens)", service, maxTokens),
		Field:      "context_length",
	}
}

func NewContentPolicyError(service string, violation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusBadRequest,
		err:        ErrContentPolicyViolation,
		Details:    fmt.Sprintf("Content policy violation in %s service: %s", service, violation),
		Field:      "content_policy",
	}
}

func NewBillingQuotaError(service string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusPaymentRequired,
		err:        ErrBillingQuotaExhausted,
		Details:    fmt.Sprintf("Billing quota exhausted for %s service", service),
		Field:      "billing",
	}
}

func NewStreamingError(service string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrStreamingChunkDropped,
		Details:    fmt.Sprintf("Streaming chunk dropped for %s service", service),
		Field:      "streaming",
	}
}

// Error Type Checkers
func IsRateLimitError(err error) bool {
	return errors.Is(err, ErrRateLimitExceeded)
}

func IsModelOverloadedError(err error) bool {
	return errors.Is(err, ErrModelOverloaded)
}

func IsContextLengthError(err error) bool {
	return errors.Is(err, ErrContextLengthExceeded)
}

func IsContentPolicyError(err error) bool {
	return errors.Is(err, ErrContentPolicyViolation)
}

func IsBillingQuotaError(err error) bool {
	return errors.Is(err, ErrBillingQuotaExhausted)
}

func IsStreamingError(err error) bool {
	return errors.Is(err, ErrStreamingChunkDropped)
}
