package main

import (
	"encoding/json"
	"math/rand/v2"
	"time"

	"github.com/go-xlan/go-mqtt/internal/utils"
	"github.com/go-xlan/go-mqtt/mqttgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func main() {
	const mqttTopic = "mqtt-go-demo1-topic"

	config := &mqttgo.Config{
		BrokerServer: "ws://127.0.0.1:8083/mqtt",
		Username:     "username",
		Password:     "password",
		OrderMatters: false,
	}
	client1 := rese.V1(mqttgo.NewClientWithCallback(config, utils.NewUUID(), mqttgo.NewCallback().
		OnConnect(func(c mqttgo.Client, retryTimes uint64) (mqttgo.CallbackState, error) {
			return mqttgo.CallbackSuccess, nil
		}),
	))
	defer client1.Disconnect(500)

	client2 := rese.V1(mqttgo.NewClientWithCallback(config, utils.NewUUID(), mqttgo.NewCallback().
		OnConnect(func(c mqttgo.Client, retryTimes uint64) (mqttgo.CallbackState, error) {
			token := c.Subscribe(mqttTopic, 1, func(client mqttgo.Client, message mqttgo.Message) {
				zaplog.SUG.Debugln("subscribe-msg:", neatjsons.SxB(message.Payload()))
			})
			must.Same(rese.C1(mqttgo.CheckToken(token, time.Minute)), mqttgo.TokenStateSuccess)
			return mqttgo.CallbackSuccess, nil
		}),
	))
	defer client2.Disconnect(500)

	type MessageType struct {
		A string
		B int
		C float64
	}

	for i := 0; i < 10; i++ {
		msg := &MessageType{
			A: time.Now().String(),
			B: i,
			C: rand.Float64(),
		}
		contentBytes := rese.A1(json.Marshal(msg))

		zaplog.SUG.Debugln("publish-msg:", neatjsons.SxB(contentBytes))

		token := client1.Publish(mqttTopic, 1, false, contentBytes)
		must.Same(rese.C1(mqttgo.CheckToken(token, time.Second*3)), mqttgo.TokenStateSuccess)
		time.Sleep(time.Second)
	}
}
