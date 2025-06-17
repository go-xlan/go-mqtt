package mqttgo

import (
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

type Log interface {
	ErrorLog(msg string, fields ...zap.Field)
	DebugLog(msg string, fields ...zap.Field)
}

var log Log = newLog()

func SetLog(mqttLog Log) {
	log = mqttLog
}

type zapLog struct{}

func newLog() *zapLog {
	return &zapLog{}
}

func (z *zapLog) ErrorLog(msg string, fields ...zap.Field) {
	zaplog.LOGS.Skip(1).Error(msg, fields...)
}

func (z *zapLog) DebugLog(msg string, fields ...zap.Field) {
	zaplog.LOGS.Skip(1).Debug(msg, fields...)
}
