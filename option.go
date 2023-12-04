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

type array []string

func (arr *array) String() string {
	return fmt.Sprintf("[%s]", strings.Join(*arr, " "))
}

func (arr *array) Set(value string) error {
	*arr = append(*arr, value)

	return nil
}
