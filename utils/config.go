package utils

import "github.com/spf13/viper"

type Config struct {
	DBUser                             string `mapstructure:"DB_USER"`
	DBPassword                         string `mapstructure:"DB_PASSWORD"`
	DBDatabase                         string `mapstructure:"DB_DATABASE"`
	DBPort                             string `mapstructure:"DB_PORT"`
	DBHost                             string `mapstructure:"DB_HOST"`
	ApiPort                            string `mapstructure:"API_PORT"`
	PrometheusPushgateway              string `mapstructure:"PROMETHEUS_PUSHGATEWAY"`
	JwtSecret                          string `mapstructure:"JWT_SECRET"`
	ObjectStorageServerUseSsl          bool   `mapstructure:"OBJECT_STORAGE_SERVER_USE_SSL"`
	ObjectStorageServerAccessKeyId     string `mapstructure:"OBJECT_STORAGE_SERVER_ACCESS_KEY_ID"`
	ObjectStorageServerSecretAccessKey string `mapstructure:"OBJECT_STORAGE_SERVER_SECRET_ACCESS_KEY"`
	ObjectStorageServerEndpoint        string `mapstructure:"OBJECT_STORAGE_SERVER_ENDPOINT"`
	BusinessAvatarBulkName             string `mapstructure:"BUSINESS_AVATAR_BULK_NAME"`
	ItemsBulkName                      string `mapstructure:"ITEMS_BULK_NAME"`
	UsersBulkName                      string `mapstructure:"USERS_BULK_NAME"`
	UsersDeleteObjectsBulkName         string `mapstructure:"USERS_DELETE_OBJECTS_BULK_NAME"`
	ItemsDeletedBulkName               string `mapstructure:"ITEMS_DELETED_BULK_NAME"`
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
