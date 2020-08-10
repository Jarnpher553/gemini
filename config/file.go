package config

import (
	"github.com/Jarnpher553/micro-core/log"
	"github.com/Jarnpher553/viper"
)

// File 构造函数
func File(options ...Option) Factory {
	return func() {
		v := viper.GetViper()
		for i := range options {
			options[i](v)
		}

		generate()
	}
}

// generate 配置生成
func generate() {
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal(log.Message(err))
	}
}
