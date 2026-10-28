package http

import (
	"errors"
	"net/http"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/2ips/marketplace/cart/pkg/api"
)

func toAddItemErrorResponse(err error) api.AddItemResponseObject {
	var (
		invalidUserID   *domain.ErrInvalidUserID
		invalidSKU      *domain.ErrInvalidSKU
		invalidCount    *domain.ErrInvalidCount
		countOverflow   *domain.ErrCountOverflow
		productNotFound *port.ErrProductNotFound
	)

	switch {
	case errors.As(err, &invalidUserID),
		errors.As(err, &invalidSKU),
		errors.As(err, &invalidCount),
		errors.As(err, &countOverflow):
		return api.AddItem400JSONResponse{Message: err.Error()}
	case errors.As(err, &productNotFound):
		return api.AddItem412JSONResponse{Message: productNotFound.Error()}
	default:
		return api.AddItemdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       api.Error{Message: "internal error"},
		}
	}
}

func toGetCartErrorResponse(err error) api.GetCartResponseObject {
	var (
		invalidUserID   *domain.ErrInvalidUserID
		userNotFound    *port.ErrUserNotFound
		productNotFound *port.ErrProductNotFound
	)

	switch {
	case errors.As(err, &invalidUserID):
		return api.GetCart400JSONResponse{Message: err.Error()}
	case errors.As(err, &userNotFound):
		return api.GetCart404JSONResponse{Message: userNotFound.Error()}
	case errors.As(err, &productNotFound):
		return api.GetCart412JSONResponse{Message: productNotFound.Error()}
	default:
		return api.GetCartdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       api.Error{Message: "internal error"},
		}
	}
}
