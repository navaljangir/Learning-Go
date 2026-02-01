package main

import (
    "context"
    "fmt"
    "time"
)

func doWork(ctx context.Context, name string) {
	time.Sleep(5 * time.Second) // Simulate initial work
}

func main() {
    // Create context with 3.5 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	fmt.Printf("Main: starting work with 3.5s timeout...\n")
    cancel()  // Always call cancel!
	doWork(ctx , "Worker")
	time.Sleep(1 * time.Second)
    fmt.Println()

    // doWork(ctx, "Worker")

    fmt.Println()
    fmt.Println("Main done")
}