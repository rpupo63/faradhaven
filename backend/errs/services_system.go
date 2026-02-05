package errs

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Configuration & Environment Errors
var (
	ErrConfigMissing       = errors.New("configuration missing")
	ErrConfigInvalid       = errors.New("configuration invalid")
	ErrRegionNotSupported  = errors.New("region not supported")
	ErrSecretMismatch      = errors.New("secret mismatch")
	ErrEnvironmentVariable = errors.New("environment variable error")
)

// Data Consistency & Integrity Errors
var (
	ErrSchemaVersionMismatch = errors.New("schema version mismatch")
	ErrClockSkew             = errors.New("clock skew detected")
	ErrPartialFailure        = errors.New("partial failure")
	ErrDataCorruption        = errors.New("data corruption")
)

// Computational & Business Logic Errors
var (
	ErrDivideByZero = errors.New("divide by zero")
	ErrOverflow     = errors.New("arithmetic overflow")
	ErrUnderflow    = errors.New("arithmetic underflow")
	ErrNilPointer   = errors.New("nil pointer dereference")
	ErrInvalidInput = errors.New("invalid input")
)

// Serialization & Encoding Errors
var (
	ErrCircularStructure = errors.New("circular structure detected")
	ErrBase64Decode      = errors.New("base64 decode error")
	ErrCharsetConversion = errors.New("charset conversion error")
	ErrJSONMarshal       = errors.New("JSON marshal error")
	ErrJSONUnmarshal     = errors.New("JSON unmarshal error")
)

// Configuration & Environment Error Constructors
func NewConfigError(configName string, cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrConfigMissing,
		Details:    fmt.Sprintf("Configuration error for %s", configName),
		Cause:      cause,
	}
}

func NewEnvironmentVariableError(varName string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrEnvironmentVariable,
		Details:    fmt.Sprintf("Environment variable %s is not set or invalid", varName),
		Field:      varName,
	}
}

func NewRegionNotSupportedError(region string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusBadRequest,
		err:        ErrRegionNotSupported,
		Details:    fmt.Sprintf("Region %s is not supported", region),
		Field:      "region",
	}
}

// Data Consistency & Integrity Error Constructors
func NewSchemaVersionMismatchError(service string, expected, actual string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusBadRequest,
		err:        ErrSchemaVersionMismatch,
		Details:    fmt.Sprintf("Schema version mismatch in %s service: expected %s, got %s", service, expected, actual),
		Field:      "schema_version",
	}
}

func NewClockSkewError(service string, skew time.Duration) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusBadRequest,
		err:        ErrClockSkew,
		Details:    fmt.Sprintf("Clock skew detected in %s service: %v", service, skew),
		Field:      "clock",
	}
}

func NewPartialFailureError(operation string, failedSteps []string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrPartialFailure,
		Details:    fmt.Sprintf("Partial failure in %s operation. Failed steps: %v", operation, failedSteps),
		Field:      "partial_failure",
	}
}

// Computational & Business Logic Error Constructors
func NewDivideByZeroError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrDivideByZero,
		Details:    fmt.Sprintf("Divide by zero error in %s", operation),
		Field:      "arithmetic",
	}
}

func NewOverflowError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrOverflow,
		Details:    fmt.Sprintf("Arithmetic overflow in %s", operation),
		Field:      "arithmetic",
	}
}

func NewNilPointerError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrNilPointer,
		Details:    fmt.Sprintf("Nil pointer dereference in %s", operation),
		Field:      "pointer",
	}
}

// Serialization & Encoding Error Constructors
func NewCircularStructureError(operation string) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrCircularStructure,
		Details:    fmt.Sprintf("Circular structure detected in %s", operation),
		Field:      "serialization",
	}
}

func NewBase64DecodeError(operation string, cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusBadRequest,
		err:        ErrBase64Decode,
		Details:    fmt.Sprintf("Base64 decode error in %s", operation),
		Cause:      cause,
		Field:      "encoding",
	}
}

func NewJSONMarshalError(operation string, cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusInternalServerError,
		err:        ErrJSONMarshal,
		Details:    fmt.Sprintf("JSON marshal error in %s", operation),
		Cause:      cause,
		Field:      "json",
	}
}

func NewJSONUnmarshalError(operation string, cause error) *ApiErr {
	return &ApiErr{
		StatusCode: http.StatusBadRequest,
		err:        ErrJSONUnmarshal,
		Details:    fmt.Sprintf("JSON unmarshal error in %s", operation),
		Cause:      cause,
		Field:      "json",
	}
}

// Type Checkers
func IsConfigError(err error) bool {
	return errors.Is(err, ErrConfigMissing) || errors.Is(err, ErrConfigInvalid)
}

func IsEnvironmentVariableError(err error) bool {
	return errors.Is(err, ErrEnvironmentVariable)
}

func IsSchemaVersionMismatchError(err error) bool {
	return errors.Is(err, ErrSchemaVersionMismatch)
}

func IsClockSkewError(err error) bool {
	return errors.Is(err, ErrClockSkew)
}

func IsPartialFailureError(err error) bool {
	return errors.Is(err, ErrPartialFailure)
}

func IsDivideByZeroError(err error) bool {
	return errors.Is(err, ErrDivideByZero)
}

func IsOverflowError(err error) bool {
	return errors.Is(err, ErrOverflow)
}

func IsNilPointerError(err error) bool {
	return errors.Is(err, ErrNilPointer)
}

func IsCircularStructureError(err error) bool {
	return errors.Is(err, ErrCircularStructure)
}

func IsBase64DecodeError(err error) bool {
	return errors.Is(err, ErrBase64Decode)
}

func IsJSONMarshalError(err error) bool {
	return errors.Is(err, ErrJSONMarshal)
}

func IsJSONUnmarshalError(err error) bool {
	return errors.Is(err, ErrJSONUnmarshal)
}
