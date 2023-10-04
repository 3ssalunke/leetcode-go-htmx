package util

import "github.com/spf13/viper"

type Config struct {
	RabbitMQHost      string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQQueueName string `mapstructure:"RABBITMQ_QUEUE_NAME"`
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
