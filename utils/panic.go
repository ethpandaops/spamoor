package utils

import (
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

func RecoverPanic(logger *logrus.Entry, routineName string, restartFn func()) bool {
	if err := recover(); err != nil {
		logger.Errorf("uncaught panic in %s: %v, stack: %v", routineName, err, string(debug.Stack()))

		if restartFn != nil {
			go func() {
				time.Sleep(5 * time.Second)
				restartFn()
			}()
		}

		return true
	}
	return false
}
