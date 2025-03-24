package spamoor

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ScenarioRunner struct {
	ctx    context.Context
	logger *logrus.Logger
}

func NewScenarioRunner(ctx context.Context, logger *logrus.Logger) *ScenarioRunner {
	return &ScenarioRunner{
		ctx:    ctx,
		logger: logger,
	}
}
