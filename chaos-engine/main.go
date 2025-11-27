package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"os"
)

// -- CONFIGURATION --

const (
	// the URL of our python service
	// if running locally with 'go run', use localhost
	// if running in docker, this will be overridden by the ENV variables.

	TotalWorkers  = 5  // number of concurrent "users"
	TotalRequests = 200 // total messages to send
)

var TargetURL = "http://localhost:8000/chat"

// -- DATA MODEL --
// this matches the python "ChatRequest" model

type Payload struct {
	Prompt string `json:"prompt"`
	Persona string `json:"persona"`
}

func main() {
	fmt.Println("Starting Simulation Chaos Engine...")
	fmt.Println("Starting Simulation Chaos Engine...")
	if envURL := os.Getenv("TARGET_URL"); envURL != "" {
		TargetURL = envURL
	}
	fmt.Printf("Targeting: %s\n", TargetURL)
	fmt.Printf("Spawning %d concurrent workers...\n", TotalWorkers)

	// A WaitGroup to wait for all workers to finish before quitting
	// This is so we can see the final metrics before the program exits.
	var wg sync.WaitGroup
	// A channel is like a conveyor belt. We'll use it to pass messages between the workers and the main thread.
	// A channel is like a conveyor belt. We'll use it to pass messages between the workers and the main thread.
	jobs := make(chan int, TotalRequests)

	// 1. Start the workers
	// we spin up the "TotalWorkers" separate goroutines.
	// they all start listening to the "jobs" channel for tasks.
	for w := 1; w <= TotalWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	// 2. Fill the queue.

	// WE PUSH THE JOB IDs onto the conveyor belt.
	for j := 1; j <= TotalRequests; j++ {
		jobs <- j
	}
	close(jobs)

	// 3. Wait for all workers to finish.
	wg.Wait()
	fmt.Println("Simulation complete.")	
	
}

// This function runs inside each Goroutine (thread)
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	// List of "personas" to simulate different user types.
	personas := []string{"Angry User", "Hacker", "Curious Student", "Developer"}

	for j := range jobs {
		// Pick a random persona
		 persona := personas[rand.Intn(len(personas))]

		 
		 // create the JSON payload
		 payload := Payload{
			Prompt: fmt.Sprintf("Stress test message #%d from worker %d", j,id),
			Persona: persona,
		 }

		 jsonValue, _ := json.Marshal(payload)

		 // send the http request
		 start := time.Now()
		 resp, err := http.Post(TargetURL, "application/json", bytes.NewBuffer(jsonValue))
		 duration := time.Since(start)

		 if err != nil{
			log.Printf("Worker %d Error calling API: %v", id, err)
			continue
		 }

		 // log the result
		 // 200 OK means the python app handled it.
		 // 500 means we broke it or it timed out

		 fmt.Printf("Worker %d | Sent: %s | Status: %s | Time: %v\n",id,payload.Persona, resp.Status, duration)
		 resp.Body.Close()
	}
}
