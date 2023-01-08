package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	url      = flag.String("u", "", "Url of target domain")
	rps      = flag.Int("n", 50, "Number of requests per second")
	duration = flag.Int("d", 2, "Duration of testing in seconds")
)

// create struct of color codes
type colors struct {
	yellow string
	green  string
	red    string
	reset  string
}

// assign color codes to color struct
var color = colors{
	yellow: "\033[33m",
	green:  "\033[32m",
	red:    "\033[31m",
	reset:  "\033[0m",
}

func main() {
	fmt.Println("\n\n--------------------------- Welcome to shex ---------------------------")
	flag.Parse()

	if *url == "" {
		fmt.Println("Please provide a target domain")
		return
	}

	if *rps < 1 {
		fmt.Println("Number of requests must be greater than 0")
		return
	}

	if *duration < 1 {
		fmt.Println("Duration must be greater than 0")
		return
	}

	fmt.Printf("Sending %d requests per second to %s for %d seconds\n", *rps, *url, *duration)

	// create a list of bool for responses, (success/failure)
	results := make([]bool, 0)

	var wg sync.WaitGroup

	// record starting time
	start := time.Now()

	// loop through the number of seconds
	for i := 0; i < *duration; i++ {
		// loop for the number of requests per second
		for j := 0; j < *rps; j++ {
			wg.Add(1)
			go sendRequest(&wg, &sync.Mutex{}, &results)
		}
		time.Sleep(time.Second * 1)
	}

	// time elapsed since start of sending requests
	finished := time.Now()
	timeUsed := finished.Sub(start)
	fmt.Printf("\nRequests sent in %v, waiting for responses...\n\n", timeUsed)
	wg.Wait()

	// record finishing time
	end := time.Now()

	// calculate time spent
	elapsed := end.Sub(start)

	fmt.Println("Load Testing took", elapsed)

	// calculate metrics
	metrics(&results)
}

func sendRequest(wg *sync.WaitGroup, mutex *sync.Mutex, results *[]bool) {
	mutex.Lock()
	//start := time.Now()
	_, err := http.Get(*url)
	//end := time.Now()
	mutex.Unlock()

	if err != nil {
		*results = append(*results, false)
		wg.Done()
	} else {
		*results = append(*results, true)
		wg.Done()
	}
}

func metrics(results *[]bool) {
	// calculate number of requests sent
	requestsSent := *rps * *duration

	// calculate number of requests that succeeded
	requestsSucceeded := 0
	for _, i := range *results {
		if i {
			requestsSucceeded++
		}
	}

	// calculate number of requests that failed
	requestsFailed := requestsSent - requestsSucceeded

	// calculate success rate
	successRate := float64(requestsSucceeded) / float64(requestsSent) * 100

	// calculate failure rate in two decimal places
	failureRate := float64(requestsFailed) / float64(requestsSent) * 100

	fmt.Println(color.yellow, "Requests sent:", requestsSent, color.reset)
	fmt.Println(color.green, "Requests succeeded:", requestsSucceeded, color.reset)
	fmt.Println(color.red, "Requests failed:", requestsFailed, color.reset)
	fmt.Println(color.green, "Success rate:", successRate, "%", color.reset)
	fmt.Println(color.red, "Failure rate:", failureRate, "%", color.reset)
}
