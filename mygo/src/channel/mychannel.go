package channel

import (
	"fmt"
	"time"
)

func MychannelTest() {
	message := make(chan string)
	go func() {
		message <- "ping"
	}()
	msg := <-message
	fmt.Println(msg)
}

//channel cache
func MychannelWithBuf() {
	message := make(chan string, 2)
	message <- "buffered"
	message <- "channel"
	fmt.Println(<-message)
	fmt.Println(<-message)
}

//channel sync
func MychannelWithSync() {
	done := make(chan bool, 1)
	go worker(done)
	<-done
}

// goroutine code
func worker(done chan bool) {
	fmt.Println("working....")
	time.Sleep(time.Second)
	fmt.Println("done")
	done <- true
}

// channel direction
func MychannelWithDirection() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
func ping(pings chan<- string, msg string) {
	pings <- msg
}
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

// channel selector
func MychannelWithSelector() {
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		time.Sleep(time.Second * 1)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "two"
	}()
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received", msg1)
		case msg2 := <-c2:
			fmt.Println("received", msg2)
		}
	}
}
