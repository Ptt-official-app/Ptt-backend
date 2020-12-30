package main

import (
	"log"
)

type Logger struct{}

var logger = Logger{}

// According RFC5424
// Emergency: system is unusable
func (l *Logger) Emergencyf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Alert: action must be taken immediately
func (l *Logger) Alertf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Critical: critical conditions
func (l *Logger) Criticalf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Error: error conditions
func (l *Logger) Errorf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Warning: warning conditions
func (l *Logger) Warningf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Notice: normal but significant condition
func (l *Logger) Noticef(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Informational: informational messages
func (l *Logger) Informationalf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

// Debug: debug-level messages
func (l *Logger) Debugf(f string, v ...interface{}) {
	log.Printf(f, v...)
}
