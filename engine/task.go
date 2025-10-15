package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mjc-gh/pisces/internal/browser"
	"github.com/rs/zerolog"
)

type Task struct {
	action    string
	url       string
	userAgent string
	winHeight int
	winWidth  int
	received  time.Time
}

func NewTask(action, input string) Task {
	url := input

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("http://%s", url)
	}

	return Task{
		action:   action,
		url:      url,
		received: time.Now(),
	}
}

func (t *Task) SetDevice(deviceType, deviceSize string) {
	t.winHeight, t.winWidth = browser.DimensionsFromDeviceProfile(deviceType, deviceSize)
}

func (t *Task) SetUserAgent(deviceType, userAgentAlias string) {
	t.userAgent = browser.UserAgent(deviceType, userAgentAlias)
}

type Result struct {
	Action  string        `json:"action"`
	Elapsed time.Duration `json:"elapsed"`
	Error   error         `json:"error"`
	URL     string        `json:"url"`
	Result  Payload       `json:"result"`
}

func (r *Result) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "\t")
}

type Payload interface{}

func performTask(ctx context.Context, task *Task, logger *zerolog.Logger) Result {
	logger.Debug().Msgf("perform task: %+v", task)
	defer logger.Debug().Msgf("performed task: %+v", task)

	result := Result{
		Action: task.action,
		URL:    task.url,
	}

	switch task.action {
	case "analyze":
		analysis, err := performAnalyzeTask(ctx, task, logger)
		if err != nil {
			return errorResult(task, err)
		}

		result.Elapsed = time.Since(task.received)
		result.Result = &analysis

	default:
		return errorResult(task, fmt.Errorf("unknown action: %s", task.action))
	}

	return result
}

func errorResult(task *Task, err error) Result {
	return Result{
		task.action,
		time.Since(task.received),
		err,
		task.url,
		nil,
	}
}
