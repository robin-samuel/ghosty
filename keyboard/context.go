package keyboard

import (
	"context"
	"time"
)

type ContextKey string

var typeDelayKey = ContextKey("ghosty-type-delay")

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, typeDelayKey, random(minDelay, maxDelay))
}

func fromContext(ctx context.Context) time.Duration {
	if delay, ok := ctx.Value(typeDelayKey).(time.Duration); ok {
		return delay
	}
	return random(minDelay, maxDelay)
}
