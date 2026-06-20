package planet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/galaxy-empire-team/bridge-api/internal/service/planet/mocks"
	"github.com/galaxy-empire-team/bridge-api/pkg/consts"
)

func TestService_findAvailablePositionZ(t *testing.T) {
	tests := []struct {
		name                   string
		colonizedSystemPlanets []consts.PlanetPositionZ
		setupRand              func(rng *mocks.RandGenerator)
		wantPos                consts.PlanetPositionZ
		wantFound              bool
	}{
		{
			name: "skip colonized position at start",
			colonizedSystemPlanets: []consts.PlanetPositionZ{
				1, 5, 10,
			},
			setupRand: func(rng *mocks.RandGenerator) {
				rng.EXPECT().Uint32().Return(uint32(0)).Once()
			},
			wantPos:   2,
			wantFound: true,
		},
		{
			name: "wrap when start position is colonized",
			colonizedSystemPlanets: []consts.PlanetPositionZ{
				15,
			},
			setupRand: func(rng *mocks.RandGenerator) {
				rng.EXPECT().Uint32().Return(uint32(14)).Once()
			},
			wantPos:   1,
			wantFound: true,
		},
		{
			name: "wrap when tail is colonized",
			colonizedSystemPlanets: []consts.PlanetPositionZ{
				10, 11, 12, 14, 15,
			},
			setupRand: func(rng *mocks.RandGenerator) {
				rng.EXPECT().Uint32().Return(uint32(9)).Once()
			},
			wantPos:   1,
			wantFound: true,
		},
		{
			name: "no free positions in system",
			colonizedSystemPlanets: []consts.PlanetPositionZ{
				1, 2, 4, 5, 6, 8, 9, 10, 11, 12, 14, 15,
			},
			setupRand: func(rng *mocks.RandGenerator) {
				rng.EXPECT().Uint32().Return(uint32(0)).Once()
			},
			wantPos:   0,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := mocks.NewRandGenerator(t)
			tt.setupRand(rng)

			svc := &Service{
				randomGenerator: rng,
			}

			gotPos, gotFound := svc.findAvailablePositionZ(tt.colonizedSystemPlanets)

			require.Equal(t, tt.wantFound, gotFound)
			assert.Equal(t, tt.wantPos, gotPos)
		})
	}
}
