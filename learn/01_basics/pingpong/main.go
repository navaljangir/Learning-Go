// File: learn/07_concurrency_examples/15_ping_pong.go
package main

import (
    "fmt"
    "time"
)

func player(name string, receive <-chan string, send chan<- string) {
    for {
        ball := <-receive  // Wait for ball
        fmt.Printf("%s received: %s\n", name, ball)
        time.Sleep(2000 * time.Millisecond)
        send <- ball  // Hit it back
    }
}

func main() {
    ping := make(chan string)
    pong := make(chan string)

    go player("Ping", ping, pong)
    go player("Pong", pong, ping)

    // Start the game
    ping <- "🏓"

    // Let them play for 2 seconds
    time.Sleep(2 * time.Second)
    fmt.Println("Game over!")
}