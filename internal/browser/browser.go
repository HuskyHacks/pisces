package browser

import (
	"context"

	"github.com/chromedp/chromedp"
)

func StartLocal(ctx context.Context) (context.Context, context.CancelFunc) {
	return chromedp.NewExecAllocator(ctx, chromedp.Headless)
}

func StartRemote(ctx context.Context, url string) (context.Context, context.CancelFunc) {
	return chromedp.NewRemoteAllocator(ctx, url)
}
