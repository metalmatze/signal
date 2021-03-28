package log

import (
	"io"

	"github.com/go-kit/kit/log"
	"go.opentelemetry.io/otel/attribute"
)

type jsonLogger struct {
	l log.Logger
}

// NewJSONLogger returns a Logger that encodes keyvals to the Writer as a
// single JSON object. Each log event produces no more than one call to
// w.Write. The passed Writer must be safe for concurrent use by multiple
// goroutines if the returned Logger will be used concurrently.
func NewJSONLogger(w io.Writer) Logger {
	return &jsonLogger{l: log.NewJSONLogger(w)}
}

func (j *jsonLogger) Log(keyvals ...attribute.KeyValue) error {
	return j.l.Log(kvToInterface(keyvals...)...)
}
