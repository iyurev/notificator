package sender

import (
	"github.com/spf13/viper"
)

type Config struct {
	Telegram *viper.Viper
}

const (
	configFileName    = "config.toml"
	envPrefix         = "TELEGRAM"
	defaultTgApiUrl   = ""
	defaultTgApiToken = ""
)

func NewTelegramConfig() (*viper.Viper, error) {
	config := viper.New()
	config.SetDefault("api_url", defaultTgApiUrl)
	config.SetDefault("api_token", defaultTgApiToken)
	//Config file
	config.SetConfigName(configFileName)

	config.AllowEmptyEnv(false)
	config.SetEnvPrefix(envPrefix)
	config.AutomaticEnv()
	return config, nil
}
