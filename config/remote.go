package config

import (
	"github.com/Jarnpher553/micro-core/log"
	"github.com/Jarnpher553/viper"
)
import _ "github.com/Jarnpher553/viper/remote"

func Remote(opts ...Option) {
	v := viper.GetViper()

	for _, opt := range opts {
		opt(v)
	}

	generateRemote()
}

// generate 配置生成
func generateRemote() {
	if err := viper.ReadRemoteConfig(); err != nil {
		log.Logger.Mark("Config").Fatalln(err)
	}
}
