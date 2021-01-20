package logging

import (
	"log"
)

type Logger interface {
	Emergencyf(f string, v ...interface{})
	Alertf(f string, v ...interface{})
	Criticalf(f string, v ...interface{})
	Errorf(f string, v ...interface{})
	Warningf(f string, v ...interface{})
	Noticef(f string, v ...interface{})
	Informationalf(f string, v ...interface{})
	Debugf(f string, v ...interface{})
}

type logger struct{}

func NewLogger() Logger {
	return &logger{}
}

func (l *logger) Emergencyf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Alert: action must be taken immediately
func (l *logger) Alertf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Critical: critical conditions
func (l *logger) Criticalf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Error: error conditions
func (l *logger) Errorf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Warning: warning conditions
func (l *logger) Warningf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Notice: normal but significant condition
func (l *logger) Noticef(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Informational: informational messages
func (l *logger) Informationalf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Debug: debug-level messages
func (l *logger) Debugf(f string, v ...interface{}) {
	log.Printf(f, v...)
}
