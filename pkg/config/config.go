package config

import (
	"flag"
	"github.com/Jarnpher553/gemini/pkg/log"
	"github.com/Jarnpher553/viper"
)

type Config = viper.Viper

var logger = log.Zap.Mark("config")

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

func F(f string) Option {
	return func(v *viper.Viper) {
		v.SetConfigFile(f)
	}
}

func Provider(provider string, endpoint string, keyOrPath string) Option {
	return func(v *viper.Viper) {
		if err := v.AddRemoteProvider(provider, endpoint, keyOrPath); err != nil {
			logger.Fatal(log.Message(err))
		}
	}
}

func File() {
	filename := flag.String("conf", "", "config file name")
	sub := flag.String("sub", "", "multi sub of single config")
	flag.Parse()

	if *filename == "" {
		logger.Fatal("config file hasn't set")
	}

	file(F(*filename))

	if *sub != "" {
		conf = viper.GetViper().Sub(*sub)
	} else {
		conf = viper.GetViper()
	}

	if conf == nil {
		logger.Fatal("sub of config is error")
	}
}
