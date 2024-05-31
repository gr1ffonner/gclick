package config

import (
	"gclick/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env-default:"localhost"`
		Port string `yaml:"port" env-default:"8080"`
	}
	Database struct {
		Host     string `yaml:"host" env-required:"true"`
		Port     string `yaml:"port" env-required:"true"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database" env-required:"true"`
	} `yaml:"database" env-required:"true"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger("warn")
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
