package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetUserResources(ctx context.Context, userID uuid.UUID) (models.UserResources, error) {
	const getResourcesQuery = `
		SELECT 
			matter,
			doreye
		FROM session_beta.user_resources
		WHERE user_id = $1;
	`

	var (
		matter *uint64
		doreye *uint64
	)
	err := r.DB.QueryRow(ctx, getResourcesQuery, userID).Scan(
		&matter,
		&doreye,
	)
	if err != nil {
		return models.UserResources{}, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	var res models.UserResources
	if matter != nil {
		res.Matter = *matter
	}

	if doreye != nil {
		res.Doreye = *doreye
	}

	return res, nil
}
