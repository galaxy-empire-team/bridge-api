package traderhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

type TraderService interface {
	Buy(ctx context.Context, userID uuid.UUID, planetID uuid.UUID, itemID consts.TraderItemID) error
}
