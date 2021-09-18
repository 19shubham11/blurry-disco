package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
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

func setEnvConfig(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = fallback
	}

	return value
}

func loadConfig() (*Config, error) {
	// Relative on runtime DIR:
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("error resolving pwd")
	}

	dir := path.Join(path.Dir(b))

	file, err := os.Open(dir + "/config.json")
	if err != nil {
		return nil, err
	}

	appConfig := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(appConfig)

	return appConfig, err
}

func GetApplicationConfig() *Config {
	conf, err := loadConfig()
	if err != nil {
		fmt.Printf("unable to load application config %v", err)
		panic(err)
	}

	// set env vars
	conf.Redis.Password = setEnvConfig("REDIS_PASS", "default")

	return conf
}
