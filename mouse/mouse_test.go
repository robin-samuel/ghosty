package mouse_test

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/robin-samuel/ghosty/mouse"
)

func TestClick(t *testing.T) {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoFirstRun,
		chromedp.Flag("disable-search-engine-choice-screen", true),
		chromedp.WindowSize(1920, 1080),
	)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	if err := chromedp.Run(ctx,
		mouse.ShowCursor(),
		chromedp.Navigate("https://www.example.com"),
		mouse.Click("body > div > p:nth-child(3) > a"),
		chromedp.Sleep(time.Minute),
	); err != nil {
		t.Fatal(err)
	}
}
