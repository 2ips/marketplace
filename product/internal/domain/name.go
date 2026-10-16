package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const MaxNameLength = 255

type Name struct {
	value string
}

type ErrNameEmpty struct {
	Value string
}

type ErrNameTooLong struct {
	Length int
}

func (e *ErrNameEmpty) Error() string {
	return fmt.Sprintf("name must be non-empty, got %s", e.Value)
}

func (e *ErrNameTooLong) Error() string {
	return fmt.Sprintf("name must be at most %d characters, got %d", MaxNameLength, e.Length)
}

func NewName(name string) (Name, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return Name{}, &ErrNameEmpty{Value: name}
	}

	length := utf8.RuneCountInString(name)
	if length > MaxNameLength {
		return Name{}, &ErrNameTooLong{Length: length}
	}

	return Name{value: name}, nil
}

func (n Name) Value() string {
	return n.value
}
