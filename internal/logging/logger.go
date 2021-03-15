package logging

import (
	"log"
)

// Logger provides the basic Logging methods, following the Syslog Message Severities defined in RFC5424.
// https://tools.ietf.org/html/rfc5424#section-6.2.1
type Logger struct {}

func NewLogger() *Logger {
	return &Logger{}
}

// Emergencyf prints Emergency level messages to standard output.
// RFC5424 definition: "system is unusable"
func (l *Logger) Emergencyf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Alertf prints Alert level messages to standard output.
// RFC5424 definition: "action must be taken immediately"
func (l *Logger) Alertf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Criticalf prints Critical level messages to standard output.
// RFC5424 definition: "critical conditions"
func (l *Logger) Criticalf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Errorf prints Error level messages to standard output.
// RFC5424 definition: "error conditions"
func (l *Logger) Errorf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Warningf prints Warning level messages to standard output.
// RFC5424 definition: "warning conditions"
func (l *Logger) Warningf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Noticef prints Notice level messages to standard output.
// RFC5424 definition: "normal but significant condition"
func (l *Logger) Noticef(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Informationalf prints Informational level messages to standard output.
// RFC5424 definition: "informational messages"
func (l *Logger) Informationalf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Debugf prints Debug level messages to standard output.
// RFC5424 definition: "debug-level messages"
func (l *Logger) Debugf(f string, v ...interface{}) {
	log.Printf(f, v...)
}
