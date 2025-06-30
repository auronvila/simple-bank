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
	ServerAddress        string        `mapstructure:"address"`
	TokenSymmetricKey    string        `mapstructure:"token_symmetric_key"`
	AccessTokenDuration  string        `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Optional: read file
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found, using environment variables only")
	}

	// Bind explicitly
	_ = viper.BindEnv("db_driver")
	_ = viper.BindEnv("db_source")
	_ = viper.BindEnv("address")
	_ = viper.BindEnv("token_symmetric_key")
	_ = viper.BindEnv("access_token_duration")
	_ = viper.BindEnv("refresh_token_duration")

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
	} else {
		fmt.Println("Unmarshal did not have any errors")
	}
	return
}
