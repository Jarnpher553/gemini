package config

import (
	"github.com/Jarnpher553/gemini/log"
	"github.com/Jarnpher553/viper"
)
import _ "github.com/Jarnpher553/viper/remote"

func Remote(opts ...Option) Factory {
	return func() {
		v := viper.GetViper()

		for _, opt := range opts {
			opt(v)
		}

		generateRemote()
	}
}

// generate 配置生成
func generateRemote() {
	if err := viper.ReadRemoteConfig(); err != nil {
		logger.Fatal(log.Message(err))
	}
}
