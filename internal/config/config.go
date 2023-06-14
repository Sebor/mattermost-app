package config

import (
	"strconv"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/rs/zerolog/log"
)

var productGroup = map[string]string{
	"pr1": "PR1_DevOps",
	"pr2": "PR1_DevOps",
	"pr3": "PR2_DevOps",
	"pr4": "PR3_DevOps",
	"pr5": "PR4_DevOps",
}

type Config struct {
	AppPort               int           `env:"APP_PORT" envDefault:"4000"`
	LogLevel              int           `env:"LOG_LEVEL" envDefault:"1"` // https://pkg.go.dev/github.com/rs/zerolog@v1.29.1#Level
	LogTimestampFieldName string        `env:"LOG_TIMESTAMP_FIELD_NAME" envDefault:"@timestamp"`
	AppEnv                string        `env:"APP_ENV" envDefault:"development"`
	AppDisplayName        string        `env:"APP_DISPLAY_NAME" envDefault:"Jira bot"`
	AppRootURL            string        `env:"APP_ROOT_URL" envDefault:"http://mattermost-bot"`
	IsRootUrlIngress      bool          `env:"IS_ROOT_URL_INGRESS" envDefault:"false"`
	AppHomePage           string        `env:"APP_HOME_PAGE" envDefault:"https://github.com/Sebor/mattermost-app"`
	APIURL                string        `env:"API_URL" envDefault:"http://host.docker.internal:9999/services/v1"`
	InsecureSkipVerify    bool          `env:"INSECURE_SKIP_VERIFY" envDefault:"false"`
	HTTPClientTimeout     time.Duration `env:"HTTP_CLIENT_TIMEOUT" envDefault:"30s"`
	TeamName              string        `env:"TEAM_NAME" envDefault:"dev"`
	PostEmojiName         string        `env:"POST_EMOJI_NAME" envDefault:"ballot_box_with_check"`
	ChannelAllowedIDs     channels
	ProductGroup          map[string]string
}

type channels struct {
	IDs []string `env:"CHANNEL_ALLOWED_IDS" envSeparator:"," envDefault:"jd6j73r4546647567fe1b6r"`
}

func (ch channels) GetChannels() map[string]struct{} {
	m := make(map[string]struct{})
	for _, v := range ch.IDs {
		m[v] = struct{}{}
	}
	return m
}

func GetConfig() Config {
	cfg := Config{}
	cfg.ProductGroup = productGroup
	if err := env.Parse(&cfg); err != nil {
		log.Error().Err(err).Msg("")
	}
	if !cfg.IsRootUrlIngress {
		cfg.AppRootURL = cfg.AppRootURL + ":" + strconv.Itoa(cfg.AppPort)
	}
	return cfg
}
