package configs

import "github.com/spf13/viper"

type Conf struct {
	AppPort          int    `mapstructure:"APP_PORT"`
	DBDriver         string `mapstructure:"DB_DRIVER"`
	DBHost           string `mapstructure:"DB_HOST"`
	DBPort           int    `mapstructure:"DB_PORT"`
	DBUser           string `mapstructure:"DB_USER"`
	DBPassword       string `mapstructure:"DB_PASSWORD"`
	DBName           string `mapstructure:"DB_NAME"`
	RabbitmqHost     string `mapstructure:"RABBITMQ_HOST"`
	RabbitmqPort     int    `mapstructure:"RABBITMQ_PORT"`
	RabbitmqUser     string `mapstructure:"RABBITMQ_USER"`
	RabbitmqPassword string `mapstructure:"RABBITMQ_PASSWORD"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
