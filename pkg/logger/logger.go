package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger Default init for logging system
func InitLogger(name, profile string) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	log.Logger = log.With().Str("name", name).Str("logEnv", profile).Logger().Output(os.Stdout)
}

func Stack(err error) string {
	return fmt.Sprintf("%+v", err)
}
