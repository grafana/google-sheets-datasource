package ext

import (
	"log"
	"strings"

	"github.com/go-logr/logr"
)

var _ logr.LogSink = (*logAdapter)(nil)

type logAdapter struct {
	log log.Logger
}

func newLogAdapter() *logAdapter {
	return &logAdapter{} //log: backend.Logger} //.New("k8s.apiserver")}
}

func (l *logAdapter) WithName(name string) logr.LogSink {
	//l.log = l.log.New("name", name)
	return l
}

func (l *logAdapter) WithValues(keysAndValues ...interface{}) logr.LogSink {
	//l.log = l.log.New(keysAndValues...)
	return l
}

func (l *logAdapter) Init(_ logr.RuntimeInfo) {
	// TODO: shrug emoji
}

func (l *logAdapter) Enabled(level int) bool {
	return level <= 5
}

func (l *logAdapter) Info(level int, msg string, keysAndValues ...interface{}) {
	// msg = strings.TrimSpace(msg)
	// if level < 1 {
	// 	l.log.Info(msg, keysAndValues...)
	// 	return
	// }
	// l.log.Debug(msg, keysAndValues...)
}

func (l *logAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	msg = strings.TrimSpace(msg)
	//l.log.Error(msg, keysAndValues...)
}
