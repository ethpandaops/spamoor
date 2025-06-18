package utils

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

func RecoverPanic(logger *logrus.Entry, routineName string) bool {
	if err := recover(); err != nil {
		logger.Errorf("uncaught panic in %s: %v, stack: %v", routineName, err, string(debug.Stack()))
		return true
	}
	return false
}
