package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type ServerConf struct {
	Host string
	Port int
}

type RedisConf struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}

type Config struct {
	Server ServerConf
	Redis  RedisConf
}

func GetApplicationConfig() *Config {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.SetConfigName("config")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	err = viper.BindEnv("redis.password", "REDIS_PASS")
	if err != nil {
		log.Fatal("err!", err)
	}

	viper.AddConfigPath(pwd + "/config")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("err!", err)
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Fatal("unable to decode config", err)
	}

	return conf
}
