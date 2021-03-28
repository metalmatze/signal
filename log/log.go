package log

import (
	"go.opentelemetry.io/otel/attribute"
)

type Logger interface {
	Log(keyvals ...attribute.KeyValue) error
}

func kvToInterface(keyvals ...attribute.KeyValue) []interface{} {
	kvs := []interface{}{}
	for _, kv := range keyvals {
		kvs = append(kvs, string(kv.Key))
		switch kv.Value.Type() {
		case attribute.BOOL:
			kvs = append(kvs, kv.Value.AsBool())
		case attribute.INT64:
			kvs = append(kvs, kv.Value.AsInt64())
		case attribute.FLOAT64:
			kvs = append(kvs, kv.Value.AsFloat64())
		case attribute.STRING:
			kvs = append(kvs, kv.Value.AsString())
		case attribute.ARRAY:
			kvs = append(kvs, kv.Value.AsArray())
		}
	}
	return kvs
}
