package mqttgo

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yyle88/erero"
)

type TokenState string

const (
	TokenStateUnknown TokenState = "unknown"
	TokenStateTimeout TokenState = "timeout"
	TokenStateSuccess TokenState = "success"
)

func CheckToken(token mqtt.Token, timeout time.Duration) (TokenState, error) {
	if success := token.WaitTimeout(timeout); !success {
		return TokenStateTimeout, erero.New(string(TokenStateTimeout))
	}
	if err := token.Error(); err != nil {
		return TokenStateUnknown, erero.Wro(err)
	}
	return TokenStateSuccess, nil
}

func WaitToken(token mqtt.Token) (TokenState, error) {
	if !token.Wait() {
		return TokenStateUnknown, erero.New("wait-token-not-complete")
	}
	if err := token.Error(); err != nil {
		return TokenStateUnknown, erero.Wro(err)
	}
	return TokenStateSuccess, nil
}
