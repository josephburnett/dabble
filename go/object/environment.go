package object

import (
	"fmt"
	"strings"
)

type Binding struct {
	Symbol Symbol
	Value  Value
	Next   *Binding
}

func (b *Binding) Resolve(symbol Symbol) Value {
	if b == nil {
		return Error(fmt.Sprintf("symbol not bound: %q", symbol))
	}
	if b.Symbol == symbol {
		return b.Value
	}
	return b.Next.Resolve(symbol)
}

func (b *Binding) String() string {
	var rest bool
	var sb strings.Builder
	sb.WriteString("(")
	for b != nil {
		if rest {
			sb.WriteString(" ")
		}
		sb.WriteString("(")
		sb.WriteString(b.Symbol.String())
		sb.WriteString(" ")
		sb.WriteString(b.Value.String())
		sb.WriteString(")")
		b = b.Next
		rest = true
	}
	sb.WriteString(")")
	return sb.String()
}
