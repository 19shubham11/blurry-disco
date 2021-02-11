package conf

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type ServerConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RedisConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
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

	confContent, err := ioutil.ReadFile(pwd + "/cmd/conf/config.yml")
	if err != nil {
		panic(err)
	}
	confContent = []byte(os.ExpandEnv(string(confContent)))
	conf := &Config{}
	if err := yaml.Unmarshal(confContent, conf); err != nil {
		panic(err)
	}
	return conf
}
