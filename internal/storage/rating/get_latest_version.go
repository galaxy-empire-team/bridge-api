package rating

import (
	"context"
	"fmt"
)

func (r *RatingStorage) GetLatestVersion(ctx context.Context, versionType string) (uint32, error) {
	const getLatestVersionQuery = `
		SELECT MAX(version)
		FROM session_beta.r_versions
		WHERE status = 'success' AND type = $1;
	`

	var version *uint32
	err := r.DB.QueryRow(ctx, getLatestVersionQuery, versionType).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("DB.QueryRow.Scan(): %w", err)
	}

	if version == nil {
		return 0, nil
	}

	return *version, nil
}
