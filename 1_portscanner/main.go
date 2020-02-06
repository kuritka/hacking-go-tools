package main


/*
	port-scanner without waitgroups.
    If you use waiting groups it starts massive scanning by running
	goroutines and causes error states in workers
*/

import (
	"fmt"
	"net"
	"sort"
	"time"
)

/*
host = "localhost"

80 open
631 open
6942 open
63342 open

execution time 1.363029864s
Process finished with exit code 0
*/


const host = "localhost"
const numOfPorts = 65535
const closedOrTimeout = 0


func main() {
	start := time.Now()

	portsCh := make(chan int, 100)
	resultCh := make(chan int)
	var openPorts []int

	for i := 0; i < cap(portsCh); i++ {
		go worker(portsCh, resultCh)
	}

	go func(portNum int) {
		for i := 1;	i <= portNum; i++	{
			portsCh <- i
		}
	}(numOfPorts)

	for r := 0 ; r < numOfPorts; r++{
		port := <- resultCh
		if port != closedOrTimeout {
			openPorts = append(openPorts, port)
		}
	}

	close(portsCh)
	close(resultCh)

	sort.Ints(openPorts)
	for _,openPort := range openPorts{
		fmt.Printf("%d open \n",openPort)
	}

	fmt.Printf("\nexecution time %s", time.Since(start))
}


func worker(ports <- chan int, result chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%v",host,p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			//if closed or timeout than result = 0
			result <- closedOrTimeout
			continue
		}
		conn.Close()
		result <- p
	}
}