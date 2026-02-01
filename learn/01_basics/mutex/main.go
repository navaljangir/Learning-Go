package main

import (
    "fmt"
    "sync"
)

func main() {
    // WITHOUT MUTEX - WRONG!
    counter := 0
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++  // DATA RACE!
        }()
    }
    wg.Wait()
    fmt.Println("Without mutex:", counter)  // Wrong! (less than 1000)

    // WITH MUTEX - CORRECT!
    counter = 0
    var mu sync.Mutex

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()    // Only one goroutine can enter
            counter++
            mu.Unlock()  // Let others in
        }()
    }
    wg.Wait()
    fmt.Println("With mutex:", counter)  // Always 1000!
}