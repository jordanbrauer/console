package console

const (
	ArgumentOptional Bits = 1 << iota
)

type Argument struct {
	Name        string
	Description string

	// modes: optional, array
	Mode Bits
}

func (argument *Argument) Optional() bool {
	return argument.Mode.Has(ArgumentOptional)
}
