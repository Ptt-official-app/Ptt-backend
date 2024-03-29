package logging

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type LogSetting struct {
	Color string
	Name  string
	Level uint
}

// level defined in RFC 5424
var (
	Emergency LogSetting = LogSetting{
		Color: "\x1b[41m",
		Name:  "EMERGENCY",
		Level: 0,
	}
	Alert LogSetting = LogSetting{
		Color: "\x1b[30;43m",
		Name:  "ALERT",
		Level: 1,
	}
	Critical LogSetting = LogSetting{
		Color: "\x1b[41m",
		Name:  "CRITICAL",
		Level: 2,
	}
	Error LogSetting = LogSetting{
		Color: "\x1b[31;103m",
		Name:  "ERROR",
		Level: 3,
	}
	Warning LogSetting = LogSetting{
		Color: "\x1b[33m",
		Name:  "WARN",
		Level: 4,
	}
	Notice LogSetting = LogSetting{
		Color: "\x1b[32m",
		Name:  "NOTICE",
		Level: 5,
	}
	Informational LogSetting = LogSetting{
		Color: "\x1b[34m",
		Name:  "INFO",
		Level: 6,
	}
	Debug LogSetting = LogSetting{
		Color: "\x1b[32m",
		Name:  "DEBUG",
		Level: 7,
	}
)

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
	// function for apply color and output to target file
	Logf(f string, setting LogSetting, v ...interface{})
}

type logger struct {
	output *os.File
	level  uint
}

func NewLogger() Logger {
	level, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		fmt.Println(err)
		level = 7
	}
	return &logger{
		output: os.Stderr,
		level:  uint(level),
	}
}

// Emergencyf implements Logger.Emergencyf by printing Emergency level messages to Logger.output.
func (l *logger) Emergencyf(f string, v ...interface{}) {
	l.Logf(f, Emergency, v...)
}

// Alertf implements Logger.Alertf by printing Alert level messages to Logger.output.
func (l *logger) Alertf(f string, v ...interface{}) {
	l.Logf(f, Alert, v...)
}

// Criticalf implements Logger.Criticalf by printing Critical level messages to Logger.output.
func (l *logger) Criticalf(f string, v ...interface{}) {
	l.Logf(f, Critical, v...)
}

// Errorf implements Logger.Errorf by printing Error level messages to Logger.output.
func (l *logger) Errorf(f string, v ...interface{}) {
	l.Logf(f, Error, v...)
}

// Warningf implements Logger.Warningf by printing Warning level messages to Logger.output.
func (l *logger) Warningf(f string, v ...interface{}) {
	l.Logf(f, Warning, v...)
}

// Noticef implements Logger.Noticef by printing Notice level messages to Logger.output.
func (l *logger) Noticef(f string, v ...interface{}) {
	l.Logf(f, Notice, v...)
}

// Informationalf implements Logger.Informationalf by printing Informational level messages to Logger.output.
func (l *logger) Informationalf(f string, v ...interface{}) {
	l.Logf(f, Informational, v...)
}

// Debugf implements Logger.Debugf by printing Debug level messages to logger.output.
func (l *logger) Debugf(f string, v ...interface{}) {
	l.Logf(f, Debug, v...)
}

// Logf implements Logger.Logf by printing messages to logger.output with LogSetting.
func (l *logger) Logf(f string, setting LogSetting, v ...interface{}) {
	if l.level >= setting.Level {
		fileinfo, err := l.output.Stat()
		outStr := time.Now().Format("2006-01-02 15:04:05")
		if err != nil || (fileinfo.Mode()&os.ModeCharDevice) == 0 {
			outStr += " " + setting.Name + " "
		} else {
			outStr += " " + setting.Color + setting.Name + "\x1b[0m "
		}

		outStr += fmt.Sprintf(f, v...)
		outStr += "\n"
		_, err = l.output.WriteString(outStr)
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}
}
