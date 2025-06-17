package main

import (
	"math/rand/v2"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-xlan/go-mqtt/internal/utils"
	"github.com/go-xlan/go-mqtt/mqttgo"
	"github.com/pkg/errors"
	"github.com/yyle88/erero"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func main() {
	const topic = "demo1_send_msg_topic"

	config := &mqttgo.Config{
		BrokerServer: "ws://127.0.0.1:8083/mqtt",
		Username:     "username",
		Password:     "password",
		OrderMatters: false,
	}
	onConnect := func(c mqtt.Client, retryTimes uint64) (mqttgo.RetryType, error) {
		if retryTimes > 10 {
			return mqttgo.RetryTypeTimeout, nil
		}
		if rand.IntN(100) >= 10 {
			time.Sleep(time.Second * 3)
			return mqttgo.RetryTypeUnknown, erero.New("random-rate-not-success")
		}
		return mqttgo.RetryTypeSuccess, nil
	}
	client1 := rese.V1(mqttgo.NewClient(config, utils.NewUUID(), onConnect))
	defer client1.Disconnect(500)

	client2 := rese.V1(mqttgo.NewClient(config, utils.NewUUID(), func(c mqtt.Client, retryTimes uint64) (mqttgo.RetryType, error) {
		token := c.Subscribe(topic, 1, func(client mqtt.Client, message mqtt.Message) {
			zaplog.SUG.Debugln("subscribe-msg:", string(message.Payload()))
		})
		if ok := token.Wait(); !ok {
			return mqttgo.RetryTypeRetries, errors.New("subscribe-is-wrong")
		}
		return mqttgo.RetryTypeSuccess, nil
	}))
	defer client2.Disconnect(500)

	for i := 0; i < 100; i++ {
		content := time.Now().String()
		zaplog.SUG.Debugln("publish-msg:", content)
		token := client1.Publish(topic, 1, false, content)
		if success := token.WaitTimeout(time.Second * 3); !success {
			zaplog.LOG.Debug("publish-msg-timeout")
		}
		if err := token.Error(); err != nil {
			zaplog.LOG.Error("publish-msg", zap.Error(err))
		}
		time.Sleep(time.Second)
	}
}
