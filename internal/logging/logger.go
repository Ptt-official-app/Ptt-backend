package logging

import (
	"log"
)

var DefaultLogger Logger = NewLogger()

// Logger is the interface that wraps the basic Logging methods.
//
// Logger follows the Syslog Message Severities in RFC5424.
type Logger interface {
	// Emergency: system is unusable
	Emergencyf(f string, v ...interface{})
	// Alert: action must be taken immediately
	Alertf(f string, v ...interface{})
	// Critical: critical conditions
	Criticalf(f string, v ...interface{})
	// Error: error conditions
	Errorf(f string, v ...interface{})
	// Warning: warning conditions
	Warningf(f string, v ...interface{})
	// Notice: normal but significant condition
	Noticef(f string, v ...interface{})
	// Informational: informational messages
	Informationalf(f string, v ...interface{})
	// Debug: debug-level messages
	Debugf(f string, v ...interface{})
}

type logger struct{}

func NewLogger() Logger {
	return &logger{}
}

// Emergencyf implements Logger.Emergencyf by printing Emergency level messages to standard output.
func (l *logger) Emergencyf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Alertf implements Logger.Alertf by printing Alert level messages to standard output.
func (l *logger) Alertf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Criticalf implements Logger.Criticalf by printing Critical level messages to standard output.
func (l *logger) Criticalf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Errorf implements Logger.Errorf by printing Error level messages to standard output.
func (l *logger) Errorf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Warningf implements Logger.Warningf by printing Warning level messages to standard output.
func (l *logger) Warningf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Noticef implements Logger.Noticef by printing Notice level messages to standard output.
func (l *logger) Noticef(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Informationalf implements Logger.Informationalf by printing Informational level messages to standard output.
func (l *logger) Informationalf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Debugf implements Logger.Debugf by printing Debug level messages to standard output.
func (l *logger) Debugf(f string, v ...interface{}) {
	log.Printf(f, v...)
}
