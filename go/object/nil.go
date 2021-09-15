package object

type nilT struct{}

var Nil nilT = struct{}{}

func (n nilT) isValue() {}

func (n nilT) Inspect() string {
	return "()"
}
