package configs

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port    string `yaml:"port" env:"PORT" env-default:"8080"`
	DbUrl   string `yaml:"db_url" env:"DB_URL"`
	Tg      TelegramConfig
	Payment PaymentServiceConfig
	Partner PartnerParamsConfig
}

type TelegramConfig struct {
	Debug bool   `env:"DEBUG"`
	Token string `yaml:"token" env:"TELEGRAM_TOKEN"`
}

type PaymentServiceConfig struct {
	BaseUrl     string `yaml:"base_url" env:"PAYMENT_SERVICE_BASE_URL"`
	PostbackUrl string `yaml:"postback_url" env:"PAYMENT_SERVICE_POSTBACK_URL"`
	SuccessUrl  string `yaml:"success_url" env:"PAYMENT_SERVICE_SUCCESS_URL"`
	FailUrl     string `yaml:"fail_url" env:"PAYMENT_SERVICE_FAIL_URL"`
}

type PartnerParamsConfig struct {
	DiscountPercent   uint8 `yaml:"discount_percent" env:"PARTNER_DISCOUNT_PERCENT"`
	CommissionPercent uint8 `yaml:"commission_percent" env:"PARTNER_COMMISSION_PERCENT"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg

}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
