package mqttgo

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yyle88/tern/zerotern"
	"go.uber.org/zap"
)

type RetryType string

const (
	RetryTypeUnknown RetryType = "unknown"
	RetryTypeRetries RetryType = "retries"
	RetryTypeTimeout RetryType = "timeout"
	RetryTypeSuccess RetryType = "success"
)

func OnConnectWithRetries(c mqtt.Client, onConnect func(c mqtt.Client, retryTimes uint64) (RetryType, error)) {
	for retryTimes := uint64(0); c.IsConnected(); retryTimes++ {
		action, err := onConnect(c, retryTimes)
		if err != nil {
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
		log.DebugLog("run-on-connect-success-done", zap.String("action", string(action)), zap.Uint64("retry_times", retryTimes))
		return
	}
}
