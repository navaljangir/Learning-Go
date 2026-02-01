// package main

// import "fmt"

// func main() {
//     // UNBUFFERED - send blocks until receive
//     unbuffered := make(chan int)

//     go func() {
//         fmt.Println("Sending to unbuffered...")
//         unbuffered <- 1  // BLOCKS here until main receives
//         fmt.Println("Sent to unbuffered!")
//     }()

//     fmt.Println("Receiving from unbuffered:", <-unbuffered)

//     // BUFFERED - can hold values without receiver
//     buffered := make(chan int, 5)  // Can hold 3 values

//     buffered <- 1  // Doesn't block
//     buffered <- 2  // Doesn't block
//     fmt.Println("Buffered:", <-buffered)
//     buffered <- 3  // Doesn't block
//     buffered <- 4  // Would BLOCK - buffer full!

//     fmt.Println("Buffered:", <-buffered, <-buffered, <-buffered )
// }


// File: learn/01_basics/buffer/pipeline.go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Stage 1: Unbuffered - tight synchronization
    unbuffered := make(chan int)

    // Stage 2: Buffered - allows some slack
    buffered := make(chan int, 5)

    // Producer -> sends to unbuffered (must wait for consumer)
    go func() {
        for i := 1; i <= 5; i++ {
            fmt.Printf("Producer: sending %d to unbuffered...\n", i)
            unbuffered <- i  // Blocks until middle stage receives
            fmt.Printf("Producer: sent %d!\n", i)
        }
        close(unbuffered)
    }()

    // Middle stage: unbuffered -> buffered
    go func() {
        for val := range unbuffered {
            doubled := val * 2
            fmt.Printf("  Middle: putting %d into buffer\n", doubled)
            buffered <- doubled  // Won't block until buffer is full
        }
        close(buffered)
    }()

    // Consumer: reads from buffered (can be slow)
    time.Sleep(100 * time.Millisecond)  // Simulate slow start
    for result := range buffered {
        fmt.Printf("    Consumer: got %d\n", result)
        time.Sleep(50 * time.Millisecond)  // Simulate slow processing
    }
}