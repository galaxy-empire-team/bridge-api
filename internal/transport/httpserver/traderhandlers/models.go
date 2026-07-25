package traderhandlers

import (
	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type ErrorResponse struct {
	Err string `json:"err"`
}

type BuyItemRequest struct {
	PlanetID uuid.UUID           `json:"planetID"`
	ItemID   consts.TraderItemID `json:"itemID"`
}

type BuyItemResponse struct {
	Message string `json:"message"`
}
