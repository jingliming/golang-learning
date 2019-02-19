package statuscoroutine

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

type readOp struct {
	key  int
	resp chan int
}

type writeOp struct {
	key  int
	val  int
	resp chan bool
}

func MyStatusCoroutine() {
	var ops int64
	reads := make(chan *readOp)
	writes := make(chan *writeOp)

	go func() {
		states := make(map[int]int)
		for {
			select {
			case read := <-reads:
				resp := states[read.key]
				read.resp <- resp
			case write := <-writes:
				states[write.key] = write.val
				write.resp <- true
			}
		}
	}()

	for r := 0; r < 100; r++ {
		go func() {
			for {
				read := &readOp{
					key:  rand.Intn(5),
					resp: make(chan int),
				}
				reads <- read
				<-read.resp
				atomic.AddInt64(&ops, 1)
			}

		}()
	}

	for w := 0; w < 10; w++ {
		go func() {
			for {
				write := &writeOp{
					key:  rand.Intn(5),
					val:  rand.Intn(100),
					resp: make(chan bool),
				}
				writes <- write
				<-write.resp
				atomic.AddInt64(&ops, 1)
			}
		}()
	}
	time.Sleep(time.Second)
	opsFinal := atomic.LoadInt64(&ops)
	fmt.Println("opsFinal:", opsFinal)
}
