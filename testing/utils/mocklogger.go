package testingutils

import "github.com/sirupsen/logrus"

// MockLogger is a mock implementation of logrus.FieldLogger for testing
type MockLogger struct{}

func (m *MockLogger) WithField(key string, value interface{}) *logrus.Entry {
	return &logrus.Entry{}
}

func (m *MockLogger) WithFields(fields logrus.Fields) *logrus.Entry {
	return &logrus.Entry{}
}

func (m *MockLogger) WithError(err error) *logrus.Entry {
	return &logrus.Entry{}
}

func (m *MockLogger) Debugf(format string, args ...interface{})   {}
func (m *MockLogger) Infof(format string, args ...interface{})    {}
func (m *MockLogger) Printf(format string, args ...interface{})   {}
func (m *MockLogger) Warnf(format string, args ...interface{})    {}
func (m *MockLogger) Warningf(format string, args ...interface{}) {}
func (m *MockLogger) Errorf(format string, args ...interface{})   {}
func (m *MockLogger) Fatalf(format string, args ...interface{})   {}
func (m *MockLogger) Panicf(format string, args ...interface{})   {}

func (m *MockLogger) Debug(args ...interface{})   {}
func (m *MockLogger) Info(args ...interface{})    {}
func (m *MockLogger) Print(args ...interface{})   {}
func (m *MockLogger) Warn(args ...interface{})    {}
func (m *MockLogger) Warning(args ...interface{}) {}
func (m *MockLogger) Error(args ...interface{})   {}
func (m *MockLogger) Fatal(args ...interface{})   {}
func (m *MockLogger) Panic(args ...interface{})   {}

func (m *MockLogger) Debugln(args ...interface{})   {}
func (m *MockLogger) Infoln(args ...interface{})    {}
func (m *MockLogger) Println(args ...interface{})   {}
func (m *MockLogger) Warnln(args ...interface{})    {}
func (m *MockLogger) Warningln(args ...interface{}) {}
func (m *MockLogger) Errorln(args ...interface{})   {}
func (m *MockLogger) Fatalln(args ...interface{})   {}
func (m *MockLogger) Panicln(args ...interface{})   {}
