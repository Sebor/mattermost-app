package logger

import (
	"io"
	"mattermost-bot/internal/config"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	once sync.Once
	log  zerolog.Logger
	cfg  = config.GetConfig()
)

func GetLogger() zerolog.Logger {
	once.Do(func() {
		zerolog.TimestampFieldName = cfg.LogTimestampFieldName

		logLevel := cfg.LogLevel

		var output io.Writer = zerolog.MultiLevelWriter(os.Stderr)

		if (cfg.AppEnv == "development") || (cfg.AppEnv == "develop") || (cfg.AppEnv == "dev") {
			output = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}
			logLevel = -1
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Logger()
	})

	return log
}
