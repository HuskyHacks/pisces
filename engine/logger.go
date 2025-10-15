package engine

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(debug bool) *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("source", "go").
		Logger().
		Sample(
			zerolog.LevelSampler{
				TraceSampler: &zerolog.BurstSampler{
					Burst:       1,
					Period:      2 * time.Second,
					NextSampler: &zerolog.BasicSampler{N: 100},
				},
				WarnSampler: &zerolog.BurstSampler{
					Burst:       4,
					Period:      1 * time.Second,
					NextSampler: &zerolog.BasicSampler{N: 100},
				},
			},
		)

	return &logger
}
