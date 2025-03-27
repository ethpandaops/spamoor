package logscope

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

type LogScope struct {
	options      *ScopeOptions
	logger       *logrus.Logger
	parentLogger *logrus.Logger
	parentFields logrus.Fields
	memBuffer    *logMemBuffer
}

type ScopeOptions struct {
	Parent     logrus.FieldLogger
	BufferSize uint64
}

type logForwarder struct {
	logger *LogScope
}

func NewLogger(options *ScopeOptions) *LogScope {
	if options == nil {
		options = &ScopeOptions{}
	}

	if options.BufferSize == 0 {
		options.BufferSize = 100
	}

	logger := &LogScope{
		options: options,
		logger:  logrus.New(),
	}

	logger.logger.SetOutput(io.Discard)
	logger.logger.SetLevel(logrus.DebugLevel)

	if options.Parent != nil {
		tmpEntry := options.Parent.WithFields(logrus.Fields{})
		logger.parentLogger = tmpEntry.Logger
		logger.parentFields = tmpEntry.Data
		logger.logger.AddHook(&logForwarder{
			logger: logger,
		})
	}

	logger.memBuffer = newLogMemBuffer(logger, options.BufferSize)
	logger.logger.AddHook(logger.memBuffer)

	return logger
}

func (lf *logForwarder) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (lf *logForwarder) Fire(entry *logrus.Entry) error {
	entry2 := entry.Dup()
	entry2.Logger = lf.logger.parentLogger

	entry2.WithFields(lf.logger.parentFields).Log(entry.Level, entry.Message)

	return nil
}

func (ls *LogScope) GetLogger() *logrus.Logger {
	return ls.logger
}

func (ls *LogScope) Flush() {
}

func (ls *LogScope) GetLogEntryCount() int {
	if ls.memBuffer != nil {
		return ls.memBuffer.GetLogEntryCount()
	}

	return 0
}

func (ls *LogScope) GetLogEntries(from time.Time, limit int) []*logrus.Entry {
	if ls.memBuffer != nil {
		return ls.memBuffer.GetLogEntries(from, limit)
	}

	return nil
}
