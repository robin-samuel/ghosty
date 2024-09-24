package mouse

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/robin-samuel/mimic"
)

// Click is an element query action that moves the cursor and sends a mouse
// click event to the first element node matching the selector.
func Click(sel string, opts ...func(*chromedp.Selector)) chromedp.QueryAction {
	return chromedp.QueryAfter(sel, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector %q did not return any nodes", sel)
		}

		n := nodes[0]

		if err := dom.ScrollIntoViewIfNeeded().WithNodeID(n.NodeID).Do(ctx); err != nil {
			return err
		}

		boxes, err := dom.GetContentQuads().WithNodeID(n.NodeID).Do(ctx)
		if err != nil {
			return err
		}

		if len(boxes) == 0 {
			return chromedp.ErrInvalidDimensions
		}

		content := boxes[0]

		if len(content) != 8 {
			return chromedp.ErrInvalidDimensions
		}

		x1, y1 := content[0], content[1]
		x2, y2 := content[2], content[3]
		x3, y3 := content[4], content[5]
		x4, y4 := content[6], content[7]

		minX := min(x1, x2, x3, x4)
		maxX := max(x1, x2, x3, x4)
		minY := min(y1, y2, y3, y4)
		maxY := max(y1, y2, y3, y4)

		padX := (maxX - minX) * 0.25
		padY := (maxY - minY) * 0.25

		minX += padX
		maxX -= padX
		minY += padY
		maxY -= padY

		xf := minX + rand.Float64()*(maxX-minX)
		yf := minY + rand.Float64()*(maxY-minY)

		pos := fromContext(ctx)
		x0, y0 := pos.X, pos.Y

		delta := math.Sqrt(math.Pow(xf-x0, 2) + math.Pow(yf-y0, 2))

		path := mimic.Generate(mimic.Config{
			Points: []mimic.Point{
				{X: x0, Y: y0},
				{X: xf, Y: yf},
			},
			Duration:  time.Duration(delta*3) * time.Millisecond,
			Noise:     0.2,
			Frequency: 60,
		})

		var tasks chromedp.Tasks
		for _, point := range path {
			tasks = append(tasks, chromedp.MouseEvent(input.MouseMoved, point.X, point.Y))
		}

		tasks = append(tasks, MouseClickXY(xf, yf))
		pos.X, pos.Y = xf, yf

		return tasks.Do(ctx)
	}, append(opts, chromedp.NodeVisible)...)
}

// MouseClickXY is an action that sends a left mouse button click (i.e.,
// mousePressed and mouseReleased event) to the X, Y location.
func MouseClickXY(x, y float64, opts ...chromedp.MouseOption) chromedp.MouseAction {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		p := &input.DispatchMouseEventParams{
			Type:       input.MousePressed,
			X:          x,
			Y:          y,
			Button:     input.Left,
			ClickCount: 1,
		}

		// apply opts
		for _, o := range opts {
			p = o(p)
		}

		if err := p.Do(ctx); err != nil {
			return err
		}

		if err := chromedp.Sleep(random(75, 95)).Do(ctx); err != nil {
			return err
		}

		p.Type = input.MouseReleased
		return p.Do(ctx)
	})
}

func Select(sel string, v string, opts ...func(*chromedp.Selector)) chromedp.Action {
	return chromedp.Tasks{
		Click(sel, opts...),
		chromedp.Sleep(random(500, 1500)),
		chromedp.SendKeys(sel, v, opts...),
	}
}

func SelectValue(sel string, v string, opts ...func(*chromedp.Selector)) chromedp.Action {
	return chromedp.Tasks{
		Click(sel, opts...),
		chromedp.Sleep(random(500, 1500)),
		chromedp.SetValue(sel, v, opts...),
	}
}
