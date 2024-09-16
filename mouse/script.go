package mouse

import (
	"context"
	_ "embed"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//go:embed mouse.js
var MouseJS string

func ShowCursor() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		if _, err := page.AddScriptToEvaluateOnNewDocument(MouseJS).Do(ctx); err != nil {
			return err
		}
		return nil
	})
}
