package mission

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/galaxy-empire-team/bridge-api/internal/models"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Service) Attack(ctx context.Context, mission models.MissionStart) error {
	if !mission.Cargo.IsEmpty() {
		return models.ErrCargoIsNotEmpty
	}

	err := s.checkFleetValid(mission.Fleet)
	if err != nil {
		return fmt.Errorf("checkFleetValid(): %w", err)
	}

	planetExists, err := s.planetStorage.CheckPlanetExists(ctx, mission.PlanetTo)
	if err != nil {
		return fmt.Errorf("planetStorage.CheckPlanetExists(): %w", err)
	}
	if !planetExists && !s.isPlanetNPC(mission.PlanetTo.Z) {
		return models.ErrPlanetNotFound
	}

	allowed, err := s.checkAttackAllowed(ctx, mission.UserID, mission.PlanetTo)
	if err != nil {
		return fmt.Errorf("checkAttackAllowed(): %w", err)
	}
	if !allowed {
		return fmt.Errorf("attack not allowed: %w", err)
	}

	if err := s.repository.CheckPlanetOwner(ctx, mission.UserID, mission.PlanetFrom); err != nil {
		return fmt.Errorf("CheckPlanetOwner(): %w", err)
	}

	missionID, err := s.registry.GetMissionIDByType(consts.MissionTypeAttack)
	if err != nil {
		return fmt.Errorf("registry.GetMissionIDByType(): %w", err)
	}

	planetFromCoordinates, err := s.planetStorage.GetCoordinates(ctx, mission.PlanetFrom)
	if err != nil {
		return fmt.Errorf("planetStorage.GetCoordinates(): %w", err)
	}

	missionDuration, err := s.calculateMissionDuration(durationInput{
		PlanetFrom:      planetFromCoordinates,
		PlanetTo:        mission.PlanetTo,
		Fleet:           mission.Fleet,
		SpeedMultiplier: mission.SpeedMultiplier,
		IsSpyMission:    false,
	})
	if err != nil {
		return fmt.Errorf("calculateMissionDuration(): %w", err)
	}

	return s.txManager.ExecMissionTx(ctx, func(ctx context.Context, storages TxStorages) error {
		err = s.removeFleetFromPlanet(ctx, mission.PlanetFrom, mission.Fleet, storages)
		if err != nil {
			return fmt.Errorf("removeFleetFromPlanet(): %w", err)
		}

		startedAt := time.Now().UTC()
		attackEvent := models.MissionEvent{
			UserID:      mission.UserID,
			PlanetFrom:  mission.PlanetFrom,
			PlanetTo:    mission.PlanetTo,
			Type:        missionID,
			Fleet:       mission.Fleet,
			IsReturning: false,
			StartedAt:   startedAt,
			FinishedAt:  startedAt.Add(missionDuration),
		}

		err = storages.CreateMissionEvent(ctx, attackEvent)
		if err != nil {
			return fmt.Errorf("missionStorage.CreateMissionEvent(): %w", err)
		}

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
