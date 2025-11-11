package engine

import (
	"context"
	"errors"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog"
)

type Input struct {
	Name  string
	Type  string
	Value string
}

type Form struct {
	Action string `json:"action"`
	Method string `json:"method"`
	Class  string `json:"class"`
	ID     string `json:"id"`
	Inputs []Input
}

type AnalyzeResult struct {
	Forms []Form `json:"forms"`
	*Visit
}

func performAnalyzeTask(ctx context.Context, task *Task, logger *zerolog.Logger) (AnalyzeResult, error) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	crawler := NewCrawler(task.userAgent, int64(task.winWidth), int64(task.winHeight))
	err := crawler.Visit(ctx, task.url, logger)
	if err != nil {
		return AnalyzeResult{}, err
	}

	visit := crawler.LastVisit()
	if visit == nil {
		return AnalyzeResult{}, errors.New("no visit from crawler")
	}

	result := AnalyzeResult{}

	if err = runFormAnalysis(ctx, &result); err != nil {
		logger.Warn().Msgf("analyze task form analysis error: %-v", err)
	}

	// TODO assign Visit

	return result, nil
}

func runFormAnalysis(ctx context.Context, result *AnalyzeResult) error {
	return chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var formNodes []*cdp.Node

		if err := chromedp.Nodes("form", &formNodes, chromedp.ByQueryAll).Do(ctx); err != nil {
			return err
		}

		formAttrs, err := attributesFromNodes(ctx, formNodes, []string{"action", "method", "class", "id"})
		if err != nil {
			return err
		}

		for _, attributes := range formAttrs {
			form := Form{
				Action: attributes[0],
				Method: attributes[1],
				Class:  attributes[2],
				ID:     attributes[3],
			}

			result.Forms = append(result.Forms, form)
		}

		return nil
	}))
}

/*
	var nodes []*chromedp.Node
	if err := chromedp.Nodes("a[href]", &nodes).Do(ctx); err != nil {
		return err
	}
	for _, node := range nodes {
		if href, ok := node.Attributes["href"]; ok {
			hrefs = append(hrefs, href)
		}
	}
	return nil
*/

func attributesFromNodes(ctx context.Context, nodes []*cdp.Node, attributes []string) ([][]string, error) {
	values := make([][]string, len(nodes))

	for i, node := range nodes {
		values[i] = make([]string, len(attributes))

		for j, attribute := range attributes {
			err := chromedp.Run(ctx,
				chromedp.JavascriptAttribute(
					node.FullXPath(), attribute, &values[i][j], chromedp.BySearch,
				),
			)

			// Fallback to Attribute method on the Node type
			if err != nil {
				if val, ok := node.Attribute(attribute); ok {
					values[i][j] = val
				}
			}
		}
	}

	return values, nil
}
