package utils

import "github.com/spf13/viper"

type Config struct {
	DBUser                string `mapstructure:"DB_USER"`
	DBPassword            string `mapstructure:"DB_PASSWORD"`
	DBDatabase            string `mapstructure:"DB_DATABASE"`
	DBPort                string `mapstructure:"DB_PORT"`
	DBHost                string `mapstructure:"DB_HOST"`
	ApiPort               string `mapstructure:"API_PORT"`
	PrometheusPushgateway string `mapstructure:"PROMETHEUS_PUSHGATEWAY"`
	JwtSecret             string `mapstructure:"JWT_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
