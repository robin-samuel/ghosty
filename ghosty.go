package ghosty

import (
	"context"

	"github.com/robin-samuel/ghosty/keyboard"
	"github.com/robin-samuel/ghosty/mouse"
)

func WithContext(ctx context.Context, width, height int) context.Context {
	ctx = keyboard.WithContext(ctx)
	ctx = mouse.WithContext(ctx, width, height)
	return ctx
}
