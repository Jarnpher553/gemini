package config

import (
	"github.com/Jarnpher553/micro-core/log"
	"github.com/Jarnpher553/viper"
)

func init() {
	conf = &Config{
		Viper: viper.New(),
	}

	conf.AddConfigPath(".")
	conf.SetConfigName("config")
	conf.SetConfigType("yaml")

	conf.Generate()
}

// File 构造函数
func File(options ...Option) {
	conf = &Config{
		Viper: viper.New(),
	}

	for i := range options {
		options[i](conf.Viper)
	}

	conf.Generate()
}

// generate 配置生成
func (c *Config) Generate() {
	err := c.ReadInConfig()
	if err != nil {
		log.Logger.Mark("Config").Fatalln(err)
	}
}
