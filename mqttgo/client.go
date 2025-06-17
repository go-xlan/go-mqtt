package mqttgo

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/tern/zerotern"
	"go.uber.org/zap"
)

// NewClient 客户端 cannot re-subscribe on reconnect，即，假如有订阅 topic，在客户端断开连接时就会断开订阅，而重连时却不能自动重新订阅。
// 在 https://github.com/eclipse/paho.mqtt.golang/issues/22 提到了解决方案
// 认为在 OnConnect 中调用订阅，是简单可靠的方案
func NewClient(cfg *Config, clientID string, onConnects ...func(c mqtt.Client, retryTimes uint64) (RetryType, error)) (mqtt.Client, error) {
	clientOptions := NewClientOptions(must.Full(cfg), clientID)
	if len(onConnects) > 0 {
		clientOptions.OnConnect = func(client mqtt.Client) {
			log.DebugLog("client-connected-reconnected-run", zap.String("clientID", clientID))
			for _, onceOnConnect := range onConnects {
				runOnConnectWithRetries(client, onceOnConnect)
			}
		}
	}
	return NewClientConnect(clientOptions)
}

type RetryType string

const (
	RetryTypeUnknown RetryType = "unknown"
	RetryTypeRetries RetryType = "retries"
	RetryTypeTimeout RetryType = "timeout"
	RetryTypeSuccess RetryType = "success"
)

func runOnConnectWithRetries(c mqtt.Client, onceOnConnect func(c mqtt.Client, retryTimes uint64) (RetryType, error)) {
	for retryTimes := uint64(0); c.IsConnected(); retryTimes++ {
		if action, err := onceOnConnect(c, retryTimes); err != nil {
			action = zerotern.VV(action, RetryTypeUnknown)
			switch action {
			case RetryTypeUnknown, RetryTypeRetries:
				log.ErrorLog("run-on-connect-with-retries", zap.String("action", string(action)), zap.Uint64("retry_times", retryTimes), zap.Error(err))
				time.Sleep(time.Millisecond * 100)
				continue
			case RetryTypeTimeout, RetryTypeSuccess:
				log.DebugLog("run-on-connect-with-retries", zap.String("action", string(action)), zap.Uint64("retry_times", retryTimes))
				return
			}
		}
		log.DebugLog("run-on-connect-success-done", zap.Uint64("retry_times", retryTimes))
		return
	}
}

func NewClientConnect(clientOptions *mqtt.ClientOptions) (mqtt.Client, error) {
	client := mqtt.NewClient(clientOptions) //连接断开时会自动重连
	token := client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, erero.Wro(err) //返回错误的同时打印错误日志
	}
	return client, nil
}
