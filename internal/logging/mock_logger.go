package logging

var _ Logger = &DummyLogger{}

// Dummy logger methods do nothing, and is used for testing
type DummyLogger struct{}

func (l *DummyLogger) Emergencyf(f string, v ...interface{}) {}

func (l *DummyLogger) Alertf(f string, v ...interface{}) {}

func (l *DummyLogger) Criticalf(f string, v ...interface{}) {}

func (l *DummyLogger) Errorf(f string, v ...interface{}) {}

func (l *DummyLogger) Warningf(f string, v ...interface{}) {}

func (l *DummyLogger) Noticef(f string, v ...interface{}) {}

func (l *DummyLogger) Informationalf(f string, v ...interface{}) {}

func (l *DummyLogger) Debugf(f string, v ...interface{}) {}
