package atomiccounter

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

func MyAutomicCounter() {
	var ops uint64 = 0
	for i := 0; i < 50; i++ {
		go func() {
			for {
				atomic.AddUint64(&ops, 1)
				//让当前协程让出cpu，让其他协程也会在未来的某个时间点继续运行
				runtime.Gosched()
			}
		}()
	}
	time.Sleep(time.Second)
	opsFinal := atomic.LoadUint64(&ops)
	fmt.Println("ops:", opsFinal)
}
