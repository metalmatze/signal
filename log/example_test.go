package log

import (
	"os"

	a "go.opentelemetry.io/otel/attribute"
)

func ExampleLogftmLogger() {
	logger := NewLogftmLogger(os.Stdout)
	logger.Log(
		a.String("msg", "test"),
		a.Bool("boolean", true),
		a.Int("number", 123),
	)

	// Output:
	// msg=test boolean=true number=123
}

func ExampleJsonLogger() {
	logger := NewJSONLogger(os.Stdout)
	logger.Log(
		a.String("msg", "test"),
		a.Bool("boolean", true),
		a.Int("number", 123),
	)

	// Output:
	// {"boolean":true,"msg":"test","number":123}
}

//func ExampleContextual() {
//	logger := NewLogftmLogger(os.Stdout)
//	With(logger, a.String("foo", "bar"))
//}
