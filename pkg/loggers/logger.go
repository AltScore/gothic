package loggers

import "fmt"

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type simpleLogger struct {
}

func NewSimple() Logger {
	return &simpleLogger{}
}

func (s simpleLogger) Infof(format string, args ...interface{}) {
	fmt.Print("INFO: ")
	fmt.Printf(format, args...)
	fmt.Println()
}

func (s simpleLogger) Warnf(format string, args ...interface{}) {
	fmt.Print("WARN: ")
	fmt.Printf(format, args...)
	fmt.Println()
}

func (s simpleLogger) Errorf(format string, args ...interface{}) {
	fmt.Print("ERROR: ")
	fmt.Printf(format, args...)
	fmt.Println()
}
