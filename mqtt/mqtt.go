package mqtt

import (
	"crypto/tls"
	"github.com/Jarnpher553/gemini/log"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"time"
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

func ResumeSubs(resume bool) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetResumeSubs(resume)
	}
}

func ClientID(id string) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetClientID(id)
	}
}

type CP = MQTT.CredentialsProvider

func CredentialsProvider(p CP) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetCredentialsProvider(p)
	}
}

func CleanSession(clean bool) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetCleanSession(clean)
	}
}

func OrderMatters(order bool) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetOrderMatters(order)
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetTLSConfig(t)
	}
}

type ST = MQTT.Store

func Store(s ST) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetStore(s)
	}
}

func KeepAlive(k time.Duration) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetKeepAlive(k)
	}
}

func PingTimeout(k time.Duration) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetPingTimeout(k)
	}
}

func ProtocolVersion(pv uint) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetProtocolVersion(pv)
	}
}

func DisableWill() Option {
	return func(options *MQTT.ClientOptions) {
		options.UnsetWill()
	}
}

func Will(topic, payload string, qos byte, retained bool) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetWill(topic, payload, qos, retained)
	}
}

func BinaryWill(topic string, payload []byte, qos byte, retained bool) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetBinaryWill(topic, payload, qos, retained)
	}
}

type MH = MQTT.MessageHandler

func DefaultPublishHandler(defaultHandler MH) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetDefaultPublishHandler(defaultHandler)
	}
}

type OCH = MQTT.OnConnectHandler

func OnConnectHandler(onConn OCH) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetOnConnectHandler(onConn)
	}
}

type CLH = MQTT.ConnectionLostHandler

func ConnectionLostHandler(onLost CLH) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetConnectionLostHandler(onLost)
	}
}

type RH = MQTT.ReconnectHandler

func ReconnectiongHandler(rh RH) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetReconnectingHandler(rh)
	}
}

func WriteTimeout(t time.Duration) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetWriteTimeout(t)
	}
}

func ConnectTimeout(t time.Duration) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetConnectTimeout(t)
	}
}

func MaxReconnectInterval(t time.Duration) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetMaxReconnectInterval(t)
	}
}

func AutoReconnect(a bool) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetAutoReconnect(a)
	}
}

func ConnectRetryInterval(t time.Duration) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetConnectRetryInterval(t)
	}
}

func ConnectRetry(a bool) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetConnectRetry(a)
	}
}

func MessageChannelDepth(s uint) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetMessageChannelDepth(s)
	}
}

func HTTPHeaders(h http.Header) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetHTTPHeaders(h)
	}
}

type WsOptions = MQTT.WebsocketOptions

func WebsocketOptions(w *WsOptions) Option {
	return func(options *MQTT.ClientOptions) {
		options.SetWebsocketOptions(w)
	}
}

func Bind(options ...Option) {
	opts := MQTT.NewClientOptions()

	for _, op := range options {
		op(opts)
	}

	MqttClient = MQTT.NewClient(opts)
	if token := MqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Zap.Mark("mqtt").Fatal(log.Message(token.Error()))
	}
}
