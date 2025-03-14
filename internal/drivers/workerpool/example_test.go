package workerpool_test

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mattdowdell/sandbox/internal/drivers/workerpool"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func Example() {
	ctx := context.Background()

	handler := &Handler{}
	collector := &Collector{}

	pool, _ := workerpool.New(5 /*size*/, handler, collector)
	go pool.Start(ctx)

	for i := range 10 {
		if err := pool.Add(i); err != nil {
			slog.ErrorContext(ctx, "failed to add item to queue", slogx.Err(err))
			return
		}
	}

	if err := pool.Wait(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to wait", slogx.Err(err))
	}
}

type Handler struct{}

func (h *Handler) Handle(_ context.Context, input int) int {
	return input * 2
}

type Collector struct{}

func (c *Collector) Collect(output int) {
	fmt.Println(output)
}
