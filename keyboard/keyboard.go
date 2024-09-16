package keyboard

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strconv"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

const (
	typoRate       = 0.03 // Probability of a typo
	minDelay       = 60   // Min delay in milliseconds between key presses
	maxDelay       = 120  // Max delay in milliseconds between key presses
	backspaceDelay = 300  // Time taken to realize a mistake and press backspace
	beginDelayMin  = 500  // Min delay in milliseconds before starting to type
	beginDelayMax  = 800  // Max delay in milliseconds before starting to type
)

func Type(sel interface{}, v string, opts ...chromedp.QueryOption) chromedp.QueryAction {
	return chromedp.QueryAfter(sel, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return fmt.Errorf("selector %q did not return any nodes", sel)
		}

		n := nodes[0]

		// grab type attribute from node
		typ, attrs := "", n.Attributes
		n.RLock()
		for i := 0; i < len(attrs); i += 2 {
			if attrs[i] == "type" {
				typ = attrs[i+1]
			}
		}
		n.RUnlock()

		// when working with input[type="file"], call dom.SetFileInputFiles
		if n.NodeName == "INPUT" && typ == "file" {
			return dom.SetFileInputFiles([]string{v}).WithNodeID(n.NodeID).Do(ctx)
		}

		// focus on the input node
		err := dom.Focus().WithNodeID(n.NodeID).Do(ctx)
		if err != nil {
			return err
		}

		// Check if the input is a number
		var onlyDigits bool
		if _, err := strconv.Atoi(v); err == nil {
			onlyDigits = true
		}

		// Get the delay from the context
		delay := fromContext(ctx)

		// Create a list of actions to simulate typing
		var actions []chromedp.Action

		// Simulate delay in finding the keyboard/start position
		beginDelay := random(beginDelayMin, beginDelayMax)
		actions = append(actions, chromedp.Sleep(beginDelay))

		// Simulate typing
		for _, r := range v {
			// Simulate typo and backspace
			if rand.Float64() < typoRate {
				var wrong rune
				if onlyDigits {
					wrong = rune(rand.IntN(10) + '0')
				} else {
					wrong = rune(rand.IntN(26) + 'a')
				}

				// Press wrong key
				actions = append(actions, chromedp.KeyEvent(string(wrong)))
				actions = append(actions, chromedp.Sleep(delay+random(0, 20)))

				// Realize mistake and press backspace
				actions = append(actions, chromedp.Sleep(backspaceDelay+random(0, 20)))
				actions = append(actions, chromedp.KeyEvent(kb.Backspace))
				actions = append(actions, chromedp.Sleep(delay+random(0, 20)))
			}

			// Press correct key
			actions = append(actions, chromedp.KeyEvent(string(r)))
			actions = append(actions, chromedp.Sleep(delay+random(0, 20)))
		}

		return chromedp.Run(ctx, actions...)
	}, append(opts, chromedp.NodeVisible)...)
}
