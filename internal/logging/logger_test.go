package logging

import (
	"log"
	"regexp"
	"strings"
	"testing"
)

const (
	patternLogPrefix = `\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2}`
	testMessage      = "test message"
)

var (
	testOutput = &strings.Builder{}
)

func init() {
	log.SetOutput(testOutput)
}

func TestLoggerEmergencyf(t *testing.T) {
	testOutput.Reset()
	DefaultLogger.Emergencyf(testMessage)
	b, err := regexp.MatchString(patternLogPrefix, testOutput.String())

	if err != nil {
		t.Errorf("MatchString expect nil, got %v", err)
	} else if !b {
		t.Errorf(`string: %q does not match pattern "%s"`, testOutput.String(), patternLogPrefix)
	}

	b = strings.Contains(testOutput.String(), testMessage)

	if !b {
		t.Errorf("expect %q contains %q but not", testOutput.String(), testMessage)
	}
}

func TestLoggerAlertf(t *testing.T) {
	testOutput.Reset()
	DefaultLogger.Alertf(testMessage)
	b, err := regexp.MatchString(patternLogPrefix, testOutput.String())

	if err != nil {
		t.Errorf("MatchString expect nil, got %v", err)
	} else if !b {
		t.Errorf(`string: %q does not match pattern "%s"`, testOutput.String(), patternLogPrefix)
	}

	b = strings.Contains(testOutput.String(), testMessage)

	if !b {
		t.Errorf("expect %q contains %q but not", testOutput.String(), testMessage)
	}
}

func TestLoggerCriticalf(t *testing.T) {
	testOutput.Reset()
	DefaultLogger.Criticalf(testMessage)
	b, err := regexp.MatchString(patternLogPrefix, testOutput.String())

	if err != nil {
		t.Errorf("MatchString expect nil, got %v", err)
	} else if !b {
		t.Errorf(`string: %q does not match pattern "%s"`, testOutput.String(), patternLogPrefix)
	}

	b = strings.Contains(testOutput.String(), testMessage)

	if !b {
		t.Errorf("expect %q contains %q but not", testOutput.String(), testMessage)
	}
}

func TestLoggerErrorf(t *testing.T) {
	testOutput.Reset()
	DefaultLogger.Errorf(testMessage)
	b, err := regexp.MatchString(patternLogPrefix, testOutput.String())

	if err != nil {
		t.Errorf("MatchString expect nil, got %v", err)
	} else if !b {
		t.Errorf(`string: %q does not match pattern "%s"`, testOutput.String(), patternLogPrefix)
	}

	b = strings.Contains(testOutput.String(), testMessage)

	if !b {
		t.Errorf("expect %q contains %q but not", testOutput.String(), testMessage)
	}
}

func TestLoggerWarningf(t *testing.T) {
	testOutput.Reset()
	DefaultLogger.Warningf(testMessage)
	b, err := regexp.MatchString(patternLogPrefix, testOutput.String())

	if err != nil {
		t.Errorf("MatchString expect nil, got %v", err)
	} else if !b {
		t.Errorf(`string: %q does not match pattern "%s"`, testOutput.String(), patternLogPrefix)
	}

	b = strings.Contains(testOutput.String(), testMessage)

	if !b {
		t.Errorf("expect %q contains %q but not", testOutput.String(), testMessage)
	}
}

func TestLoggerNoticef(t *testing.T) {
	testOutput.Reset()
	DefaultLogger.Noticef(testMessage)
	b, err := regexp.MatchString(patternLogPrefix, testOutput.String())

	if err != nil {
		t.Errorf("MatchString expect nil, got %v", err)
	} else if !b {
		t.Errorf(`string: %q does not match pattern "%s"`, testOutput.String(), patternLogPrefix)
	}

	b = strings.Contains(testOutput.String(), testMessage)

	if !b {
		t.Errorf("expect %q contains %q but not", testOutput.String(), testMessage)
	}
}

func TestLoggerInformationalf(t *testing.T) {
	testOutput.Reset()
	DefaultLogger.Informationalf(testMessage)
	b, err := regexp.MatchString(patternLogPrefix, testOutput.String())

	if err != nil {
		t.Errorf("MatchString expect nil, got %v", err)
	} else if !b {
		t.Errorf(`string: %q does not match pattern "%s"`, testOutput.String(), patternLogPrefix)
	}

	b = strings.Contains(testOutput.String(), testMessage)

	if !b {
		t.Errorf("expect %q contains %q but not", testOutput.String(), testMessage)
	}
}

func TestLoggerDebugf(t *testing.T) {
	testOutput.Reset()
	DefaultLogger.Debugf(testMessage)
	b, err := regexp.MatchString(patternLogPrefix, testOutput.String())

	if err != nil {
		t.Errorf("MatchString expect nil, got %v", err)
	} else if !b {
		t.Errorf(`string: %q does not match pattern "%s"`, testOutput.String(), patternLogPrefix)
	}

	b = strings.Contains(testOutput.String(), testMessage)

	if !b {
		t.Errorf("expect %q contains %q but not", testOutput.String(), testMessage)
	}
}
