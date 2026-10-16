package domain

import "fmt"

type Count struct {
	value uint32
}

type ErrInvalidCount struct {
	Value uint32
}

func (e *ErrInvalidCount) Error() string {
	return fmt.Sprintf("count must be non zero, got %d", e.Value)
}

func NewCount(count uint32) (Count, error) {
	if count == 0 {
		return Count{}, &ErrInvalidCount{Value: count}
	}

	return Count{value: count}, nil
}

func (c Count) Value() uint32 {
	return c.value
}
