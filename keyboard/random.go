package keyboard

import (
	"math/rand/v2"
	"time"
)

func random(min int, max int) time.Duration {
	return time.Duration(rand.IntN(max-min)+min) * time.Millisecond
}
