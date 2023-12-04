package console

import (
	"fmt"
	"strings"
)

type Option struct {
	Name        string
	Description string
	Default     any
	Type        any

	value any
	// modes: array, none/bool, required (value), optional (may have value or no), negatable (no-*)
}
