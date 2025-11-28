package engine

import (
	"context"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog"
)

type ScreenshotResult struct {
	Buffer *[]byte
}

func performScreenshotTask(ctx context.Context, task *Task, logger *zerolog.Logger) (ScreenshotResult, error) {
	var buf []byte

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.EmulateViewport(int64(task.winWidth), int64(task.winHeight)),
		emulation.SetUserAgentOverride(task.userAgent),
		chromedp.Navigate(task.url),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		return ScreenshotResult{}, err
	}

	return ScreenshotResult{&buf}, nil
}
