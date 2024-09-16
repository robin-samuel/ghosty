package mouse

import (
	"context"
	"math/rand/v2"

	"github.com/robin-samuel/mimic"
)

type ContextKey string

var mousePositionKey = ContextKey("ghosty-mouse-position")

func WithContext(ctx context.Context, width, height int) context.Context {
	x := rand.IntN(width)
	y := rand.IntN(height)
	return context.WithValue(ctx, mousePositionKey, &mimic.Point{X: float64(x), Y: float64(y)})
}

func fromContext(ctx context.Context) *mimic.Point {
	if pos, ok := ctx.Value(mousePositionKey).(*mimic.Point); ok {
		return pos
	}
	return &mimic.Point{X: 0, Y: 0}
}
