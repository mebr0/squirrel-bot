package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
)

var (
	cfg *Config
)

type Config struct {
	Log struct {
		Level string `yaml:"level" envconfig:"LOG_LEVEL"`
	} `yaml:"log"`

	Telegram struct {
		BotToken string `yaml:"bot-token" envconfig:"TELEGRAM_BOT_TOKEN"`
	} `yaml:"telegram"`

	DB struct {
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     string `yaml:"port" envconfig:"DB_PORT"`
		Name     string `yaml:"name" envconfig:"DB_NAME"`
		User     string `yaml:"user" envconfig:"DB_USER"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD"`
		SSLMode  string `yaml:"ssl-mode" envconfig:"DB_SSL_MODE"`
	} `yaml:"db"`

	Game Game `yaml:"game"`
}

type Game struct {
	Debug bool          `yaml:"debug" envconfig:"GAME_DEBUG"`
	Speed time.Duration `yaml:"speed" envconfig:"GAME_SPEED"`
}

func LoadConfig(configPath string) *Config {
	if cfg == nil {
		cfg = &Config{}

		cfg.readFile(configPath)
		cfg.readEnv()
	}

	return cfg
}

// File configs with values from configs file
func (c *Config) readFile(path string) {
	f, err := os.Open(path)

	if err != nil {
		processError(err)
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(c)

	if err != nil {
		log.Println(c)
		processError(err)
	}
}

// Read configs with values from env variables
func (c *Config) readEnv() {
	loadFromEnvFile()

	err := envconfig.Process("", c)

	if err != nil {
		processError(err)
	}
}

// Load values from .env file to system
func loadFromEnvFile() {
	if err := godotenv.Load(); err != nil {
		zap.L().Debug("Error loading .env file")
	}
}

func processError(err error) {
	zap.L().Error(err.Error())
	os.Exit(2)
}
