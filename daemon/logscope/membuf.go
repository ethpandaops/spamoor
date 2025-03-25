package logscope

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type logMemBuffer struct {
	logger     *LogScope
	bufferSize uint64
	bufIdx     uint64
	bufMtx     sync.Mutex
	buf        []*logrus.Entry
	lastIdx    uint64
}

func newLogMemBuffer(logger *LogScope, bufferSize uint64) *logMemBuffer {
	return &logMemBuffer{
		logger:     logger,
		bufferSize: bufferSize,
		buf:        make([]*logrus.Entry, 0, bufferSize),
	}
}

func (lmb *logMemBuffer) Levels() []logrus.Level {
	return logrus.AllLevels
}

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

func (lmb *logMemBuffer) GetLogEntryCount() int {
	return int(lmb.lastIdx)
}

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

	for !from.IsZero() && len(entries) > 0 && entries[0].Time.Before(from) {
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
