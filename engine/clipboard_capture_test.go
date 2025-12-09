package engine

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClipboardCapture_AddToAndValues(t *testing.T) {
	cc := NewClipboardCapture()

	cc.AddTo("first")
	cc.AddTo("second")
	cc.AddTo("first")
	cc.AddTo("third")

	values := cc.Values()
	sort.Strings(values)

	assert.Equal(t, []string{"first", "second", "third"}, values)
}
