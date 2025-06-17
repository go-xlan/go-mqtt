package mqttgo

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yyle88/must"
	"go.uber.org/zap"
)

type Config struct { //配置1个或者1批客户端信息
	BrokerServer string
	Username     string
	Password     string
	OrderMatters bool
}

func NewClientOptions(cfg *Config, clientID string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().
		AddBroker(must.Nice(cfg.BrokerServer)).
		SetClientID(must.Nice(clientID)).
		SetUsername(cfg.Username).
		SetPassword(cfg.Password).
		SetOrderMatters(cfg.OrderMatters) //在订阅端是否是有序的，假如设为true则订阅的处理函数不要阻塞消费，否则阻塞所有流程

	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.DebugLog("default-publish-handle-function", zap.String("clientID", clientID), zap.String("topic", msg.Topic()))
	}) // 设置消息回调处理函数
	opts.SetPingTimeout(1 * time.Second)
	opts.OnConnect = func(client mqtt.Client) {
		log.DebugLog("client-on-connected-on-reconnected-callback", zap.String("clientID", clientID))
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.ErrorLog("client-on-lost-connection-callback", zap.String("clientID", clientID), zap.Error(err))
	}
	return opts
}
