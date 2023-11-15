package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUSER     string `mapstructure:"DB_USER"`
	DBPASSWORD string `mapstructure:"DB_PASSWORD"`
	DBNAME     string `mapstructure:"DB_NAME"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// logger.Error("Error while read config:"+ err.Error())
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		// logger.Error("Error while unmarshal:"+ err.Error())
		return nil, err
	}
	return config, nil
}