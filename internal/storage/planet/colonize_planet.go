package planet

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (r *PlanetStorage) ColonizePlanet(ctx context.Context, planet models.Planet, operationID uint64) error {
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

	if operationID != 0 {
		isAvailable, err := r.CheckIsCreationAvailable(ctx, operationID)
		if err != nil {
			return fmt.Errorf("CheckIsCreationAvailable(): %w", err)
		}

		if !isAvailable {
			return nil
		}
	}

	planetToColonize := fromPlanetModel(planet)

	err = r.createPlanetRow(ctx, tx, planetToColonize)
	if err != nil {
		return fmt.Errorf("createPlanet(): %w", err)
	}

	err = r.createResourcesRow(ctx, tx, planetToColonize)
	if err != nil {
		return fmt.Errorf("createResources(): %w", err)
	}

	err = r.createMineRows(ctx, tx, planetToColonize.ID)
	if err != nil {
		return fmt.Errorf("createMineRows(): %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("tx.Commit(): %w", err)
	}

	return nil
}

const planetCreatedStatus = "created"

func (r *PlanetStorage) CheckIsCreationAvailable(ctx context.Context, operationID uint64) (bool, error) {
	const tryToInsertNewOperationIDQuery = `
		INSERT INTO session_beta.em_colonizations (id, status)
		VALUES ($1, $2);
	`

	_, err := r.DB.Exec(ctx, tryToInsertNewOperationIDQuery, operationID, planetCreatedStatus)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolationCode {
				return false, nil
			}
		}

		return false, fmt.Errorf("DB.Pool.Exec(): %w", err)
	}

	return true, nil
}

func (r *PlanetStorage) createPlanetRow(ctx context.Context, tx pgx.Tx, planet planetToColonize) error {
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
		planet.UserID,
		planet.Coordinates.X,
		planet.Coordinates.Y,
		planet.Coordinates.Z,
		planet.HasMoon,
		planet.IsCapitol,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case uniqueViolationCode:
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

func (r *PlanetStorage) createResourcesRow(ctx context.Context, tx pgx.Tx, planet planetToColonize) error {
	const createResourcesQuery = `
		INSERT INTO session_beta.planet_resources (
			planet_id,
			metal,
			crystal,
			gas,
			updated_at
		) VALUES (
			$1,    --- planet.ID
			$2,     --- metal
			$3,     --- crystal
			$4,     --- gas
			NOW()  --- updated_at
		);
	`

	_, err := tx.Exec(
		ctx,
		createResourcesQuery,
		planet.ID,
		planet.Resources.Metal,
		planet.Resources.Crystal,
		planet.Resources.Gas,
	)
	if err != nil {
		return fmt.Errorf("DB.Pool.Exec(): %w", err)
	}

	return nil
}

func (r *PlanetStorage) createMineRows(ctx context.Context, tx pgx.Tx, planetID uuid.UUID) error {
	const createMinesQuery = `
		INSERT INTO session_beta.planet_buildings (
			planet_id,
			building_id
		)
		SELECT $1, id
		FROM session_beta.s_buildings sb
		WHERE sb.building_type = ANY($2) AND sb.level = 1;
	`

	_, err := tx.Exec(
		ctx,
		createMinesQuery,
		planetID,
		[]consts.BuildingType{consts.BuildingTypeMetalMine, consts.BuildingTypeCrystalMine, consts.BuildingTypeGasMine},
	)
	if err != nil {
		return fmt.Errorf("DB.Pool.Exec(): %w", err)
	}

	return nil
}
