package config

import (
	"flag"
	"github.com/Jarnpher553/gemini/pkg/log"
	"github.com/Jarnpher553/viper"
	"os"
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
	flag.Parse()

	if *filename == "" {
		logger.Fatal("config file hasn't set")
	}

	file(F(*filename))
	conf = viper.GetViper()
}

func Args() string {
	osArgs := os.Args

	file(Path("."), Name("config"), Type("toml"))

	var env string
	if len(osArgs) == 2 {
		env = osArgs[1]
	} else {
		env = "dev"
	}
	conf = viper.GetViper().Sub(env)

	if conf == nil {
		logger.Fatal("the config of os args is nil")
	}
	return env
}
