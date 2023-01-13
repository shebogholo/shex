package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"os"
	"sort"
	"strconv"
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

type responseItem struct {
	status      int
	latency     float64
	connectTime float64
}

// transport for http client
var transport = &http.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 100,
}

var client = &http.Client{
	Timeout:   time.Second * 90,
	Transport: transport,
}

func main() {
	fmt.Println("\n\n--------------------------- Welcome to shex ---------------------------")
	flag.Parse()

	if *url == "" {
		fmt.Println("Provide a target domain url")
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

	responseItems := make([]responseItem, 0)

	var wg sync.WaitGroup

	// record starting time
	start := time.Now()

	// loop through the number of seconds
	for i := 0; i < *duration; i++ {
		// loop for the number of requests per second
		for j := 0; j < *rps; j++ {
			wg.Add(1)
			go sendRequest(&wg, &sync.Mutex{}, &results, &responseItems)
		}
		time.Sleep(time.Second * 1)
	}

	// time elapsed since start of sending requests
	finished := time.Now()
	timeUsed := finished.Sub(start)

	// format time used to send requests
	if int(timeUsed.Minutes()) > 0 {
		minutes := int(timeUsed.Minutes())
		seconds := int(timeUsed.Seconds()) % 60
		fmt.Printf("\nRequests sent in %d minute(s) and %d seconds, waiting for responses...\n\n", minutes, seconds)
	} else {
		fmt.Printf("\nRequests sent in %d seconds, waiting for responses...\n\n", int(timeUsed.Seconds()))
	}

	wg.Wait()

	// record finishing time
	end := time.Now()
	// calculate time spent
	elapsed := end.Sub(start)

	// format elapsed time
	if int(elapsed.Minutes()) > 0 {
		minutes := int(elapsed.Minutes())
		seconds := int(elapsed.Seconds()) % 60
		fmt.Printf("\nLoad Testing took %d minute(s) and %d seconds\n", minutes, seconds)
	} else {
		fmt.Printf("\nLoad Testing took %d seconds\n", int(elapsed.Seconds()))
	}

	// calculate metrics
	metrics(&results)
	advancedMetrics(responseItems)
	fmt.Println("------------------------- Thanks for using shex -----------------------")

}

func sendRequest(wg *sync.WaitGroup, mutex *sync.Mutex, results *[]bool, responseItems *[]responseItem) {
	req, err := http.NewRequest(http.MethodGet, *url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	mutex.Lock()
	start := time.Now()
	// define connect time variable
	var connectTime float64

	// client trace
	trace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			connectTime = time.Since(start).Seconds()
		},
	}

	clientTraceCtx := httptrace.WithClientTrace(req.Context(), trace)
	req = req.WithContext(clientTraceCtx)
	resp, err := client.Do(req)

	//resp, err := http.Get(*url)
	end := time.Now()
	mutex.Unlock()

	elapsed := end.Sub(start)

	if err != nil {
		*results = append(*results, false)
		wg.Done()
	} else {
		*results = append(*results, true)

		// create response item from response
		respItem := responseItem{resp.StatusCode, elapsed.Seconds(), connectTime}

		// append resp to response items
		*responseItems = append(*responseItems, respItem)

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

func advancedMetrics(responseItems []responseItem) {
	numOfResponses := len(responseItems)

	// sort by latency
	sort.Slice(responseItems, func(i, j int) bool { return responseItems[i].latency < responseItems[j].latency })

	fastest := responseItems[0].latency
	slowest := responseItems[numOfResponses-1].latency

	fmt.Println(color.green, "Fastest request elapsed time: ", fastest, "seconds", color.reset)
	fmt.Println(color.red, "Slowest request elapsed time: ", slowest, "seconds", color.reset)

	// sort by connect time
	sort.Slice(responseItems, func(i, j int) bool { return responseItems[i].connectTime < responseItems[j].connectTime })

	fastestConnect := responseItems[0].connectTime
	slowestConnect := responseItems[numOfResponses-1].connectTime

	fmt.Println(color.green, "Fastest request connect time: ", fastestConnect, "seconds", color.reset)
	fmt.Println(color.red, "Slowest request connect time: ", slowestConnect, "seconds", color.reset)

	// save responseItems to csv file
	saveToCSV(responseItems)
}

func saveToCSV(responseItems []responseItem) {
	// create csv file
	file, err := os.Create("./results.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// create csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write headers
	writer.Write([]string{"Status", "Latency", "Connect"})

	// write data
	for _, item := range responseItems {
		writer.Write([]string{strconv.Itoa(item.status), strconv.FormatFloat(item.latency, 'f', 6, 64), strconv.FormatFloat(item.connectTime, 'f', 6, 64)})
	}

	fmt.Println("Results saved to ./results.csv")
}
