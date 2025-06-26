package mqttgo

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"go.uber.org/zap"
)

// NewClientWithCallback 客户端 cannot re-subscribe on reconnect，即，假如有订阅 topic，在客户端断开连接时就会断开订阅，而重连时却不能自动重新订阅。
// 在 https://github.com/eclipse/paho.mqtt.golang/issues/22 提到了解决方案
// 认为在 OnConnect 中调用订阅，是简单可靠的方案
func NewClientWithCallback(config *Config, clientID string, callback *Callback) (mqtt.Client, error) {
	clientOptions := NewClientOptions(must.Full(config), clientID)
	must.Full(callback)
	if onConnects := callback.onConnects; len(onConnects) > 0 {
		// 这里只是对 on-connect 做个简单的封装，直接用原版也是可以的，因为原版已经做的足够好也没什么可封装的
		// 根据 https://github.com/eclipse-paho/paho.mqtt.golang/blob/35b84c5b6910d3125376886939d0b36a8284d22a/client.go#L614
		// 这里 on-connect 是异步执行的
		clientOptions.OnConnect = func(client mqtt.Client) {
			log.DebugLog("client-connected-reconnected-run", zap.String("clientID", clientID))
			for _, onConnect := range onConnects {
				OnConnectWithRetries(client, onConnect)
			}
		}
	}
	return NewClient(clientOptions)
}

type Callback struct {
	onConnects []func(client mqtt.Client, retryTimes uint64) (CallbackState, error)
}

func NewCallback() *Callback {
	return &Callback{}
}

func (C *Callback) OnConnect(onConnect func(client mqtt.Client, retryTimes uint64) (CallbackState, error)) *Callback {
	C.onConnects = append(C.onConnects, onConnect)
	return C
}

func NewClient(clientOptions *mqtt.ClientOptions) (mqtt.Client, error) {
	client := mqtt.NewClient(clientOptions) //连接断开时会自动重连
	token := client.Connect()
	if ok := token.Wait(); !ok {
		return nil, erero.New("can-not-connect-mqtt")
	}
	if err := token.Error(); err != nil {
		return nil, erero.Wro(err) //返回错误的同时打印错误日志
	}
	return client, nil
}
