package erc20bloater

import (
	"io"

	"github.com/sirupsen/logrus"
)

func newTestLogger() *logrus.Entry {
	lg := logrus.New()
	lg.SetOutput(io.Discard)

	return lg.WithField("test", ScenarioName)
}
