package utils

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Host []*Host `yaml:"host"`
}

type Host struct {
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`
	Auth Auth   `yaml:"auth"`
}

type Auth struct {
	Username     string `yaml:"username"`
	Passwd       string `yaml:"passwd"`
	GoogleSecret string `yaml:"googleSecret"`
}

func InitConfig(configFile string) *Config {

	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	var _config *Config
	err = viper.Unmarshal(&_config)
	if err != nil {
		fmt.Println(err.Error())
	}
	return _config
}
