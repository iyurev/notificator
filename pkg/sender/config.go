package sender

import (
	"github.com/iyurev/notificator/pkg/errors"
	"github.com/spf13/viper"
	"log"
)

const (
	configFileName    = "config.toml"
	envPrefix         = "TELEGRAM"
	defaultConfigPath = "/etc/hermes"
	defaultTgApiUrl   = "https://api.telegram.org"
	defaultTgApiToken = ""
)

var (
	globalConfig *GlobalConfig
	tgConfig     *TgConifg
)

func init() {
	var err error
	globalConfig, err = NewGlobalConfig()
	if err != nil {
		log.Fatalln(err)
	}
	tgConfig, err = NewTelegramConfig()
	if err != nil {
		log.Fatalln(err)
	}
}

type TgConifg struct {
	cfg *viper.Viper
}

func (tc *TgConifg) ApiUrl() string {
	return tc.cfg.GetString("api_url")
}

func (tc *TgConifg) ApiToken() string {
	return tc.cfg.GetString("api_token")
}

type TgRecipient struct {
	ChatID int64 `json:"chatID"`
}

type GlobalConfig struct {
	Telegram map[string]TgRecipient `json:"telegram"`
}

func (cfg *GlobalConfig) GetTgProjectRecipient(projectName string) (*TgRecipient, error) {
	recipient, ok := cfg.Telegram[projectName]
	if !ok {
		return nil, errors.NoSuchRecipient
	}
	return &recipient, nil
}

func NewGlobalConfig() (*GlobalConfig, error) {
	var globalConfig GlobalConfig
	config := viper.New()
	//Config file
	config.SetConfigName(configFileName)
	config.AddConfigPath(defaultConfigPath)
	config.AddConfigPath(".")
	config.AddConfigPath("$HOME/hermes")
	config.SetConfigType("toml")
	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}
	err := config.Unmarshal(&globalConfig)
	if err != nil {
		return nil, err
	}
	return &globalConfig, nil
}

func NewTelegramConfig() (*TgConifg, error) {

	config := viper.New()
	config.SetDefault("api_url", defaultTgApiUrl)
	config.SetDefault("api_token", defaultTgApiToken)
	config.AllowEmptyEnv(false)
	config.SetEnvPrefix(envPrefix)
	config.AutomaticEnv()
	return &TgConifg{
		cfg: config,
	}, nil
}
