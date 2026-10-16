package memory

import (
	"fmt"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/domain"
)

type cartModel struct {
	Items []itemModel
}

func cartToModel(cart *domain.Cart) cartModel {
	items := cart.Items()

	models := make([]itemModel, len(items))
	for i, item := range items {
		models[i] = itemToModel(item)
	}

	return cartModel{Items: models}
}

func cartFromModel(userID domain.UserID, model cartModel) (*domain.Cart, error) {
	items := make([]*domain.Item, len(model.Items))
	for i, itemModel := range model.Items {
		item, err := itemFromModel(itemModel)
		if err != nil {
			return nil, err
		}

		items[i] = item
	}

	cart, err := domain.NewCart(userID, items)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	return cart, nil
}
