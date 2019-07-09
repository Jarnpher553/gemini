package config

import (
	"github.com/Jarnpher553/micro-core/log"
	"github.com/Jarnpher553/viper"
)

type Config struct {
	*viper.Viper
}

var (
	conf *Config
)

func Conf() *Config {
	if conf == nil {
		log.Logger.Mark("Config").Fatalln("must be constructed before")
	}
	return conf
}

// Option 配置项方法
type Option func(*viper.Viper)

// Path 路径配置
func Path(path string) Option {
	return func(v *viper.Viper) {
		v.AddConfigPath(path)
	}
}

// Name 文件名配置
func Name(name string) Option {
	return func(v *viper.Viper) {
		v.SetConfigName(name)
	}
}

// Type 文件类型配置
func Type(_type string) Option {
	return func(v *viper.Viper) {
		v.SetConfigType(_type)
	}
}

func Provider(provider string, endpoint string, keyOrPath string) Option {
	return func(v *viper.Viper) {
		if err := v.AddRemoteProvider(provider, endpoint, keyOrPath); err != nil {
			log.Logger.Mark("Config").Fatalln(err)
		}
	}
}
