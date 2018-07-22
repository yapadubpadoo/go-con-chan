package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func profilingStart() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
}

func profilingEnd() {
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}

func main() {
	profilingStart()

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

	time.Sleep(5 * time.Second)
	fmt.Println("Main done")

	profilingEnd()
}
