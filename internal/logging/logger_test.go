package logging

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// newFile returns a temporally file
func newFile(testName string, t *testing.T) (f *os.File) {
	f, err := os.CreateTemp(t.TempDir(), "_Ptt-backend_"+testName)
	if err != nil {
		t.Fatalf("TempFile %s: %s", testName, err)
	}
	return
}

var settingList []LogSetting = []LogSetting{
	Emergency,
	Alert,
	Critical,
	Error,
	Warning,
	Notice,
	Informational,
	Debug,
}

func testLoggerLevel(t *testing.T, targetLevel uint) {
	testName := fmt.Sprintf("testing level %v", targetLevel)
	tempFile := newFile(testName, t)
	t.Cleanup(func() {
		if err := tempFile.Close(); err != nil {
			t.Errorf("close temp log file: %v", err)
		}
	})
	testLogger := &logger{
		output: tempFile,
		level:  targetLevel,
	}
	testLogger.Emergencyf("")
	testLogger.Alertf("")
	testLogger.Criticalf("")
	testLogger.Errorf("")
	testLogger.Warningf("")
	testLogger.Noticef("")
	testLogger.Informationalf("")
	testLogger.Debugf("")
	_, err := tempFile.Seek(0, 0)
	if err != nil {
		log.Fatalf("Get error %s\n", err.Error())
	}
	reader := bufio.NewReader(tempFile)
	idx := uint(0)
	for line, _, err := reader.ReadLine(); ; line, _, err = reader.ReadLine() {
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Get error %s", err.Error())
		}
		result := strings.Split(string(line), " ")
		// result[0]: 2006-01-02
		// result[1]: 15:04:05
		if idx > targetLevel {
			t.Fatalf("Unexpect LogLevel %s", result[2])
		}
		if result[2] != settingList[idx].Name {
			t.Fatalf("Expect log level %s, but get %s", settingList[idx].Name, result[2])
		}
		idx++
	}

	if idx-1 != targetLevel {
		t.Fatalf("Missing Log Level %d: \"%s\"", idx, settingList[idx].Name)
	}
}

func TestLogger(t *testing.T) {
	testLoggerLevel(t, 7)
	testLoggerLevel(t, 6)
	testLoggerLevel(t, 5)
	testLoggerLevel(t, 4)
	testLoggerLevel(t, 3)
	testLoggerLevel(t, 2)
	testLoggerLevel(t, 1)
	testLoggerLevel(t, 0)
	t.Setenv("LOG_LEVEL", "4")
	testLogger := NewLogger()
	if testLogger.(*logger).level != 4 {
		t.Fatalf("Expect logger level is 4, get %d", testLogger.(*logger).level)
	}
	t.Setenv("LOG_LEVEL", "")
	testLogger = NewLogger()
	if testLogger.(*logger).level != 7 {
		t.Fatalf("Expect logger level is 7, get %d", testLogger.(*logger).level)
	}
}
