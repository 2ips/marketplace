package http

import (
	"context"
	"fmt"

	"github.com/2ips/marketplace/cart/internal/app/usecase"
	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/2ips/marketplace/cart/pkg/api"
)

var _ api.StrictServerInterface = (*CartServer)(nil)

type CartServer struct {
	addItem *usecase.AddItem
	getCart *usecase.GetCart
}

func NewCartServer(addItem *usecase.AddItem, getCart *usecase.GetCart) *CartServer {
	return &CartServer{addItem: addItem, getCart: getCart}
}

func (s *CartServer) AddItem(ctx context.Context, req api.AddItemRequestObject) (api.AddItemResponseObject, error) {
	if req.UserId < 1 {
		return api.AddItem400JSONResponse{
			Message: fmt.Sprintf("user ID must be positive, got %d", req.UserId)}, nil
	}

	if req.Sku < 1 {
		return api.AddItem400JSONResponse{
			Message: fmt.Sprintf("SKU must be positive, got %d", req.Sku)}, nil
	}

	userID, err := domain.NewUserID(uint64(req.UserId))
	if err != nil {
		return toAddItemErrorResponse(err), nil
	}

	sku, err := domain.NewSKU(uint64(req.Sku))
	if err != nil {
		return toAddItemErrorResponse(err), nil
	}

	count, err := domain.NewCount(req.Body.Count)
	if err != nil {
		return toAddItemErrorResponse(err), nil
	}

	if _, err := s.addItem.Execute(ctx, usecase.AddItemRequest{UserID: userID, SKU: sku, Count: count}); err != nil {
		return toAddItemErrorResponse(err), nil
	}

	return api.AddItem200Response{}, nil
}

func (s *CartServer) GetCart(ctx context.Context, req api.GetCartRequestObject) (api.GetCartResponseObject, error) {
	if req.UserId < 1 {
		return api.GetCart400JSONResponse{
			Message: fmt.Sprintf("user ID must be positive, got %d", req.UserId)}, nil
	}

	userID, err := domain.NewUserID(uint64(req.UserId))
	if err != nil {
		return toGetCartErrorResponse(err), nil
	}

	res, err := s.getCart.Execute(ctx, usecase.GetCartRequest{UserID: userID})
	if err != nil {
		return toGetCartErrorResponse(err), nil
	}

	items := res.Cart.Items()
	if len(items) == 0 {
		return api.GetCart404JSONResponse{
			Message: fmt.Sprintf("cart must be non-empty, user ID %d", userID.Value()),
		}, nil
	}

	models := make([]api.CartItem, len(items))
	for i, item := range items {
		model, err := toCartItem(item)
		if err != nil {
			return toGetCartErrorResponse(err), nil
		}

		models[i] = model
	}

	totalPrice, err := res.Cart.TotalPrice()
	if err != nil {
		return toGetCartErrorResponse(err), nil
	}

	return api.GetCart200JSONResponse{Items: models, TotalPrice: totalPrice.Value()}, nil
}
