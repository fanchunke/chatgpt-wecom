package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App          `mapstructure:"app"`
	HTTP         `mapstructure:"http"`
	Logger       `mapstructure:"logger"`
	WeCom        `mapstructure:"wecom"`
	GPT          `mapstructure:"gpt"`
	Database     `mapstructure:"database"`
	Conversation `mapstructure:"conversation"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type HTTP struct {
	Port string `mapstructure:"port"`
}

type Logger struct {
	Level                 string `mapstructure:"level"`
	FileLoggingEnabled    bool   `mapstructure:"file_enabled"`
	ConsoleLoggingEnabled bool   `mapstructure:"console_enabled"`
	Filename              string `mapstructure:"filename"`
}

type GPT struct {
	ApiKey string `mapstructure:"api_key"`
}

type WeCom struct {
	CorpId         string `mapstructure:"corp_id"`
	CorpSecret     string `mapstructure:"corp_secret"`
	AgentId        int64  `mapstructure:"agent_id"`
	EncodingAESKey string `mapstructure:"encoding_aes_key"`
	Token          string `mapstructure:"token"`
	Url            string `mapstructure:"url"`
}

type Database struct {
	Driver     string `mapstructure:"driver"`
	DataSource string `mapstructure:"dataSource"`
}

type Conversation struct {
	CloseSessionFlag  string `mapstructure:"closeSessionFlag"`
	CloseSessionReply string `mapstructure:"closeSessionReply"`
	EnableEnterEvent  bool   `mapstructure:"enableEnterEvent"`
	EnterEventReply   string `mapstructure:"enterEventReply"`
}

func New(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Failed to Read configuration: %s", err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("Failed to Unmarshal configuration: %s", err)
	}
	return cfg, nil
}
