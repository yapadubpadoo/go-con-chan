package main

import (
	"fmt"
	"time"
)

func worker1(x int, c chan string) {
	time.Sleep(3 * time.Second)
	c <- fmt.Sprintf("worker1 result is %d", x*2)
}

func worker2(y int, c chan string) {
	c <- fmt.Sprintf("worker2 result is %d", y*3)
}

func main() {
	fmt.Println("This is main function")
	c1 := make(chan string)
	// c2 := make(chan int)
	go worker1(5, c1)
	go worker2(10, c1)
	fmt.Println(<-c1)
	fmt.Println(<-c1)
}
