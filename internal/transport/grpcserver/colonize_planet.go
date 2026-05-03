package grpcserver

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/galaxy-empire-team/bridge-api/api/gen/go/planet/v1"
	"github.com/galaxy-empire-team/bridge-api/internal/models"
	planetservice "github.com/galaxy-empire-team/bridge-api/internal/service/planet"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func (s *Server) ColonizePlanet(ctx context.Context, req *planetpb.ColonizePlanetRequest) (*planetpb.ColonizePlanetResponse, error) {
	userID, err := uuid.Parse(req.GetUserID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	err = s.planetService.ColonizePlanet(ctx, userID, planetservice.CreatePlanetRequest{
		OperationID: req.GetOperationID(),
		Coordinates: models.Coordinates{
			X: consts.PlanetPositionX(req.GetCoordinates().GetX()),
			Y: consts.PlanetPositionY(req.GetCoordinates().GetY()),
			Z: consts.PlanetPositionZ(req.GetCoordinates().GetZ()),
		},
		Resources: models.Resources{
			Metal:   req.GetResources().GetMetal(),
			Crystal: req.GetResources().GetCrystal(),
			Gas:     req.GetResources().GetGas(),
		},
		IsCapitol: req.GetIsCapitol(),
	})
	if err != nil {
		if errors.Is(err, models.ErrPlanetCoordinatesAlreadyTaken) {
			return nil, status.Error(codes.AlreadyExists, "another user colonized the planet")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &planetpb.ColonizePlanetResponse{}, nil
}
