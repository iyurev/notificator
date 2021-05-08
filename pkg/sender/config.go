package sender

import (
	"errors"
	"github.com/spf13/viper"
)

const (
	configFileName    = "config.toml"
	envPrefix         = "TELEGRAM"
	defaultConfigPath = "/etc/hermes"
	defaultTgApiUrl   = ""
	defaultTgApiToken = ""
)

var (
	globalConfig GlobalConfig
)

type Recipient struct {
	ChatID int64 `json:"chatID"`
}

type GlobalConfig struct {
	Telegram map[string]Recipient `json:"telegram"`
}

func (cfg *GlobalConfig) GetProjectRecipient(projectName string) (*Recipient, error) {
	recipient, ok := cfg.Telegram[projectName]
	if !ok {
		return nil, errors.New("there's no recipient with such name")
	}
	return &recipient, nil
}

func NewGlobalConfig() (*viper.Viper, error) {
	config := viper.New()
	//Config file
	config.SetConfigName(configFileName)
	config.AddConfigPath(defaultConfigPath)
	config.AddConfigPath(".")
	config.SetConfigType("toml")
	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}
	err := config.Unmarshal(&globalConfig)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func NewTelegramConfig() (*viper.Viper, error) {
	config := viper.New()
	config.SetDefault("api_url", defaultTgApiUrl)
	config.SetDefault("api_token", defaultTgApiToken)
	config.AllowEmptyEnv(false)
	config.SetEnvPrefix(envPrefix)
	config.AutomaticEnv()
	return config, nil
}
