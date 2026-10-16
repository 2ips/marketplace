package domain

import "fmt"

type UserID struct {
	value uint64
}

type ErrInvalidUserID struct {
	Value uint64
}

func (e *ErrInvalidUserID) Error() string {
	return fmt.Sprintf("user ID must be non zero, got %d", e.Value)
}

func NewUserID(userID uint64) (UserID, error) {
	if userID == 0 {
		return UserID{}, &ErrInvalidUserID{Value: userID}
	}

	return UserID{value: userID}, nil
}

func (u UserID) Value() uint64 {
	return u.value
}
