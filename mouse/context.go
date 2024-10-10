package mouse

import (
	"context"
	"math/rand/v2"

	"github.com/robin-samuel/mimic"
)

type ContextKey string

var mousePositionKey = ContextKey("ghosty-mouse-position")
var viewportKey = ContextKey("ghosty-viewport")

func WithContext(ctx context.Context, width, height int) context.Context {
	x := rand.IntN(width)
	y := rand.IntN(height)
	ctx = context.WithValue(ctx, mousePositionKey, &mimic.Point{X: float64(x), Y: float64(y)})
	ctx = context.WithValue(ctx, viewportKey, &mimic.Viewport{Width: float64(width), Height: float64(height)})
	return ctx
}

func fromContext(ctx context.Context) *mimic.Point {
	if pos, ok := ctx.Value(mousePositionKey).(*mimic.Point); ok {
		return pos
	}
	return &mimic.Point{X: 0, Y: 0}
}

func viewportFromContext(ctx context.Context) *mimic.Viewport {
	if viewport, ok := ctx.Value(viewportKey).(*mimic.Viewport); ok {
		return viewport
	}
	return &mimic.Viewport{Width: 0, Height: 0}
}
