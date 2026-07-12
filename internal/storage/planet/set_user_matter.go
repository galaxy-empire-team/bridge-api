package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *PlanetStorage) SetUserMatter(ctx context.Context, userID uuid.UUID, matter uint64) error {
	const setMatterQuery = `
		UPDATE session_beta.user_resources
		SET 
			matter = $2,
			updated_at = now()
		WHERE user_id = $1;
	`

	_, err := r.DB.Exec(ctx, setMatterQuery, userID, matter)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
