package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *PlanetStorage) SetHasMoon(ctx context.Context, planetID uuid.UUID, hasMoon bool) error {
	const setHasMoonQuery = `
		UPDATE session_beta.planets
		SET has_moon = $2
		WHERE id = $1;
	`

	_, err := r.DB.Exec(ctx, setHasMoonQuery, planetID, hasMoon)
	if err != nil {
		return fmt.Errorf("DB.Exec(): %w", err)
	}

	return nil
}
