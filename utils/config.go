package utils

import "github.com/spf13/viper"

type Config struct {
	Environment                        string `mapstructure:"ENVIRONMENT"`
	DBUser                             string `mapstructure:"DB_USER"`
	Tls                                string `mapstructure:"TLS"`
	DBDsn                              string `mapstructure:"DB_DSN"`
	DBPassword                         string `mapstructure:"DB_PASSWORD"`
	DBDatabase                         string `mapstructure:"DB_DATABASE"`
	DBPort                             int    `mapstructure:"DB_PORT"`
	DBHost                             string `mapstructure:"DB_HOST"`
	ApiPort                            int    `mapstructure:"API_PORT"`
	PrometheusPushgateway              string `mapstructure:"PROMETHEUS_PUSHGATEWAY"`
	JwtSecret                          string `mapstructure:"JWT_SECRET"`
	ObjectStorageServerUseSsl          string `mapstructure:"OBJECT_STORAGE_SERVER_USE_SSL"`
	ObjectStorageServerAccessKeyId     string `mapstructure:"OBJECT_STORAGE_SERVER_ACCESS_KEY_ID"`
	ObjectStorageServerSecretAccessKey string `mapstructure:"OBJECT_STORAGE_SERVER_SECRET_ACCESS_KEY"`
	ObjectStorageServerEndpoint        string `mapstructure:"OBJECT_STORAGE_SERVER_ENDPOINT"`
	BusinessAvatarBulkName             string `mapstructure:"BUSINESS_AVATAR_BULK_NAME"`
	BusinessAvatarDeletedBulkName      string `mapstructure:"BUSINESS_AVATAR_DELETED_BULK_NAME"`
	ItemsBulkName                      string `mapstructure:"ITEMS_BULK_NAME"`
	UsersBulkName                      string `mapstructure:"USERS_BULK_NAME"`
	UsersDeletedBulkName               string `mapstructure:"USERS_DELETED_BULK_NAME"`
	ItemsDeletedBulkName               string `mapstructure:"ITEMS_DELETED_BULK_NAME"`
	EmailHostname                      string `mapstructure:"EMAIL_HOSTNAME"`
	EmailSmtpPort                      int    `mapstructure:"EMAIL_SMTP_PORT"`
	EmailAddress                       string `mapstructure:"EMAIL_ADDRESS"`
	EmailAddressPassword               string `mapstructure:"EMAIL_ADDRESS_PASSWORD"`
	AppName                            string `mapstructure:"APP_NAME"`
	RedisHost                          string `mapstructure:"REDIS_HOST"`
	RedisPort                          int    `mapstructure:"REIDS_PORT"`
	RedisPassword                      string `mapstructure:"REDIS_PASSWORD"`
	RedisDb                            int    `mapstructure:"REDIS_DB"`
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
