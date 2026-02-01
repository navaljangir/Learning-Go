// File: learn/07_concurrency_examples/13_worker_pool.go
package main

import (
    "fmt"
    "time"
)

type Job struct {
    ID      int
    Payload string
}

type Result struct {
    JobID    int
    Output   string
    WorkerID int
    Duration time.Duration
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        start := time.Now()

        // Simulate work
        time.Sleep(time.Duration(100+job.ID*10) * time.Millisecond)
        output := fmt.Sprintf("Processed: %s", job.Payload)

        results <- Result{
            JobID:    job.ID,
            Output:   output,
            WorkerID: id,
            Duration: time.Since(start),
        }
    }
}

func main() {
    numWorkers := 3
    numJobs := 10

    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)

    // Start workers
    fmt.Printf("Starting %d workers...\n", numWorkers)
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    fmt.Printf("Sending %d jobs...\n\n", numJobs)
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j, Payload: fmt.Sprintf("Task-%d", j)}
    }
    close(jobs)

    // Collect results
    for i := 0; i < numJobs; i++ {
        r := <-results
        fmt.Printf("Job %2d → Worker %d completed in %v\n", r.JobID, r.WorkerID, r.Duration)
    }
}