package util

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Config struct {
	DbDriver             string        `mapstructure:"db_driver"`
	DbSource             string        `mapstructure:"db_source"`
	HttpServerAddress    string        `mapstructure:"http_server_address"`
	GrpcServerAddress    string        `mapstructure:"grpc_server_address"`
	MigrationUrl         string        `mapstructure:"migration_url"`
	TokenSymmetricKey    string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration  string        `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigFile("app.env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	fmt.Println("Trying to load config from path:", path)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found, using environment variables only")
	}

	config = Config{
		DbDriver:            strings.TrimSpace(viper.GetString("db_driver")),
		DbSource:            strings.TrimSpace(viper.GetString("db_source")),
		HttpServerAddress:   strings.TrimSpace(viper.GetString("http_server_address")),
		GrpcServerAddress:   strings.TrimSpace(viper.GetString("grpc_server_address")),
		MigrationUrl:        strings.TrimSpace(viper.GetString("migration_url")),
		TokenSymmetricKey:   strings.TrimSpace(viper.GetString("token_symmetric_key")),
		AccessTokenDuration: strings.TrimSpace(viper.GetString("access_token_duration")),
		RefreshTokenDuration: func() time.Duration {
			dur, _ := time.ParseDuration(strings.TrimSpace(viper.GetString("refresh_token_duration")))
			return dur
		}(),
	}

	// debug purposes
	//fmt.Printf("Loaded config: %+v\n", config)
	return
}
