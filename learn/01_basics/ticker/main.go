package main

import (
    "fmt"
    "time"
)

func main() {
    // Tick every 500ms
    ticker := time.NewTicker(500 * time.Millisecond)
	jobs := make(chan int,  5)
	close(jobs)
	fmt.Println("Jobs", jobs)
	  // ✅ Correct ways to print:
    fmt.Printf("Jobs (%%v): %v\n", jobs)        // Jobs (%v): 0xc0000240a0
    fmt.Printf("Jobs (%%p): %p\n", jobs)        // Jobs (%p): 0xc0000240a0
    fmt.Printf("Jobs type: %T\n", jobs)        // Jobs type: chan int
    fmt.Println("Jobs:", jobs)                 // Jobs: 0xc0000240a0
    
    // done := make(chan bool)

    // go func() {
    //     for {
    //         select {
    //         case <-done:
    //             return
    //         case t := <-ticker.C:
    //             fmt.Println("Tick at", t.Format("15:04:05.000"))
    //         }
    //     }
    // }()

    // // Run for 2.5 seconds
    // time.Sleep(2500 * time.Millisecond)
	fmt.Printf("%v\n", <-ticker.C)
    ticker.Stop()
    // done <- true
    fmt.Println("Ticker stopped")
}