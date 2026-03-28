package research

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ResearchStorage) CheckResearchInProgress(ctx context.Context, user_id uuid.UUID) (bool, error) {
	const getResearchProgressCountQuery = `
		SELECT EXISTS (
			SELECT 1
			FROM session_beta.event_researches p
			WHERE user_id = $1
		);
	`

	var exists bool
	err := s.DB.QueryRow(ctx, getResearchProgressCountQuery, user_id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("DB.QueryRow(): %w", err)
	}

	return exists, nil
}
