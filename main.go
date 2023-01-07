package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	targetDomain      = flag.String("u", "", "Target domain url")
	requestsPerSecond = flag.Int("n", 50, "Number of requests per second")
	duration          = flag.Int("d", 2, "Duration of the attack in seconds")
)

func main() {
	fmt.Println("--------------------------- Welcome to shex ---------------------------")
	flag.Parse()

	if *targetDomain == "" {
		fmt.Println("Please provide a target domain")
		return
	}

	if *requestsPerSecond < 1 {
		fmt.Println("Number of requests must be greater than 0")
		return
	}

	if *duration < 1 {
		fmt.Println("Duration must be greater than 0")
		return
	}

	fmt.Printf("Sending %d requests per second to %s for %d seconds\n", *requestsPerSecond, *targetDomain, *duration)

	// Create channel to receive results
	results := make([]bool, 0)

	// record starting time
	start := time.Now()

	var wg sync.WaitGroup

	// loop through the number of seconds
	for i := 0; i < *duration; i++ {
		// loop for the number of requests per second
		for j := 0; j < *requestsPerSecond; j++ {
			wg.Add(1)
			go sendRequest(&wg, &sync.Mutex{}, &results)
		}
		time.Sleep(time.Second * 1)
	}
	fmt.Println("Requests sent, waiting for responses...")
	wg.Wait()

	// record finishing time
	end := time.Now()

	// calculate time spent
	elapsed := end.Sub(start)

	fmt.Println("Load Testing took", elapsed)

	// calculate number of requests sent
	requestsSent := *requestsPerSecond * *duration

	// calculate number of requests that succeeded
	requestsSucceeded := 0
	for _, i := range results {
		if i {
			requestsSucceeded++
		}
	}

	// calculate number of requests that failed
	requestsFailed := requestsSent - requestsSucceeded

	// calculate success rate
	successRate := float64(requestsSucceeded) / float64(requestsSent) * 100

	// calculate failure rate
	failureRate := float64(requestsFailed) / float64(requestsSent) * 100

	fmt.Println("\033[33mRequests sent:", requestsSent, "\033[0m")
	fmt.Println("\033[32mRequests succeeded:", requestsSucceeded, "\033[0m")
	fmt.Println("\033[31mRequests failed:", requestsFailed, "\033[0m")
	fmt.Println("\033[32mSuccess rate:", successRate, "%", "\033[0m")
	fmt.Println("\033[31mFailure rate:", failureRate, "%", "\033[0m")

	fmt.Println("------------------------- Thanks for using shex -----------------------")

}

func sendRequest(wg *sync.WaitGroup, mutex *sync.Mutex, results *[]bool) {
	mutex.Lock()

	//start := time.Now()
	_, err := http.Get(*targetDomain)
	//end := time.Now()
	mutex.Unlock()

	//fmt.Println("Request took", end.Sub(start))

	if err != nil {
		*results = append(*results, false)
		wg.Done()
	} else {
		*results = append(*results, true)
		wg.Done()
	}
}
