package eval

import (
	"dabble/object"
	"fmt"
	"strings"
)

var NilFrame *Frame = nil

type Frame struct {
	caller Function
	symbol object.Symbol
	value  object.Value
	next   *Frame
}

func (f *Frame) Bind(symbol object.Symbol, value object.Value) *Frame {
	return &Frame{
		symbol: symbol,
		value:  value,
		next:   f,
	}
}

func (f *Frame) Resolve(symbol object.Symbol) object.Value {
	if f == nil {
		return object.Error(fmt.Sprintf("symbol not bound: %q", symbol))
	}
	if f.symbol == symbol {
		return f.value
	}
	return f.next.Resolve(symbol)
}

func (f *Frame) Call(caller Function) *Frame {
	return &Frame{
		caller: caller,
		next:   f,
	}
}

func (f *Frame) LastCaller() Function {
	if f == nil {
		return func(_ *Frame, _ ...object.Value) object.Value {
			return object.Error(fmt.Sprintf("no caller"))
		}
	}
	if f.caller != nil {

		return f.caller
	}
	return f.next.LastCaller()
}

func (f *Frame) BindAll(f2 *Frame) *Frame {
	for f2 != nil {
		if f2.caller != nil {
			f = f.Call(f2.caller)
		} else {
			f = f.Bind(f2.symbol, f2.value)
		}
		f2 = f2.next
	}
	return f
}

func (f *Frame) String() string {
	var rest bool
	var sb strings.Builder
	sb.WriteString("(")
	for f != nil {
		if f.caller != nil {
			f = f.next
			continue
		}
		if rest {
			sb.WriteString(" ")
		}
		sb.WriteString("(")
		sb.WriteString(f.symbol.String())
		sb.WriteString(" ")
		sb.WriteString(f.value.String())
		sb.WriteString(")")
		f = f.next
		rest = true
	}
	sb.WriteString(")")
	return sb.String()
}
