package util

import "github.com/spf13/viper"

type Config struct {
	DBHost                  string `mapstructure:"DB_HOST"`
	DBPort                  string `mapstructure:"DB_PORT"`
	DBUser                  string `mapstructure:"DB_USER"`
	DBPass                  string `mapstructure:"DB_PASS"`
	DBName                  string `mapstructure:"DB_NAME"`
	TokenSecret             string `mapstructure:"TOKEN_SECRET"`
	AppPort                 string `mapstructure:"APP_PORT"`
	GoogleOAuthClientId     string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOAuthClientSecret string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	RabbitMQHost            string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQQueueName       string `mapstructure:"RABBITMQ_QUEUE_NAME"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
