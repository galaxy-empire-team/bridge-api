package grpcserver

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/galaxy-empire-team/bridge-api/api/gen/go/planet/v1"
)

func (s *Server) UpdatePlanetResources(ctx context.Context, req *planetpb.UpdatePlanetResourcesRequest) (*planetpb.UpdatePlanetResourcesResponse, error) {
	userID, err := uuid.Parse(req.GetUserID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	planetID, err := uuid.Parse(req.GetPlanetID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planet_id")
	}

	var updatedTime *time.Time
	if req.Time != nil {
		rt := req.Time.AsTime()
		updatedTime = &rt
	}

	if err := s.planetService.UpdatePlanetResources(ctx, userID, planetID, updatedTime); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &planetpb.UpdatePlanetResourcesResponse{}, nil
}
