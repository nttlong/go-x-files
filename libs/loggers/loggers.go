package loggers

import (
	"fmt"
	"log"
	"runtime/debug"
)

type tracedError struct {
	message string
	stack   []byte
	cause   error
}

func (e *tracedError) Error() string {
	return fmt.Sprintf("%s\nStack trace:\n%s", e.message, e.stack)
}

func (e *tracedError) Unwrap() error {
	return e.cause
}
func HandlePanic() {
	if r := recover(); r != nil {
		var err error
		switch rt := r.(type) {
		case string:
			err = &tracedError{message: rt, stack: debug.Stack()}
		case error:
			err = &tracedError{message: rt.Error(), stack: debug.Stack(), cause: rt}
		default:
			err = &tracedError{message: fmt.Sprintf("panic: %v", rt), stack: debug.Stack()}
		}
		log.Println(err) // Log the full error with stack trace
		fmt.Println(err) // Print the full error with stack trace
	}
}
