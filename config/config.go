package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"strings"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	//PgConf     PostgresConfig `yaml:"postgres"`
	AmqpConf   AmqpConfig `yaml:"amqp"`
	OtlpConfig OtlpConfig `yaml:"otlp_config"`
	TgConfig   TgConfig   `yaml:"tg_config"`
}

type AmqpConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	UserName     string `yaml:"user_name"`
	UserPass     string `yaml:"user_pass"`
	QueueName    string `yaml:"queue"`
	ExchangeName string `yaml:"exchange"`
	RoutingKey   string `yaml:"routing_key"`
}

type OtlpConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	ServiceName string `yaml:"service_name"`
}

type TgConfig struct {
	BaseURL string `yaml:"base_url"`
	Token   string `yaml:"token"`
}

func (r AmqpConfig) GetAmqpUri() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.UserName, r.UserPass, r.Host, r.Port)
}

// MustLoad загружает конфигурацию из нескольких файлов, переопределяя значения.
func MustLoad() *Config {
	configPaths := fetchConfigPaths()
	if len(configPaths) == 0 {
		panic("no config paths provided")
	}

	var cfg Config

	for _, path := range configPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic("config file does not exist: " + path)
		}

		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			panic("failed to read config: " + err.Error())
		}
	}

	return &cfg
}

// fetchConfigPaths получает список путей к файлам конфигурации из флага командной строки или переменной окружения.
func fetchConfigPaths() []string {
	var paths string

	flag.StringVar(&paths, "config", "", "comma-separated list of config files")
	flag.Parse()

	if paths == "" {
		paths = os.Getenv("CONFIG_PATH")
	}

	if paths == "" {
		return nil
	}

	return splitAndTrim(paths)
}

// splitAndTrim разбивает строку по запятой и удаляет лишние пробелы.
func splitAndTrim(input string) []string {
	parts := strings.Split(input, ",")
	var result []string
	for _, path := range parts {
		trimmed := strings.TrimSpace(path)
		if trimmed != "" && strings.HasSuffix(trimmed, ".yaml") {
			result = append(result, trimmed)
		}
	}
	return result
}
