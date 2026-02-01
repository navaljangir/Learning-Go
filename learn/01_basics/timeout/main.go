package main

import (
    "fmt"
    "time"
)

func slowOperation() <-chan string {
    ch := make(chan string)
    go func() {
        time.Sleep(3 * time.Second)  // Simulates slow work
        ch <- "Operation complete!"
    }()
    return ch
}

func main() {
    fmt.Println("Starting operation...")

    select {
    case result := <-slowOperation():
        fmt.Println("Success:", result)
    case <-time.After( 5* time.Second):
        fmt.Println("ERROR: Operation timed out after 5 seconds!")
    }

    fmt.Println("Moving on...")
}