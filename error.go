package werror

import (
	"fmt"
	"runtime"
	"strings"

	"google.golang.org/grpc/codes"
)

// Error is a custom error type that includes an error reason.
// If you are using protocol buffers, define the type of this field as pb.Reason.
//
// shoudReport is used to indicate whether the error should be reported.
// callers is used to identify the location where the error occurred.
//
// code is used to indicate the gRPC error code.
// This field is used to indicate the error code when using gRPC only.
// If you are not using gRPC, you can ignore this field.
// see: https://pkg.go.dev/google.golang.org/grpc/codes
//
// reason is used to indicate the reason for your application error.
// If you are using protocol buffers, define the type of this field as pb.Reason.
//
// Error implements ErrorWithCallers interface of bugsnag-go.
// see: https://github.com/bugsnag/bugsnag-go/blob/b31bbecd4eb6e307dd7738f729ab51973244d903/v2/errors/error.go#L26-L29
type Error struct {
	cause        error
	err          error
	messages     []string
	reason       ErrorReason
	shouldReport bool
	code         codes.Code
	callers      []uintptr
}

type ErrorReason interface {
	String() string
}

type emptyReason struct{}

func (e emptyReason) String() string {
	return ""
}

type Annotator func(error) error

// New creates a new error with the given text.
func New(text string) error {
	err := &Error{
		messages:     []string{text},
		reason:       emptyReason{},
		code:         codes.Unknown,
		shouldReport: true,
	}
	return WithCallers(1)(err)
}

// NewFromStandardError creates a new *Error from the given error.
// If the error is not *Error, create a new *Error.
// we do not specify err field because bugsnag needs ErrorWithCallers interface.
func NewFromStandardError(err error) *Error {
	return &Error{
		cause:        err,
		messages:     []string{err.Error()},
		code:         codes.Unknown,
		reason:       emptyReason{},
		shouldReport: true,
	}
}

// Wrap annotates an error with the given annotators.
// If the error is nil, it returns nil.
func Wrap(err error, annotators ...Annotator) error {
	if err == nil {
		return nil
	}

	for _, a := range annotators {
		err = a(err)
	}

	if _, ok := err.(*Error); !ok {
		err = WithCallers(1)(err)
	}

	return err
}

// WithCallers annotates an error with the stack trace.
// The offset parameter is used to identify the location where the error occurred.
func WithCallers(offset int) Annotator {
	return func(err error) error {
		if werr, ok := err.(*Error); ok {
			if werr.callers != nil {
				return werr
			}

			werr.callers = createCallers(offset + 1)
			return werr
		}

		werr := NewFromStandardError(err)
		werr.callers = createCallers(offset + 1)
		return werr
	}
}

// Error returns the error message.
// The error message includes the reason, code, and message.
func (e *Error) Error() string {
	return fmt.Sprintf("reason: %s, code: %d, message: %s", Reason(e).String(), Code(e), Message(e))
}

// Unwrap returns the original error.
func (e *Error) Unwrap() error {
	return e.err
}

// Callers returns the stack trace of the error.
// This is used by bugsnag to display the stack trace.
// see: https://github.com/bugsnag/bugsnag-go/blob/b31bbecd4eb6e307dd7738f729ab51973244d903/v2/errors/error.go#L26-L29
func (e *Error) Callers() []uintptr {
	if e.callers == nil {
		return []uintptr{}
	}
	return e.callers
}

// ShouldReport returns whether the error should be reported.
func ShouldReport(err error) bool {
	if r, ok := err.(*Error); ok {
		return r.shouldReport
	}
	return true
}

// WithCode annotates an error with the gRPC error code.
func WithCode(code codes.Code) Annotator {
	return func(err error) error {
		if werr, ok := err.(*Error); ok {
			werr.code = code
			return werr
		}

		werr := NewFromStandardError(err)
		werr.code = code

		return WithCallers(1)(werr)
	}
}

// WithReason annotates an error with the reason of the application error.
func WithReason(reason interface{}) Annotator {
	if errReason, ok := reason.(ErrorReason); ok {
		return func(err error) error {
			if werr, ok := err.(*Error); ok {
				werr.reason = errReason
				return werr
			}

			werr := NewFromStandardError(err)
			werr.reason = errReason

			return WithCallers(1)(werr)
		}
	}

	return func(err error) error {
		return err
	}
}

// WithIgnoreReport annotates an error with the shouldReport field set to false.
func WithIgnoreReport() Annotator {
	return func(err error) error {
		if werr, ok := err.(*Error); ok {
			werr.shouldReport = false
			return werr
		}

		werr := NewFromStandardError(err)
		werr.shouldReport = false

		return WithCallers(1)(werr)
	}
}

func WithMessage(message string) Annotator {
	return func(err error) error {
		if werr, ok := err.(*Error); ok {
			werr.messages = append(werr.messages, message)
			return werr
		}

		werr := NewFromStandardError(err)
		werr.messages = append(werr.messages, message)

		return WithCallers(1)(werr)
	}
}

// Code returns the gRPC error code.
// If the error does not have a gRPC error code, it returns codes.Unknown.
func Code(err error) codes.Code {
	if werr, ok := err.(*Error); ok {
		return werr.code
	}

	return codes.Unknown
}

// Reason returns the reason of the application error.
func Reason(err error) ErrorReason {
	if werr, ok := err.(*Error); ok {
		return werr.reason
	}

	return emptyReason{}
}

func Message(err error) string {
	if werr, ok := err.(*Error); ok {
		return strings.Join(werr.messages, " : ")
	}

	return ""
}

func createCallers(offset int) []uintptr {
	pcs := make([]uintptr, 100)
	n := runtime.Callers(offset+2, pcs[:])
	return pcs[:n]
}
