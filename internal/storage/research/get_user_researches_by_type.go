package research

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *ResearchStorage) GetUserResearchesByTypes(ctx context.Context, userID uuid.UUID, researchTypes []consts.ResearchType) (map[consts.ResearchType]consts.ResearchID, error) {
	const getUserResearchesByTypeQuery = `
		SELECT 
			r.id,
			r.research_type
		FROM session_beta.user_researches ur
		JOIN session_beta.s_researches r ON r.id = ur.research_id
		WHERE ur.user_id = $1 AND r.research_type = ANY($2);
	`

	rows, err := r.DB.Query(ctx, getUserResearchesByTypeQuery, userID, researchTypes)
	if err != nil {
		return nil, fmt.Errorf("DB.Query.Scan(): %w", err)
	}
	defer rows.Close()

	researchesByType := make(map[consts.ResearchType]consts.ResearchID)
	for rows.Next() {
		var researchID consts.ResearchID
		var researchType consts.ResearchType
		err = rows.Scan(&researchID, &researchType)
		if err != nil {
			return nil, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
		}

		researchesByType[researchType] = researchID
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return researchesByType, nil
}
