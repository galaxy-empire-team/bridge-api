package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *PlanetStorage) GetUserNPCAttackByZ(ctx context.Context, userID uuid.UUID, z consts.PlanetPositionZ) (*models.NPCAttack, error) {
	const getUserNPCAttacksQuery = `
		SELECT
			attacked_at
		FROM session_beta.user_npc_attacks
		WHERE user_id = $1 AND npc_coordinate_z = $2;
	`

	fmt.Println(userID)
	npcAttack := &models.NPCAttack{
		Z: z,
	}
	err := r.DB.QueryRow(ctx, getUserNPCAttacksQuery, userID, z).Scan(
		&npcAttack.AttackedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return npcAttack, nil
}
