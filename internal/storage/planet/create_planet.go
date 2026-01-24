package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"initialservice/internal/models"
)

func (r *PlanetStorage) CreatePlanet(ctx context.Context, userID uuid.UUID, planet models.Planet) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pool.Begin(): %w", err)
	}

	defer func() {
		rollbackErr := tx.Rollback(ctx)

		if rollbackErr == nil || errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return
		}

		if err != nil {
			err = fmt.Errorf("%w; tx.Rollback(): %w", err, rollbackErr)

			return
		}

		err = fmt.Errorf("tx.Rollback(): %w", rollbackErr)
	}()

	err = r.createPlanetRow(ctx, tx, userID, planet)
	if err != nil {
		return fmt.Errorf("createPlanet(): %w", err)
	}

	err = r.createResourcesRow(ctx, tx, planet.ID)
	if err != nil {
		return fmt.Errorf("createResources(): %w", err)
	}

	err = r.createBuildingsRows(ctx, tx, planet)
	if err != nil {
		return fmt.Errorf("createBuildings(): %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit(): %w", err)
	}

	return nil
}

func (r *PlanetStorage) createPlanetRow(ctx context.Context, tx pgx.Tx, userID uuid.UUID, planet models.Planet) error {
	const createPlanetQuery = `
		INSERT INTO session_beta.planets (
			id,
			user_id,
			x,
			y,
			z,
			has_moon,
			is_capitol,
			colonized_at
		) VALUES (
			$1,   --- planet.ID
			$2,   --- userID
			$3,	  --- planet.X
			$4,   --- planet.Y
			$5,   --- planet.Z
			$6,   --- planet.HasMoon
			$7,   --- planet.IsCapitol
			NOW() --- colonized_at
		);
	`

	_, err := tx.Exec(
		ctx,
		createPlanetQuery,
		planet.ID,
		userID,
		planet.Location.X,
		planet.Location.Y,
		planet.Location.Z,
		planet.HasMoon,
		planet.IsCapitol,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				if pgErr.ConstraintName == "planet_have_unique_x_y_z" {
					return fmt.Errorf("DB.Pool.Exec(): %w", models.ErrPlanetCoordinatesAlreadyTaken)
				}

				return fmt.Errorf("DB.Pool.Exec(): %w", models.ErrCapitolAlreadyExists)
			}
		}

		return fmt.Errorf("DB.Pool.Exec(): %w", err)
	}

	return nil
}

func (r *PlanetStorage) createResourcesRow(ctx context.Context, tx pgx.Tx, planetID uuid.UUID) error {
	const createResourcesQuery = `
		INSERT INTO session_beta.planet_resources (
			planet_id,
			metal,
			crystal,
			gas,
			updated_at
		) VALUES (
			$1,    --- planet.ID
			0,     --- metal
			0,     --- crystal
			0,     --- gas
			NOW()  --- updated_at
		);
	`

	_, err := tx.Exec(
		ctx,
		createResourcesQuery,
		planetID,
	)
	if err != nil {
		return fmt.Errorf("DB.Pool.Exec(): %w", err)
	}

	return nil
}

func (r *PlanetStorage) createBuildingsRows(ctx context.Context, tx pgx.Tx, planet models.Planet) error {
	const createBuildingQuery = `
		INSERT INTO session_beta.planet_buildings (
			planet_id,
			building_id,
			updated_at,
			finished_at
		) VALUES (
			$1,    --- planet.ID
			( SELECT id FROM session_beta.buildings WHERE type = $2 AND level = 0 ), --- building_id
			NOW(),  --- updated_at
			NULL	--- finished_at
		);
	`

	batch := &pgx.Batch{}
	for buildingType := range planet.Buildings {
		batch.Queue(
			createBuildingQuery,
			planet.ID,
			buildingType,
		)
	}

	batchResults := tx.SendBatch(ctx, batch)
	defer batchResults.Close()

	for i := 0; i < batch.Len(); i++ {
		_, err := batchResults.Exec()
		if err != nil {
			return fmt.Errorf("DB.Pool.Exec(): %w", err)
		}
	}

	if err := batchResults.Close(); err != nil {
		return fmt.Errorf("batchResults.Close(): %w", err)
	}

	return nil
}
