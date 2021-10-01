package eval

import (
	"fmt"
	"strings"
)

var t *trace

type trace struct {
	lines  []string
	indent int
}

func BeginTrace() {
	t = &trace{}
}

func EndTrace() string {
	out := t.String()
	t = nil
	return out
}

func T(msg string, args ...interface{}) string {
	return t.T(msg, args...)
}

func (t *trace) T(msg string, args ...interface{}) string {
	if t == nil {
		return msg
	}
	indent := strings.Repeat(".   ", t.indent)
	out := fmt.Sprintf(indent+msg, args...)
	t.lines = append(t.lines, out)
	return out
}

func (t *trace) String() string {
	if t == nil {
		return "<nil>"
	}
	return strings.Join(t.lines, "\n")
}

func (t *trace) In() {
	if t != nil {
		t.indent += 1
	}
}

func (t *trace) Out() {
	if t != nil {
		t.indent -= 1
	}
}
