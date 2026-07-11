package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Attack(ctx context.Context, mission models.MissionStart) (models.UserMission, error) {
	if err := s.validateAttackFleet(mission.Fleet, mission.Cargo); err != nil {
		return models.UserMission{}, fmt.Errorf("validateAttackFleet(): %w", err)
	}

	if err := s.repository.CheckPlanetOwner(ctx, mission.UserID, mission.PlanetFrom); err != nil {
		return models.UserMission{}, fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	allowed, err := s.checkAttackAllowed(ctx, mission.UserID, mission.PlanetTo)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("checkAttackAllowed(): %w", err)
	}
	if !allowed {
		return models.UserMission{}, fmt.Errorf("attack not allowed: %w", err)
	}

	planetExists, err := s.planetStorage.CheckPlanetExists(ctx, mission.PlanetTo)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("planetStorage.CheckPlanetExists(): %w", err)
	}
	if !planetExists && !s.isPlanetNPC(mission.PlanetTo.Z) {
		return models.UserMission{}, models.ErrPlanetNotFound
	}

	attackMission, err := s.prepareUserMission(ctx, mission, consts.MissionTypeAttack)
	if err != nil {
		return models.UserMission{}, fmt.Errorf("prepareUserMission(): %w", err)
	}

	return attackMission, s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err := s.removeFromPlanet(ctx, mission.UserID, mission.PlanetFrom, mission.Fleet, mission.Cargo, storages)
		if err != nil {
			return fmt.Errorf("removeFromPlanet(): %w", err)
		}

		attackEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			MissionID:   attackMission.MissionID,
			Fleet:       mission.Fleet,
			IsReturning: attackMission.IsReturning,
			StartedAt:   attackMission.StartedAt,
			FinishedAt:  attackMission.FinishedAt,
		}

		eventID, err := storages.CreateMissionEvent(ctx, attackEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

		attackMission.ID = eventID

		return nil
	})
}

func (s *Service) checkAttackAllowed(ctx context.Context, userID uuid.UUID, planet models.Coordinates) (bool, error) {
	if s.isPlanetNPC(planet.Z) {
		npcAttack, err := s.planetStorage.GetUserNPCAttackByZ(ctx, userID, planet.Z)
		if err != nil {
			return false, fmt.Errorf("planetStorage.GetUserNPCAttackByZ(): %w", err)
		}

		if npcAttack != nil && npcAttack.AttackedAt.After(time.Now().UTC().Add(-consts.NPCAttackCooldown)) {
			return false, models.ErrNPCCooldownNotExpired
		}

		exists, err := s.missionStorage.CheckNPCMissionExists(ctx, userID, planet.Z)
		if err != nil {
			return false, fmt.Errorf("missionStorage.CheckNPCMissionExists(): %w", err)
		}
		if exists {
			return false, models.ErrNPCMissionAlreadyExists
		}
	} else {
		attackedAt, err := s.planetStorage.GetPlanetAttackedAt(ctx, planet)
		if err != nil {
			return false, fmt.Errorf("planetStorage.GetPlanetAttackedAt(): %w", err)
		}

		if attackedAt != nil && attackedAt.After(time.Now().UTC().Add(-consts.PlanetAttackCooldown)) {
			return false, models.ErrPlanetAttackedCooldownNotExpired
		}
	}

	return true, nil
}

func (s *Service) validateAttackFleet(fleet []models.FleetUnitCount, cargo models.Resources) error {
	if !cargo.IsEmpty() {
		return models.ErrCargoIsNotEmpty
	}

	err := s.validateFleet(fleet, cargo)
	if err != nil {
		return fmt.Errorf("validateFleet(): %w", err)
	}

	return nil
}
