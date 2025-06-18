package logscope

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// logMemBuffer implements a circular buffer for storing log entries in memory.
// It acts as a logrus hook that captures entries and stores them in a fixed-size buffer
// with automatic rotation when the buffer is full. Thread-safe with mutex protection.
type logMemBuffer struct {
	logger     *LogScope       // Reference to the parent LogScope
	bufferSize uint64          // Maximum number of entries to store
	bufIdx     uint64          // Current buffer write index (total writes modulo bufferSize)
	bufMtx     sync.Mutex      // Protects concurrent access to buffer operations
	buf        []*logrus.Entry // Circular buffer storing the actual log entries
	lastIdx    uint64          // Total number of entries processed (never decreases)
}

// newLogMemBuffer creates a new memory buffer for log entries with the specified size.
// The buffer starts empty and grows up to bufferSize entries before rotating.
func newLogMemBuffer(logger *LogScope, bufferSize uint64) *logMemBuffer {
	return &logMemBuffer{
		logger:     logger,
		bufferSize: bufferSize,
		buf:        make([]*logrus.Entry, 0, bufferSize),
	}
}

// Levels returns all logrus levels, indicating this hook captures all log levels.
// This is required by the logrus.Hook interface for the memory buffer.
func (lmb *logMemBuffer) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire stores a log entry in the circular memory buffer.
// When the buffer is full, new entries overwrite the oldest entries.
// Each entry is assigned a sequential index for tracking total throughput.
func (lmb *logMemBuffer) Fire(entry *logrus.Entry) error {
	lmb.bufMtx.Lock()
	defer lmb.bufMtx.Unlock()

	logIdx := lmb.lastIdx + 1
	lmb.lastIdx = logIdx

	if lmb.bufIdx >= lmb.bufferSize {
		bufIdx := lmb.bufIdx % lmb.bufferSize
		lmb.buf[bufIdx] = entry
	} else {
		lmb.buf = append(lmb.buf, entry)
	}

	lmb.bufIdx++

	return nil
}

// GetLogEntryCount returns the total number of log entries processed by this buffer.
// This count includes entries that may have been overwritten due to buffer rotation.
func (lmb *logMemBuffer) GetLogEntryCount() int {
	return int(lmb.lastIdx)
}

// GetLogEntries retrieves log entries from the circular buffer with filtering options.
// The from parameter excludes entries at or before the specified time (0 time = no filter).
// The limit parameter restricts returned entries (0 = no limit).
// Returns entries in chronological order, handling circular buffer reconstruction.
func (lmb *logMemBuffer) GetLogEntries(from time.Time, limit int) []*logrus.Entry {
	lmb.bufMtx.Lock()
	defer lmb.bufMtx.Unlock()

	var entries []*logrus.Entry

	if lmb.bufIdx >= lmb.bufferSize {
		entries = make([]*logrus.Entry, lmb.bufferSize)
		firstIdx := lmb.bufIdx % lmb.bufferSize

		copy(entries, lmb.buf[firstIdx:])
		copy(entries[lmb.bufferSize-firstIdx:], lmb.buf[0:firstIdx])
	} else {
		entries = make([]*logrus.Entry, lmb.bufIdx)
		copy(entries, lmb.buf)
	}

	if len(entries) == 0 {
		return entries
	}

	for !from.IsZero() && len(entries) > 0 && (entries[0].Time.Before(from) || entries[0].Time.Equal(from)) {
		entries = entries[1:]
	}

	if limit != 0 && len(entries) > limit {
		if from.IsZero() {
			offset := len(entries) - limit
			entries = entries[offset:]
		} else {
			entries = entries[0:limit]
		}
	}

	return entries
}
