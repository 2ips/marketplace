package port

import (
	"context"
	"errors"
	"fmt"

	"github.com/2ips/marketplace/cart/internal/domain"
)

var ErrStorage = errors.New("storage error")

type ErrUserNotFound struct {
	Value domain.UserID
}

type CartRepository interface {
	Get(ctx context.Context, userID domain.UserID) (domain.Cart, error)
	AddItem(ctx context.Context, user domain.UserID, item domain.Item) error
}

func (e *ErrUserNotFound) Error() string {
	return fmt.Sprintf("cannot find user with ID %d", e.Value)
}
