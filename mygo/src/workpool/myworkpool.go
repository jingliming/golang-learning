package workpool

import (
	"fmt"
	"time"
)

/*
* 构建3个worker并发处理生产的job，一开始没有job，worker阻塞等待
 */
func MyWorkPool() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)
	for a := 1; a <= 9; a++ {
		<-results
	}
}

/*
* jobs通道指定只允许输出，results指定只允许输入
 */
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}
