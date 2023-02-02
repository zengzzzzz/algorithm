package limitFlow

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

// limit flow golang rate
func LimitFlowAllow() {
	l := rate.NewLimiter(rate.Every(time.Second/10), 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				if l.Allow() {
					fmt.Printf("allow %d\n", i)
				}
				time.Sleep(time.Second / 2)
			}
		}(i)
	}
	time.Sleep((time.Second * 10))
}

func LimitFlowWait() {
	l := rate.NewLimiter(rate.Every(time.Second/10), 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			if err := l.Wait(context.TODO()); err != nil {
			} else {
				fmt.Printf("allow %d\n", i)
			}
			time.Sleep(time.Second / 2)
		}(i)
	}
	time.Sleep((time.Second * 10))
}

func LimitFlowReserve() {
	l := rate.NewLimiter(rate.Every(time.Second/10), 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			if r := l.Reserve(); r.OK() {
				time.Sleep(r.Delay())
				fmt.Printf("allow %d\n", i)
			}
		}(i)
	}
	time.Sleep((time.Second * 10))
}
