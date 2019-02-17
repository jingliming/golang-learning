## 代码结构
```shell
mygo
├── bin
├── debug
├── main.exe
├── main.go
├── pkg
└── src
    ├── channel
    │   └── mychannel.go
    └── coroutine
        └── mycoroutine.go
```
## 携程go-coroutine
### 示例代码
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
### 遇到的问题
* 包名就是文件夹名称，不同的文件可以使用相同的包名，只要是在同一个文件夹下就可以
* 包导出的函数必须是首字母大写的

## 通道channel
### 普通channel
#### 示例代码
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
### 带缓存的channel
#### 示例代码
mychannle.go
```
package channel

import (
	"fmt"
)

.......
/*注意如果通道缓存设置为1，但是写入两个缓存，就会爆出错误
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
### 带同步的channel
#### 示例代码
mychannel.go
```
package channel

import (
	"fmt"
	"time"
)

.......

/*channel sync
* 注意此处的<-done，如果没有这行代码，可能存在worker携程没有执行就退出了
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
### 通道作为参数时约束通道方向
#### 示例代码
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
### 使用select完成通道选择
#### 示例代码
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
### channel超时处理
#### 示例代码
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
### channle非阻塞
通过select的default分支实现
### channel关闭
mychannel.go
#### 示例代码
```
/*channel close
* Go携程原理，如果携程被关闭并且没有输出，more为false
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
### channel遍历
#### 示例代码
mychannel.go
```
//channel traverse
//一个非空的通道也是可以关闭的，但是通道中剩下的值仍然可以被接收到
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
### 遇到的问题
* 使用go func启用携程时后面必须带一个()，从而成功调用该匿名函数，否则会提示`[go] expression in go must be function call`