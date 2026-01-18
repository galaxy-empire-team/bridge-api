package planet

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"initialservice/internal/models"
)

const (
	galaxyCount        = 1
	systemInGalaxyCount = 3
	planetsInSystemCount = 16
)

type planetRepository interface {
	GetCapitol(ctx context.Context, userID uuid.UUID) (models.Planet, error)
	ColonizeCapitol(ctx context.Context, userID uuid.UUID, planet models.Planet) error 
}

type Service struct {
	planetRepo planetRepository
	randomGenerator *rand.Rand
}

func New(repo planetRepository) *Service {
	gen := rand.New(rand.NewSource((time.Now().UnixNano())))

	return &Service{
		planetRepo: repo,
		randomGenerator: gen,
	}
}
