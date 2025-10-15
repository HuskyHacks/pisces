package engine

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog"
)

type AnalyzeResult struct {
	Assets map[string]*Asset
	//headers http.Header
	//cookies map[string]string
}

type Asset struct {
	URL             string
	RequestHeaders  map[string]any
	ResponseHeaders map[string]any
	Body            []byte
}

func performAnalyzeTask(ctx context.Context, task *Task, logger *zerolog.Logger) (AnalyzeResult, error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	result := AnalyzeResult{}
	result.Assets = make(map[string]*Asset)

	chromedp.ListenTarget(ctx, func(ev any) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			result.Assets[string(ev.RequestID)] = &Asset{
				URL:            ev.Request.URL,
				RequestHeaders: ev.Request.Headers,
			}
		case *network.EventResponseReceived:
			if asset, ok := result.Assets[string(ev.RequestID)]; ok {
				asset.ResponseHeaders = ev.Response.Headers
			}
		case *network.EventLoadingFinished:
			if asset, ok := result.Assets[string(ev.RequestID)]; ok {
				go func(reqID network.RequestID) {
					var body []byte

					// ActionFunc to bind body and handle error
					fn := func(ctx context.Context) (err error) {
						body, err = network.GetResponseBody(reqID).Do(ctx)
						return err
					}

					err := chromedp.Run(ctx, chromedp.ActionFunc(fn))
					if err != nil {
						logger.Warn().Msgf("analyze task: event finished error: %s", err)
						return
					}

					asset.Body = body
				}(ev.RequestID)
			}
		}
	})

	initialSteps := []chromedp.Action{
		network.Enable(),
		chromedp.EmulateViewport(int64(task.winWidth), int64(task.winHeight)),
		emulation.SetUserAgentOverride(task.userAgent),
		chromedp.Navigate(task.url),
		chromedp.Sleep(2 * time.Second),
		//chromedp.ActionFunc(func(ctx context.Context) error {
		//c, err := network.GetCookies().Do(ctx)
		//if err != nil {
		//return err
		//}
		//cookies = c
		//return nil
		//}),
	}

	if err := chromedp.Run(ctx, initialSteps...); err != nil {
		return AnalyzeResult{}, err
	}

	return result, nil
}
