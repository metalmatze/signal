package prometheus

import (
	"fmt"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	a "go.opentelemetry.io/otel/attribute"
)

func ExampleCounter() {
	c := NewCounterVec(prom.CounterOpts{
		Name: "counter",
		Help: "some counter",
	}, []a.Key{"code", "method"})

	c.With(a.Int("code", 200), a.String("method", "GET")).Inc()
	c.With(a.Int("code", 500), a.String("method", "PUT")).Inc()

	fmt.Println(testutil.CollectAndCount(c, "counter"))
	// Output:
	// 2
}
