package browser

import (
	"context"

	"github.com/chromedp/chromedp"
)

const (
	CHROME_DP_HEADLESS_SHELL = "http://127.0.0.1:9222/json/version"
)

func StartLocal(ctx context.Context) (context.Context, error) {
	allocatorContext, _ := chromedp.NewExecAllocator(ctx)
	c, _ := chromedp.NewContext(allocatorContext)

	if err := chromedp.Run(c); err != nil {
		return nil, err
	}

	return c, nil
}

func StartRemote(ctx context.Context, url string) (context.Context, error) {
	allocatorContext, _ := chromedp.NewRemoteAllocator(ctx, url)
	c, _ := chromedp.NewContext(allocatorContext)

	if err := chromedp.Run(c); err != nil {
		return nil, err
	}

	return c, nil
}
