package mqtt

import (
	"github.com/Jarnpher553/micro-core/log"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var MqttClient MQTT.Client

type Option func(options *MQTT.ClientOptions)

func Broker(broker string) Option {
	return func(options *MQTT.ClientOptions) {
		options.AddBroker(broker)
	}
}

func UserName(username string) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetUsername(username)
	}
}

func Pwd(pwd string) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetPassword(pwd)
	}
}

func Bind(options ...Option) {
	opts := MQTT.NewClientOptions()

	for _, op := range options {
		op(opts)
	}

	MqttClient = MQTT.NewClient(opts)
	if token := MqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Logger.Mark("MQTT").Fatalln(token.Error())
	}
}
