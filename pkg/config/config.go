package config

import "github.com/spf13/viper"

type Config struct {
	Source string
	Target string
	Debug  bool
	Port   string
}

var Instance = &Config{}

func Bind() {
	Instance.Source = viper.GetString("source")
	Instance.Target = viper.GetString("target")
	Instance.Debug = viper.GetBool("debug")
	Instance.Port = viper.GetString("port")
}
