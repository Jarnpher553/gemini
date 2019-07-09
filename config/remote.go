package config

import (
	"gitee.com/jarnpher_rice/micro-core/log"
	"github.com/Jarnpher553/viper"
)
import _ "github.com/Jarnpher553/viper/remote"

func Remote(opts ...Option) {
	conf = &Config{
		Viper: viper.New(),
	}

	for _, opt := range opts {
		opt(conf.Viper)
	}

	conf.generateRemote()
}

// generate 配置生成
func (c *Config) generateRemote() {
	if err := c.ReadRemoteConfig(); err != nil {
		log.Logger.Mark("Config").Fatalln(err)
	}
}
