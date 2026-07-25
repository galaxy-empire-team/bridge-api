package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) SetUserResources(ctx context.Context, userID uuid.UUID, resources models.UserResources) error {
	const setUserResourcesQuery = `
		UPDATE session_beta.user_resources
		SET 
			matter = $2,
			doreye = $3,
			updated_at = now()
		WHERE user_id = $1;
	`

	_, err := r.DB.Exec(ctx, setUserResourcesQuery, userID, resources.Matter, resources.Doreye)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
