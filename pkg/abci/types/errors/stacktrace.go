package errors

import "github.com/pkg/errors"





// stackTrace returns the first found stack trace frame carried by given error
// or any wrapped error. It returns nil if no stack trace is found.
func stackTrace(err error) errors.StackTrace {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	for {
		if st, ok := err.(stackTracer); ok {
			return st.StackTrace()
		}

		if c, ok := err.(causer); ok {
			err = c.Cause()
		} else {
			return nil
		}
	}
}