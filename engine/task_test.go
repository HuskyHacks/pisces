package engine

import (
	"context"
	"testing"

	"github.com/mjc-gh/pisces"
	"github.com/stretchr/testify/assert"
)

func TestPerformTaskUnknownType(t *testing.T) {
	task := Task{action: "huh", url: "http://example.com"}
	r := performTask(context.TODO(), &task, pisces.Logger())

	assert.Equal(t, "huh", r.Action)
	assert.Error(t, r.Error)
	assert.NotEmpty(t, r.URL)
	assert.NotEmpty(t, r.Elapsed)
}
