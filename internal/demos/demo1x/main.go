package main

import (
	"encoding/json"
	"math/rand/v2"
	"time"

	"github.com/go-xlan/go-mqtt/internal/utils"
	"github.com/go-xlan/go-mqtt/mqttgo"
	"github.com/pkg/errors"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func main() {
	const topic = "demo1_send_msg_topic"

	config := &mqttgo.Config{
		BrokerServer: "ws://127.0.0.1:8083/mqtt",
		Username:     "username",
		Password:     "password",
		OrderMatters: false,
	}
	onConnect := func(c mqttgo.Client, retryTimes uint64) (mqttgo.RetryType, error) {
		if retryTimes > 10 {
			return mqttgo.RetryTypeTimeout, nil
		}
		if rand.IntN(100) >= 10 {
			return mqttgo.RetryTypeUnknown, erero.New("random-rate-not-success")
		}
		return mqttgo.RetryTypeSuccess, nil
	}
	client1 := rese.V1(mqttgo.NewClient(config, utils.NewUUID(), onConnect))
	defer client1.Disconnect(500)

	client2 := rese.V1(mqttgo.NewClient(config, utils.NewUUID(), func(c mqttgo.Client, retryTimes uint64) (mqttgo.RetryType, error) {
		token := c.Subscribe(topic, 1, func(client mqttgo.Client, message mqttgo.Message) {
			zaplog.SUG.Debugln("subscribe-msg:", neatjsons.SxB(message.Payload()))
		})
		tokenState, err := mqttgo.WaitToken(token)
		if err != nil {
			return mqttgo.RetryTypeRetries, errors.WithMessage(err, "subscribe-is-wrong")
		}
		must.Same(tokenState, mqttgo.TokenStateSuccess)
		return mqttgo.RetryTypeSuccess, nil
	}))
	defer client2.Disconnect(500)

	type MessageType struct {
		A string
		B int
		C float64
	}

	for i := 0; i < 100; i++ {
		msg := &MessageType{
			A: time.Now().String(),
			B: i,
			C: rand.Float64(),
		}
		contentBytes := rese.A1(json.Marshal(msg))

		zaplog.SUG.Debugln("publish-msg:", neatjsons.SxB(contentBytes))

		token := client1.Publish(topic, 1, false, contentBytes)
		tokenState := rese.C1(mqttgo.CheckToken(token, time.Second*3))
		must.Same(tokenState, mqttgo.TokenStateSuccess)
		time.Sleep(time.Second)
	}
}
