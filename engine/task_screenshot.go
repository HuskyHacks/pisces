package engine

import (
	"context"

	"github.com/rs/zerolog"
)

type ScreenshotResult struct{}

func performScreenshotTask(ctx context.Context, task *Task, logger *zerolog.Logger) (ScreenshotResult, error) {
	return ScreenshotResult{}, nil
}
