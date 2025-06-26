package mqttgo

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yyle88/tern/zerotern"
	"go.uber.org/zap"
)

type CallbackState string

const (
	CallbackUnknown CallbackState = "unknown"
	CallbackRetries CallbackState = "retries"
	CallbackTimeout CallbackState = "timeout"
	CallbackSuccess CallbackState = "success"
)

func OnConnectWithRetries(client mqtt.Client, onConnect func(client mqtt.Client, retryTimes uint64) (CallbackState, error)) {
	for retryTimes := uint64(0); client.IsConnected(); retryTimes++ {
		action, err := onConnect(client, retryTimes)
		if err != nil {
			action = zerotern.VV(action, CallbackUnknown)
			switch action {
			case CallbackUnknown, CallbackRetries:
				log.ErrorLog("run-on-connect-with-retries", zap.String("action", string(action)), zap.Uint64("retry_times", retryTimes), zap.Error(err))
				time.Sleep(time.Millisecond * 100)
				continue
			case CallbackTimeout, CallbackSuccess:
				log.DebugLog("run-on-connect-with-retries", zap.String("action", string(action)), zap.Uint64("retry_times", retryTimes))
				return
			}
		}
		log.DebugLog("run-on-connect-success-done", zap.String("action", string(action)), zap.Uint64("retry_times", retryTimes))
		return
	}
}
