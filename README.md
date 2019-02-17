## ????
```shell
mygo
??? bin
??? debug
??? main.exe
??? main.go
??? pkg
??? src
    ??? channel
    ?   ??? mychannel.go
    ??? coroutine
        ??? mycoroutine.go
```
## ??go-coroutine
### ????
main.go
```golang
package main

import "mygo/src/coroutine"

func main() {
	coroutine.CorotineTest()
}

```
mycoroutine.go
```golang
package coroutine

import "fmt"

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

func CorotineTest() {
	f("direct")
	go f("goroutine")
	go func(msg string) {
		fmt.Println(msg)
	}("going")

	var input string
	fmt.Scanln(&input)
	fmt.Println("done")
}

```
### ?????
* ???????????????????????????????????????
* ???????????????

## ??channel
### ??channel
#### ????
main.go
```
package main

import (
	"mygo/src/channel"
	
)

func main() {
	channel.MychannelTest()
}

```
mychannel.go
```
package channel

import (
	"fmt"
)

func MychannelTest() {
	message := make(chan string)
	go func() {
		message <- "ping"
	}()
	msg := <-message
	fmt.Println(msg)
}

```
### ????channel
#### ????
mychannle.go
```
package channel

import (
	"fmt"
)

.......
/*???????????1????????????????
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send]:
mygo/src/channel.MychannelWithBuf()
        E:/GoPath/src/mygo/src/channel/mychannel.go:19 +0x7d
main.main()
        E:/GoPath/src/mygo/main.go:11 +0x27
*/        
func MychannelWithBuf() {
	message := make(chan string, 1)
	message <- "buffered"
	message <- "channel"
	fmt.Println(<-message)
	fmt.Println(<-message)
}

```
### ????channel
#### ????
mychannel.go
```
package channel

import (
	"fmt"
	"time"
)

.......

/*channel sync
* ?????<-done??????????????worker??????????
*/
func MychannelWithSync() {
	done := make(chan bool, 1)
	go worker(done)
	<-done
}

func worker(done chan bool) {
	fmt.Println("working....")
	time.Sleep(time.Second)
	fmt.Println("done")
	done <- true
}

```
### ?????????????
#### ????
```
package channel

.......

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

```
### ??select??????
#### ????
mychannel.go
```
package channel

import (
	"fmt"
	"time"
)

.......

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

```
### channel????
#### ????
mychannel.go
```
package channel

import (
	"fmt"
	"time"
)
// channel time expire
func MychannelWithTimeExpire() {
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		time.Sleep(time.Second * 2)
		c1 <- "result 1"
	}()
	select {
	case res := <-c1:
		fmt.Println(res)
	case <-time.After(time.Second * 2):
		fmt.Println("result 1 timeout!!!")
	}
	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "result 2"
	}()
	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(time.Second * 1):
		fmt.Println("result 2 timeout!!!")
	}
}
```
### channle???
??select?default????
### channel??
mychannel.go
#### ????
```
/*channel close
* Go???????????????????more?false
*/
func MychannelWithClose() {
	jobs := make(chan int, 5)
	done := make(chan bool)
	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()
	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("send job", j)

	}
	close(jobs)
	fmt.Println("send all jobs")
	<-done
}
```
### channel??
#### ????
mychannel.go
```
//channel traverse
//????????????????????????????????
func MychannelWithTraverse() {
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)
	for elem := range queue {
		fmt.Println(elem)
	}
}
```
### ?????
* ??go func????????????()??????????????????`[go] expression in go must be function call`
