package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func callSija(endpoint string, c chan string) {
	r, _ := http.Get(endpoint)
	body, _ := ioutil.ReadAll(r.Body)
	c <- fmt.Sprintf("%s", string(body))
}

func mainNormal() {
	fmt.Println("mainNormal start")
	endpoint := "http://192.168.1.54:3000/posts"

	maxLoop := 20
	for i := 0; i < maxLoop; i++ {
		r, _ := http.Get(fmt.Sprintf("%s/%d", endpoint, i))
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Println("mainNormal", string(body))
	}
	fmt.Println("mainNormal done")
}

func mainWithManualSleep() {
	fmt.Println("mainWithManualSleep start")

	maxLoop := 20
	endpoint := "http://192.168.1.54:3000/posts"
	c := make(chan string)

	for i := 0; i < maxLoop; i++ {
		go callSija(fmt.Sprintf("%s/%d", endpoint, i), c)
	}

	time.Sleep(5 * time.Second)
	for i := 0; i < maxLoop; i++ {
		fmt.Println("mainWithManualSleep", <-c)
	}

	fmt.Println("mainWithManualSleep done")
}

func main() {
	fmt.Println("Main start")

	maxLoop := 20
	endpoint := "http://192.168.1.54:3000/posts"
	c := make(chan string)

	go mainNormal()
	go mainWithManualSleep()

	for i := 0; i < maxLoop; i++ {
		go callSija(fmt.Sprintf("%s/%d", endpoint, i), c)
	}

	for i := 0; i < maxLoop; i++ {
		fmt.Println("main", <-c)
	}

	fmt.Println("Main done")
}
