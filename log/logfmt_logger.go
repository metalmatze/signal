package log

import (
	"io"

	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/attribute"
)

type logftmLogger struct {
	l log.Logger
}

func NewLogftmLogger(w io.Writer) Logger {
	return logftmLogger{l: log.NewLogfmtLogger(w)}
}

func (l logftmLogger) Log(keyvals ...attribute.KeyValue) error {
	return l.l.Log(kvToInterface(keyvals...)...)
}
