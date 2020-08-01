package config

import (
	"github.com/Jarnpher553/micro-core/log"
	"github.com/Jarnpher553/viper"
)

type Config = viper.Viper

var logger = log.Zap.Mark("Config")

//默认使用当前文件夹下的config.yaml文件
//func init() {
//	viper.AddConfigPath(".")
//	viper.SetConfigName("config")
//	viper.SetConfigType("yaml")
//
//	generate()
//}

var conf *Config

func Conf() *Config {
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
			logger.Fatal(log.Message(err))
		}
	}
}

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

func RunMode(osArgs []string, factories ...func()) string {
	var runMode string
	if len(osArgs) == 1 {
		runMode = DebugMode
	} else {
		switch osArgs[1] {
		case "debug":
			runMode = DebugMode
		case "test":
			runMode = TestMode
		case "release":
			runMode = ReleaseMode
		default:
			runMode = DebugMode
		}
	}
	if factories != nil && len(factories) != 0 {
		factories[0]()
	}

	conf = viper.GetViper().Sub(runMode)

	return runMode
}
