package log

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

func TestKVToInterface(t *testing.T) {
	testcases := []struct {
		in       []attribute.KeyValue
		expected []interface{}
	}{{
		in:       nil,
		expected: nil,
	}, {
		in:       []attribute.KeyValue{attribute.String("foo", "bar")},
		expected: []interface{}{"foo", "bar"},
	}, {
		in:       []attribute.KeyValue{attribute.Bool("foo", true)},
		expected: []interface{}{"foo", true},
	}, {
		in:       []attribute.KeyValue{attribute.Float64("foo", 1.234)},
		expected: []interface{}{"foo", 1.234},
	}, {
		in:       []attribute.KeyValue{attribute.Int("foo", 123)},
		expected: []interface{}{"foo", int64(123)},
	}, {
		in:       []attribute.KeyValue{attribute.Int64("foo", 123)},
		expected: []interface{}{"foo", int64(123)},
	}, {
		in:       []attribute.KeyValue{attribute.Array("foo", []int{1, 2, 3})},
		expected: []interface{}{"foo", [3]int{1, 2, 3}},
	}}
	for _, tc := range testcases {
		out := kvToInterface(tc.in...)

		require.Len(t, out, len(tc.expected))
		for i, o := range out {
			require.Equal(t, tc.expected[i], o)
		}
	}
}
