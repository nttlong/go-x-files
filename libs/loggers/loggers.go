package loggers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
func SetupLoggers(appDir string, logDir string) {
	// if logdir start with './' combine with appdir
	absLogDir := logDir
	if logDir[0:2] == "./" {
		absLogDir = filepath.Join(appDir, logDir[1:])
	}
	// create log directory if not exists
	if _, err := os.Stat(absLogDir); os.IsNotExist(err) {
		os.MkdirAll(absLogDir, 0755)
	}
	// create log file
	logFile := filepath.Join(absLogDir, "app.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Failed to create log file:", err)
	}
	// set log output to file
	log.SetOutput(f)
}
