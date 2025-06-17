package logscope

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

// LogScope provides a scoped logger with memory buffering and optional forwarding to a parent logger.
// It creates an isolated logging environment that can capture log entries in memory while
// optionally forwarding them to a parent logger with inherited fields.
type LogScope struct {
	options      *ScopeOptions
	logger       *logrus.Logger
	parentLogger *logrus.Logger
	parentFields logrus.Fields
	memBuffer    *logMemBuffer
}

// ScopeOptions defines configuration for creating a LogScope instance.
// It specifies the parent logger for forwarding and buffer size for memory storage.
type ScopeOptions struct {
	Parent     logrus.FieldLogger // Parent logger to forward entries to (optional)
	BufferSize uint64             // Maximum number of log entries to buffer in memory
}

// logForwarder is a logrus hook that forwards log entries to a parent logger.
// It preserves the parent's fields and logger instance while forwarding entries.
type logForwarder struct {
	logger *LogScope
}

// NewLogger creates a new LogScope with the specified options.
// If no options are provided, defaults to a buffer size of 100 entries.
// The logger discards output by default and relies on hooks for processing.
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

// Levels returns all logrus levels, indicating this hook handles all log levels.
// This is required by the logrus.Hook interface.
func (lf *logForwarder) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire forwards a log entry to the parent logger with inherited fields.
// It duplicates the entry, switches to the parent logger instance,
// and logs with the combined parent fields and original message.
func (lf *logForwarder) Fire(entry *logrus.Entry) error {
	entry2 := entry.Dup()
	entry2.Logger = lf.logger.parentLogger

	entry2.WithFields(lf.logger.parentFields).Log(entry.Level, entry.Message)

	return nil
}

// GetLogger returns the underlying logrus logger instance.
// This logger has memory buffering and optional parent forwarding configured.
func (ls *LogScope) GetLogger() *logrus.Logger {
	return ls.logger
}

// Flush is a no-op method provided for interface compatibility.
// The LogScope uses hooks for immediate processing, so no flushing is needed.
func (ls *LogScope) Flush() {
}

// GetLogEntryCount returns the total number of log entries processed by this scope.
// This includes entries that may have been rotated out of the memory buffer.
func (ls *LogScope) GetLogEntryCount() int {
	if ls.memBuffer != nil {
		return ls.memBuffer.GetLogEntryCount()
	}

	return 0
}

// GetLogEntries retrieves log entries from the memory buffer with optional filtering.
// The from parameter filters entries to those after the specified time (exclusive).
// The limit parameter restricts the number of returned entries (0 = no limit).
// Returns entries in chronological order, with older entries first.
func (ls *LogScope) GetLogEntries(from time.Time, limit int) []*logrus.Entry {
	if ls.memBuffer != nil {
		return ls.memBuffer.GetLogEntries(from, limit)
	}

	return nil
}
