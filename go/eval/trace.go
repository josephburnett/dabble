package eval

import (
	"dabble/object"
	"fmt"
	"runtime"
	"strings"
)

type Trace struct {
	lines []string
}

func NewTrace() *Trace {
	return &Trace{}
}

func (t *Trace) T(msg string, args ...interface{}) string {
	if t == nil {
		return msg
	}
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	var loc string
	if ok && details != nil {
		loc = details.Name()
	}
	out := fmt.Sprintf("["+loc+"] "+msg, args...)
	t.lines = append(t.lines, out)
	return out
}

func (t *Trace) Eval(env *object.Binding, value object.Value) object.Value {
	return eval(env, false, t, value)
}

func (t *Trace) String() string {
	if t == nil {
		return "<nil>"
	}
	return strings.Join(t.lines, "\n")
}
