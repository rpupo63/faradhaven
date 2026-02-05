package errs

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Networking & Transport Errors
var (
	ErrDNSResolution      = errors.New("DNS resolution failed")
	ErrTCPTimeout         = errors.New("TCP connection timeout")
	ErrTLSHandshake       = errors.New("TLS handshake failed")
	ErrConnectionReset    = errors.New("connection reset")
	ErrProxyBlocked       = errors.New("proxy/firewall blocked")
	ErrNetworkUnreachable = errors.New("network unreachable")
)

// Resource Exhaustion Errors
var (
	ErrOutOfMemory         = errors.New("out of memory")
	ErrCPUExhausted        = errors.New("CPU exhausted")
	ErrFileDescriptorLimit = errors.New("file descriptor limit exceeded")
	ErrDiskSpaceFull       = errors.New("disk space full")
)

// Timeout & Cancellation Errors
var (
	ErrContextDeadline    = errors.New("context deadline exceeded")
	ErrClientDisconnected = errors.New("client disconnected")
	ErrComputationTimeout = errors.New("computation timeout")
	ErrRequestTimeout     = errors.New("request timeout")
	ErrTimeout            = errors.New("timeout")
)

// Deployment & Runtime Platform Errors
var (
	ErrBadRollout     = errors.New("bad rollout")
	ErrSidecarCrash   = errors.New("sidecar crash")
	ErrContainerImage = errors.New("container image error")
	ErrPodUnhealthy   = errors.New("pod unhealthy")
)

// Networking & Transport Error Constructors
func NewDNSError(host string, cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusServiceUnavailable,
		err:        ErrDNSResolution,
		Details:    fmt.Sprintf("DNS resolution failed for %s", host),
		Cause:      cause,
	}
}

func NewTCPTimeoutError(host string, timeout time.Duration) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusServiceUnavailable,
		err:        ErrTCPTimeout,
		Details:    fmt.Sprintf("TCP connection timeout to %s after %v", host, timeout),
		Field:      "connection",
	}
}

func NewTLSHandshakeError(host string, cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusServiceUnavailable,
		err:        ErrTLSHandshake,
		Details:    fmt.Sprintf("TLS handshake failed for %s", host),
		Cause:      cause,
	}
}

func NewConnectionResetError(host string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusServiceUnavailable,
		err:        ErrConnectionReset,
		Details:    fmt.Sprintf("Connection reset by peer for %s", host),
		Field:      "connection",
	}
}

// Resource Exhaustion Error Constructors
func NewOutOfMemoryError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrOutOfMemory,
		Details:    fmt.Sprintf("Out of memory during %s operation", operation),
		Field:      "memory",
	}
}

func NewCPUExhaustedError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrCPUExhausted,
		Details:    fmt.Sprintf("CPU exhausted during %s operation", operation),
		Field:      "cpu",
	}
}

func NewFileDescriptorLimitError() *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrFileDescriptorLimit,
		Details:    "File descriptor limit exceeded",
		Field:      "file_descriptors",
	}
}

// Timeout & Cancellation Error Constructors
func NewContextDeadlineError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusRequestTimeout,
		err:        ErrContextDeadline,
		Details:    fmt.Sprintf("Context deadline exceeded for %s", operation),
		Field:      "timeout",
	}
}

func NewClientDisconnectedError() *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusRequestTimeout,
		err:        ErrClientDisconnected,
		Details:    "Client disconnected during request",
		Field:      "client",
	}
}

func NewComputationTimeoutError(operation string, timeout time.Duration) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusRequestTimeout,
		err:        ErrComputationTimeout,
		Details:    fmt.Sprintf("Computation timeout for %s after %v", operation, timeout),
		Field:      "computation",
	}
}

// Deployment & Runtime Platform Error Constructors
func NewBadRolloutError(service string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusServiceUnavailable,
		err:        ErrBadRollout,
		Details:    fmt.Sprintf("Bad rollout for %s service", service),
		Field:      "deployment",
	}
}

func NewSidecarCrashError(service string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusServiceUnavailable,
		err:        ErrSidecarCrash,
		Details:    fmt.Sprintf("Sidecar crash for %s service", service),
		Field:      "deployment",
	}
}

func NewContainerImageError(service string, cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrContainerImage,
		Details:    fmt.Sprintf("Container image error for %s service", service),
		Cause:      cause,
		Field:      "deployment",
	}
}

// Type Checkers
func IsDNSError(err error) bool {
	return errors.Is(err, ErrDNSResolution)
}

func IsTCPTimeoutError(err error) bool {
	return errors.Is(err, ErrTCPTimeout)
}

func IsTLSHandshakeError(err error) bool {
	return errors.Is(err, ErrTLSHandshake)
}

func IsConnectionResetError(err error) bool {
	return errors.Is(err, ErrConnectionReset)
}

func IsOutOfMemoryError(err error) bool {
	return errors.Is(err, ErrOutOfMemory)
}

func IsCPUExhaustedError(err error) bool {
	return errors.Is(err, ErrCPUExhausted)
}

func IsContextDeadlineError(err error) bool {
	return errors.Is(err, ErrContextDeadline)
}

func IsClientDisconnectedError(err error) bool {
	return errors.Is(err, ErrClientDisconnected)
}

func IsComputationTimeoutError(err error) bool {
	return errors.Is(err, ErrComputationTimeout)
}

func IsBadRolloutError(err error) bool {
	return errors.Is(err, ErrBadRollout)
}

func IsSidecarCrashError(err error) bool {
	return errors.Is(err, ErrSidecarCrash)
}

func IsContainerImageError(err error) bool {
	return errors.Is(err, ErrContainerImage)
}
