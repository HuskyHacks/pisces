package engine

type ClipboardCapture struct {
	set map[string]struct{}
}

func NewClipboardCapture() *ClipboardCapture {
	return &ClipboardCapture{
		make(map[string]struct{}),
	}
}

func (cc *ClipboardCapture) AddTo(value string) {
	if _, ok := cc.set[value]; !ok {
		cc.set[value] = struct{}{}
	}
}

func (cc *ClipboardCapture) Values() []string {
	values := make([]string, len(cc.set))
	idx := 0

	for v := range cc.set {
		values[idx] = v
		idx++
	}

	return values
}
