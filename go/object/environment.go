package object

import (
	"fmt"
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
