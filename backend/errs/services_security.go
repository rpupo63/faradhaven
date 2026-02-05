package errs

import (
	"errors"
	"fmt"
	"net/http"
)

// Dependency & Service Discovery Errors
var (
	ErrServiceUnreachable = errors.New("service unreachable")
	ErrServiceDiscovery   = errors.New("service discovery failed")
	ErrFeatureFlagBackend = errors.New("feature flag backend unavailable")
	ErrLoadBalancer       = errors.New("load balancer error")
	ErrCircuitBreakerOpen = errors.New("circuit breaker open")
)

// Security & Compliance Errors
var (
	ErrSQLInjection       = errors.New("SQL injection attempt")
	ErrRequestForgery     = errors.New("request forgery detected")
	ErrSecretsLeaked      = errors.New("secrets leaked")
	ErrPIIBreach          = errors.New("PII breach detected")
	ErrUnauthorizedAccess = errors.New("unauthorized access")
)

// Observability & Telemetry Errors
var (
	ErrTracerExporter     = errors.New("tracer exporter down")
	ErrLogPipeline        = errors.New("log pipeline error")
	ErrMetricsCardinality = errors.New("metrics cardinality explosion")
	ErrTelemetryFailure   = errors.New("telemetry failure")
)

// Dependency & Service Discovery Error Constructors
func NewServiceUnreachableError(service string, cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusServiceUnavailable,
		err:        ErrServiceUnreachable,
		Details:    fmt.Sprintf("Service %s is unreachable", service),
		Cause:      cause,
		Field:      "service_discovery",
	}
}

func NewServiceDiscoveryError(service string, cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusServiceUnavailable,
		err:        ErrServiceDiscovery,
		Details:    fmt.Sprintf("Service discovery failed for %s", service),
		Cause:      cause,
		Field:      "service_discovery",
	}
}

// Security & Compliance Error Constructors
func NewSQLInjectionError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusBadRequest,
		err:        ErrSQLInjection,
		Details:    fmt.Sprintf("SQL injection attempt detected in %s", operation),
		Field:      "security",
	}
}

func NewRequestForgeryError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusForbidden,
		err:        ErrRequestForgery,
		Details:    fmt.Sprintf("Request forgery detected in %s", operation),
		Field:      "security",
	}
}

func NewSecretsLeakedError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrSecretsLeaked,
		Details:    fmt.Sprintf("Secrets leaked in %s", operation),
		Field:      "security",
	}
}

func NewPIIBreachError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusBadRequest,
		err:        ErrPIIBreach,
		Details:    fmt.Sprintf("PII breach detected in %s", operation),
		Field:      "compliance",
	}
}

// Observability & Telemetry Error Constructors
func NewTracerExporterError(cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrTracerExporter,
		Details:    "Tracer exporter is down",
		Cause:      cause,
		Field:      "telemetry",
	}
}

func NewLogPipelineError(cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrLogPipeline,
		Details:    "Log pipeline error",
		Cause:      cause,
		Field:      "telemetry",
	}
}

func NewMetricsCardinalityError(metric string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrMetricsCardinality,
		Details:    fmt.Sprintf("Metrics cardinality explosion for %s", metric),
		Field:      "telemetry",
	}
}

// Type Checkers
func IsServiceUnreachableError(err error) bool {
	return errors.Is(err, ErrServiceUnreachable)
}

func IsServiceDiscoveryError(err error) bool {
	return errors.Is(err, ErrServiceDiscovery)
}

func IsSQLInjectionError(err error) bool {
	return errors.Is(err, ErrSQLInjection)
}

func IsRequestForgeryError(err error) bool {
	return errors.Is(err, ErrRequestForgery)
}

func IsSecretsLeakedError(err error) bool {
	return errors.Is(err, ErrSecretsLeaked)
}

func IsPIIBreachError(err error) bool {
	return errors.Is(err, ErrPIIBreach)
}

func IsTracerExporterError(err error) bool {
	return errors.Is(err, ErrTracerExporter)
}

func IsLogPipelineError(err error) bool {
	return errors.Is(err, ErrLogPipeline)
}

func IsMetricsCardinalityError(err error) bool {
	return errors.Is(err, ErrMetricsCardinality)
}
