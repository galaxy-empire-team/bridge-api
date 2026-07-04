package planet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
)

func (r *PlanetStorage) GetUserNPCAttacks(ctx context.Context, userID uuid.UUID) ([]models.NPCAttack, error) {
	const getUserNPCAttacksQuery = `
		SELECT 
			npc_coordinate_z,
			attacked_at
		FROM session_beta.user_npc_attacks
		WHERE user_id = $1;
	`

	rows, err := r.DB.Query(ctx, getUserNPCAttacksQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	var npcAttacks []models.NPCAttack
	for rows.Next() {
		var npcAttack models.NPCAttack
		err = rows.Scan(
			&npcAttack.Z,
			&npcAttack.AttackedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		npcAttacks = append(npcAttacks, npcAttack)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return npcAttacks, nil
}
