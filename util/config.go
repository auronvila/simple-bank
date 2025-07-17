package util

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type Config struct {
	Environment          string        `mapstructure:"environment"`
	DbDriver             string        `mapstructure:"db_driver"`
	DbSource             string        `mapstructure:"db_source"`
	HttpServerAddress    string        `mapstructure:"http_server_address"`
	GrpcServerAddress    string        `mapstructure:"grpc_server_address"`
	MigrationUrl         string        `mapstructure:"migration_url"`
	TokenSymmetricKey    string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration  string        `mapstructure:"access_token_duration"`
	RedisAddress         string        `mapstructure:"redis_address"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigFile("app.env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	// * because the variables are not loaded manually print the logs pretty
	prettyOutput := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	prettyOutput.Info().Msgf("Trying to load config from path: %s", path)

	if err := viper.ReadInConfig(); err != nil {
		prettyOutput.Error().Msg("No config file found, using environment variables only")
	}

	config = Config{
		Environment:         strings.TrimSpace(viper.GetString("environment")),
		DbDriver:            strings.TrimSpace(viper.GetString("db_driver")),
		DbSource:            strings.TrimSpace(viper.GetString("db_source")),
		HttpServerAddress:   strings.TrimSpace(viper.GetString("http_server_address")),
		GrpcServerAddress:   strings.TrimSpace(viper.GetString("grpc_server_address")),
		MigrationUrl:        strings.TrimSpace(viper.GetString("migration_url")),
		TokenSymmetricKey:   strings.TrimSpace(viper.GetString("token_symmetric_key")),
		AccessTokenDuration: strings.TrimSpace(viper.GetString("access_token_duration")),
		RedisAddress:        strings.TrimSpace(viper.GetString("redis_address")),
		RefreshTokenDuration: func() time.Duration {
			dur, _ := time.ParseDuration(strings.TrimSpace(viper.GetString("refresh_token_duration")))
			return dur
		}(),
	}

	// debug purposes
	//fmt.Printf("Loaded config: %+v\n", config)
	return
}
