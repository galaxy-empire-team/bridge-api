package user

import (
	"context"
	"fmt"
)

func (r *UserStorage) GetUserCount(ctx context.Context) (uint64, error) {
	const getUsersCountQuery = `
		SELECT COUNT(*) FROM session_beta.users;
	`

	var count uint64
	err := r.DB.QueryRow(ctx, getUsersCountQuery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	return count, nil
}
