# Go Concurrency: Goroutines, Channels & Context

Go's concurrency model is one of its biggest strengths. This guide covers everything with **real examples**, **actual outputs**, and **Node.js comparisons**.

---

## Key Definitions (Read This First!)

### What is Concurrency?
**Concurrency** = Dealing with multiple things at once (structure)
**Parallelism** = Doing multiple things at once (execution)

```
CONCURRENCY (managing multiple tasks):
┌─────────────────────────────────────────┐
│  Task A ─────►  Task B ─────►  Task A   │  One cook, multiple dishes
│         switch       switch             │  Switches between tasks
└─────────────────────────────────────────┘

PARALLELISM (actually simultaneous):
┌─────────────────────────────────────────┐
│  Task A ─────────────────────────────►  │  Multiple cooks
│  Task B ─────────────────────────────►  │  Working at same time
│  Task C ─────────────────────────────►  │
└─────────────────────────────────────────┘
```

Go gives you **BOTH** - concurrent code that runs in parallel on multiple CPU cores.

---

### Goroutine
**Definition:** A lightweight thread managed by Go runtime, not the OS.

| Feature | OS Thread | Goroutine |
|---------|-----------|-----------|
| Memory | ~1-8 MB stack | ~2 KB stack |
| Creation time | Slow (syscall) | Fast (Go runtime) |
| Switching | Expensive | Cheap |
| Max count | ~10,000 | ~1,000,000+ |

```go
// Create a goroutine - just add "go" keyword
go doSomething()
```

**Node.js equivalent:** `Promise`, `async/await` (but single-threaded)

---

### Channel
**Definition:** A typed pipe for goroutines to communicate and synchronize.

Think of it like a **conveyor belt** between workers:
- One worker puts items on
- Another worker takes items off
- They don't need to know about each other

```go
ch := make(chan string)  // Create pipe for strings
ch <- "hello"            // Put item on belt (send)
msg := <-ch              // Take item off belt (receive)
```

**Node.js equivalent:** `EventEmitter`, but type-safe and blocking

---

### Context
**Definition:** Carries deadlines, cancellation signals, and request-scoped values across API boundaries.

Think of it like a **cancellation token** that flows through your code:
- Parent cancelled? → All children get notified
- Timeout reached? → Everything stops gracefully

```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()  // Always clean up!
```

**Node.js equivalent:** `AbortController` / `AbortSignal`

---

### WaitGroup
**Definition:** A counter that blocks until it reaches zero. Used to wait for goroutines.

Think of it like a **checkout counter**:
- Each goroutine says "I'm starting" → `wg.Add(1)`
- Each goroutine says "I'm done" → `wg.Done()`
- Main waits until all done → `wg.Wait()`

```go
var wg sync.WaitGroup
wg.Add(1)      // +1 goroutine starting
go func() {
    defer wg.Done()  // -1 when done
}()
wg.Wait()      // Block until count = 0
```

**Node.js equivalent:** `Promise.all()`

---

### Mutex (Mutual Exclusion)
**Definition:** A lock that ensures only ONE goroutine can access shared data at a time.

Think of it like a **bathroom lock**:
- Lock the door → `mu.Lock()`
- Use the bathroom (modify data)
- Unlock the door → `mu.Unlock()`

```go
var mu sync.Mutex
mu.Lock()        // Only I can enter
counter++        // Safe to modify
mu.Unlock()      // Others can enter now
```

**Node.js equivalent:** Not needed (single-threaded), but similar to database transactions

---

### Select
**Definition:** Like a `switch` statement for channels. Waits for multiple channel operations.

```go
select {
case msg := <-ch1:     // If ch1 has data
case ch2 <- value:     // If ch2 can receive
case <-time.After(5s): // If 5 seconds pass
default:               // If nothing ready (non-blocking)
}
```

**Node.js equivalent:** `Promise.race()`

---

### Buffered vs Unbuffered Channels

| Type | Syntax | Behavior |
|------|--------|----------|
| **Unbuffered** | `make(chan int)` | Send blocks until receive (synchronous) |
| **Buffered** | `make(chan int, 5)` | Send blocks only when buffer full (async-ish) |

```
Unbuffered = Phone call (both parties must be present)
Buffered = Voicemail (leave message, pick up later)
```

---

## Table of Contents
1. [Node.js vs Go Architecture (Deep Dive)](#nodejs-vs-go-architecture)
2. [Understanding Concurrency (Visual)](#understanding-concurrency)
3. [Goroutines (Go's "Threads")](#goroutines)
3. [Channels (Communication)](#channels)
4. [Select Statement](#select-statement)
5. [Context (Cancellation & Timeouts)](#context)
6. [sync Package (WaitGroup, Mutex, etc.)](#sync-package)
7. [Real API Examples (Copy & Run!)](#real-api-examples)
8. [Common Patterns](#common-patterns)
9. [Common Mistakes](#common-mistakes)

---

## Node.js vs Go Architecture

### The Full Picture (What Tutorials Don't Tell You)

**Common misconception:** "Node.js is single-threaded"
**Reality:** Your JavaScript CODE is single-threaded, but I/O uses threads!

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              NODE.JS ARCHITECTURE                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   YOUR JAVASCRIPT CODE ──► SINGLE THREAD (V8 + Event Loop)                  │
│                                    │                                         │
│                                    │ Offloads I/O                           │
│                                    ▼                                         │
│   ┌────────────────────────────────────────────────────────────┐            │
│   │                         LIBUV                               │            │
│   │  ┌──────────────────┐    ┌──────────────────────────────┐  │            │
│   │  │   Thread Pool    │    │   OS Async (epoll/kqueue)    │  │            │
│   │  │   (4 threads)    │    │   (NO threads needed)        │  │            │
│   │  │                  │    │                              │  │            │
│   │  │  • File system   │    │  • Network I/O (HTTP, TCP)   │  │            │
│   │  │  • DNS lookup    │    │  • Sockets                   │  │            │
│   │  │  • Compression   │    │  • Timers                    │  │            │
│   │  │  • Crypto        │    │                              │  │            │
│   │  └──────────────────┘    └──────────────────────────────┘  │            │
│   └────────────────────────────────────────────────────────────┘            │
│                                    │                                         │
│                                    │ Callbacks                              │
│                                    ▼                                         │
│   BACK TO SINGLE THREAD ──► Execute callback in event loop                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### What's ACTUALLY Single-Threaded in Node.js?

| Component | Threads | Notes |
|-----------|---------|-------|
| **Your JS code** | 1 | Always single-threaded |
| **Event loop** | 1 | Processes callbacks one at a time |
| **V8 engine** | 1 | Executes JavaScript |
| **libuv thread pool** | 4 (default) | For blocking I/O (files, DNS, crypto) |
| **Network I/O** | 0 | Uses OS async primitives (no threads!) |

### Go Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              GO ARCHITECTURE                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   YOUR GO CODE ──► MULTIPLE GOROUTINES (truly parallel)                     │
│                                                                              │
│   ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│   │ Goroutine 1 │ │ Goroutine 2 │ │ Goroutine 3 │ │ Goroutine N │          │
│   │  (Request 1)│ │  (Request 2)│ │  (Request 3)│ │  (Request N)│          │
│   │             │ │             │ │             │ │             │          │
│   │ CPU work    │ │ CPU work    │ │ CPU work    │ │ CPU work    │          │
│   │ runs HERE   │ │ runs HERE   │ │ runs HERE   │ │ runs HERE   │          │
│   └──────┬──────┘ └──────┬──────┘ └──────┬──────┘ └──────┬──────┘          │
│          │               │               │               │                  │
│          └───────────────┼───────────────┼───────────────┘                  │
│                          │               │                                   │
│                    ┌─────▼───────────────▼─────┐                            │
│                    │      Go Scheduler          │                            │
│                    │  Maps goroutines to OS     │                            │
│                    │  threads (M:N scheduling)  │                            │
│                    └─────────────┬──────────────┘                            │
│                                  │                                           │
│                    ┌─────────────▼──────────────┐                            │
│                    │     OS Threads (CPUs)      │                            │
│                    │  [CPU 1] [CPU 2] [CPU 3]   │                            │
│                    └────────────────────────────┘                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### The KEY Difference: CPU-Bound Work

```javascript
// Node.js - CPU-bound work BLOCKS the event loop
app.get('/compute', (req, res) => {
    // While this loop runs, NO other request can be processed!
    // All callbacks wait in queue
    for (let i = 0; i < 1000000000; i++) {
        // Blocking the single JS thread
    }
    res.send('done');
});

// 100 concurrent requests? They ALL wait for each computation to finish
// Request 2 waits for Request 1's loop
// Request 3 waits for Request 2's loop
// etc.
```

```go
// Go - CPU-bound work runs in PARALLEL
func computeHandler(c *gin.Context) {
    // This runs in ITS OWN goroutine
    for i := 0; i < 1000000000; i++ {
        // Running on one CPU core
    }
    c.JSON(200, "done")
}

// 100 concurrent requests? Each gets its own goroutine!
// Request 1 runs on CPU 1
// Request 2 runs on CPU 2  (at the SAME time!)
// Request 3 runs on CPU 3  (at the SAME time!)
// etc.
```

### Summary Comparison

| Aspect | Node.js | Go |
|--------|---------|-----|
| **JS/Go code execution** | Single thread | Multiple goroutines (parallel) |
| **I/O operations** | libuv (threads + OS async) | Goroutines + OS async |
| **CPU-bound work** | **BLOCKS event loop** | Runs parallel on multiple cores |
| **100 CPU-heavy requests** | Sequential (slow!) | Parallel (fast!) |
| **100 I/O requests** | Fast (async) | Fast (async) |

### When Does It Matter?

| Workload | Node.js | Go | Winner |
|----------|---------|-----|--------|
| API that calls databases | Fast | Fast | Tie |
| API that calls other APIs | Fast | Fast | Tie |
| Image processing | **Slow (blocks)** | Fast | Go |
| JSON parsing (large) | **Slow (blocks)** | Fast | Go |
| Compression | **Slow (blocks)** | Fast | Go |
| Math-heavy operations | **Slow (blocks)** | Fast | Go |
| Simple CRUD | Fast | Fast | Tie |

**Rule of thumb:**
- I/O heavy (web APIs, databases) → Both are fine
- CPU heavy (processing, parsing, computation) → Go wins

---

## Understanding Concurrency

### Node.js vs Go: The Core Difference

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          NODE.JS (Single Thread)                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   ┌──────────┐                                                          │
│   │  Event   │  ──►  Task 1  ──►  Task 2  ──►  Task 3  ──►  ...        │
│   │   Loop   │       (async)      (async)      (async)                  │
│   └──────────┘                                                          │
│                                                                          │
│   One thread juggles everything. Async I/O, but CPU work blocks.        │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│                          GO (Multiple Goroutines)                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐                   │
│   │ Goroutine 1 │   │ Goroutine 2 │   │ Goroutine 3 │   ...            │
│   │   Task 1    │   │   Task 2    │   │   Task 3    │                   │
│   └─────────────┘   └─────────────┘   └─────────────┘                   │
│         │                 │                 │                            │
│         └────────────────┼─────────────────┘                            │
│                          │                                               │
│                    ┌─────▼─────┐                                        │
│                    │  Go       │  Distributes goroutines across         │
│                    │  Scheduler│  multiple CPU cores                    │
│                    └───────────┘                                        │
│                                                                          │
│   True parallelism. Multiple things running at the SAME time.           │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Quick Comparison Table

| Concept | Node.js | Go |
|---------|---------|-----|
| Async model | Single-threaded event loop | Multi-threaded goroutines |
| Starting async work | `async/await`, Promises | `go` keyword |
| Communication | Callbacks, Promises | Channels |
| Cancellation | AbortController | Context |
| Parallel execution | `Promise.all()` | Multiple goroutines + WaitGroup |
| Shared state protection | Not needed (single thread) | Mutex, Channels |
| Memory per "thread" | N/A | ~2KB per goroutine |
| Max concurrent | Limited by event loop | Millions of goroutines |

---

## Goroutines

### What is a Goroutine?

A goroutine is a **lightweight thread** managed by Go runtime. Cost: ~2KB of stack memory.

### Example 1: Basic Goroutine

```go
// File: learn/07_concurrency_examples/01_basic_goroutine.go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 1; i <= 3; i++ {
        fmt.Printf("[%s] Hello #%d\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    fmt.Println("=== Starting ===")

    go sayHello("Goroutine-1")  // Runs in background
    go sayHello("Goroutine-2")  // Runs in background

    sayHello("Main")  // Runs in main thread

    fmt.Println("=== Done ===")
}
```

**Output:** (order varies because they run concurrently!)
```
=== Starting ===
[Main] Hello #1
[Goroutine-1] Hello #1
[Goroutine-2] Hello #1
[Goroutine-2] Hello #2
[Main] Hello #2
[Goroutine-1] Hello #2
[Goroutine-1] Hello #3
[Goroutine-2] Hello #3
[Main] Hello #3
=== Done ===
```

### Visual: What's Happening

```
Time ──────────────────────────────────────────────────────────►

Main:        [Start]──[Hello#1]──[Hello#2]──[Hello#3]──[Done]
                │
Goroutine-1:   └──►[Hello#1]──[Hello#2]──[Hello#3]
                │
Goroutine-2:   └──►[Hello#1]──[Hello#2]──[Hello#3]

All three run at the SAME TIME (parallel)!
```

### Node.js Comparison

```javascript
// Node.js - NOT truly parallel (single thread)
async function sayHello(name) {
    for (let i = 1; i <= 3; i++) {
        console.log(`[${name}] Hello #${i}`);
        await new Promise(r => setTimeout(r, 100));
    }
}

async function main() {
    console.log("=== Starting ===");

    // These DON'T run in parallel - they alternate on event loop
    Promise.all([
        sayHello("Promise-1"),
        sayHello("Promise-2"),
        sayHello("Main")
    ]);
}
```

```
Node.js Output: (interleaved, but NOT truly parallel)
[Promise-1] Hello #1
[Promise-2] Hello #1
[Main] Hello #1
... (alternates on single thread)
```

### Example 2: The "Main Exits Too Fast" Problem

```go
package main

import "fmt"

func main() {
    go func() {
        fmt.Println("This might NOT print!")  // Goroutine might not finish
    }()

    // main() exits immediately, killing all goroutines
    fmt.Println("Main done")
}
```

**Output:**
```
Main done
```
(Goroutine never prints because main exits!)

**Fix with WaitGroup:**
```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup

    wg.Add(1)  // "I'm starting 1 goroutine"
    go func() {
        defer wg.Done()  // "I'm done" (when function exits)
        fmt.Println("This WILL print!")
    }()

    wg.Wait()  // Block until wg counter = 0
    fmt.Println("Main done")
}
```

**Output:**
```
This WILL print!
Main done
```

---

## Channels

### What is a Channel?

A channel is a **pipe** for goroutines to communicate. One sends, another receives.

### Visual: How Channels Work

```
                    CHANNEL (pipe)
                    ┌──────────┐
  Goroutine A ────► │   data   │ ────► Goroutine B
  (sender)          └──────────┘       (receiver)
                     ch <- val          val := <-ch
```

### Example 3: Basic Channel

```go
// File: learn/07_concurrency_examples/02_basic_channel.go
package main

import "fmt"

func main() {
    // Create a channel that carries strings
    messages := make(chan string)

    // Start goroutine that sends a message
    go func() {
        fmt.Println("Goroutine: Sending message...")
        messages <- "Hello from goroutine!"  // Send to channel
        fmt.Println("Goroutine: Message sent!")
    }()

    fmt.Println("Main: Waiting for message...")
    msg := <-messages  // Receive from channel (BLOCKS until data arrives)
    fmt.Println("Main: Received:", msg)
}
```

**Output:**
```
Main: Waiting for message...
Goroutine: Sending message...
Goroutine: Message sent!
Main: Received: Hello from goroutine!
```

### Visual: Execution Flow

```
Time ──────────────────────────────────────────────────────────►

Main:       [Create ch]──[Wait for msg]─────────────[Receive: "Hello"]──[Print]
                │              │ (blocked)              ▲
                │              │                        │
Goroutine:      └──►[Print]──[Send "Hello"]────────────┘
                              ch <- "Hello"        Main unblocks!
```

### Example 4: Unbuffered vs Buffered Channels

```go
package main

import "fmt"

func main() {
    // UNBUFFERED - send blocks until receive
    unbuffered := make(chan int)

    go func() {
        fmt.Println("Sending to unbuffered...")
        unbuffered <- 1  // BLOCKS here until main receives
        fmt.Println("Sent to unbuffered!")
    }()

    fmt.Println("Receiving from unbuffered:", <-unbuffered)

    // BUFFERED - can hold values without receiver
    buffered := make(chan int, 3)  // Can hold 3 values

    buffered <- 1  // Doesn't block
    buffered <- 2  // Doesn't block
    buffered <- 3  // Doesn't block
    // buffered <- 4  // Would BLOCK - buffer full!

    fmt.Println("Buffered:", <-buffered, <-buffered, <-buffered)
}
```

**Output:**
```
Sending to unbuffered...
Receiving from unbuffered: 1
Sent to unbuffered!
Buffered: 1 2 3
```

### Visual: Buffer Comparison

```
UNBUFFERED (make(chan int)):
┌────────┐          ┌────────┐
│ Sender │──────────│Receiver│   Must meet at the SAME TIME
└────────┘          └────────┘   Like a phone call - both parties needed


BUFFERED (make(chan int, 3)):
┌────────┐   ┌─┬─┬─┐   ┌────────┐
│ Sender │──►│1│2│3│──►│Receiver│   Buffer stores messages
└────────┘   └─┴─┴─┘   └────────┘   Like voicemail - leave message, pick up later
```

### What Happens With Different Buffer Sizes?

| Buffer Size | Behavior |
|-------------|----------|
| `make(chan int)` | Send blocks until receive (synchronous) |
| `make(chan int, 1)` | Can hold 1 value, then blocks |
| `make(chan int, 5)` | Can hold 5 values before blocking |
| `make(chan int, 100)` | Can hold 100 values - more "slack" in the system |

**Larger buffer = more decoupling between sender and receiver**, but uses more memory.

### Example 4b: Combining Buffered & Unbuffered (Pipeline)

This example shows both channel types working together - unbuffered for tight synchronization, buffered for allowing slack.

```go
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
```

**Output:**
```
Producer: sending 1 to unbuffered...
  Middle: putting 2 into buffer
Producer: sent 1!
Producer: sending 2 to unbuffered...
  Middle: putting 4 into buffer
Producer: sent 2!
Producer: sending 3 to unbuffered...
  Middle: putting 6 into buffer
Producer: sent 3!
Producer: sending 4 to unbuffered...
  Middle: putting 8 into buffer
Producer: sent 4!
Producer: sending 5 to unbuffered...
  Middle: putting 10 into buffer
Producer: sent 5!
    Consumer: got 2
    Consumer: got 4
    Consumer: got 6
    Consumer: got 8
    Consumer: got 10
```

### Visual: Pipeline with Mixed Channels

```
┌──────────┐     unbuffered      ┌────────────┐     buffered (5)      ┌──────────┐
│ Producer │────────────────────►│   Middle   │─────┬─┬─┬─┬─┬────────►│ Consumer │
│          │     (blocks each    │   Stage    │     │2│4│6│8│10       │  (slow)  │
│  1,2,3,  │      send until     │            │     └─┴─┴─┴─┴         │          │
│  4,5     │      received)      │  val * 2   │   (can queue up to 5) │          │
└──────────┘                     └────────────┘                       └──────────┘

Unbuffered: Producer MUST wait for Middle to receive each value
Buffered:   Middle can keep sending even if Consumer is slow (up to 5 values)
```

**Key insight:** The buffer acts as a "shock absorber" - the middle stage doesn't have to wait for the slow consumer, it can keep working until the buffer fills up.

### Example 5: Range Over Channel

```go
package main

import "fmt"

func main() {
    numbers := make(chan int)

    // Producer: sends numbers then closes
    go func() {
        for i := 1; i <= 5; i++ {
            fmt.Printf("Sending: %d\n", i)
            numbers <- i
        }
        close(numbers)  // IMPORTANT: close when done!
        fmt.Println("Channel closed")
    }()

    // Consumer: range automatically stops when closed
    fmt.Println("Receiving:")
    for num := range numbers {
        fmt.Printf("  Got: %d\n", num)
    }
    fmt.Println("Done receiving")
}
```

**Output:**
```
Receiving:
Sending: 1
  Got: 1
Sending: 2
  Got: 2
Sending: 3
  Got: 3
Sending: 4
  Got: 4
Sending: 5
  Got: 5
Channel closed
Done receiving
```

---

## Select Statement

### What is Select?

`select` waits on **multiple channels** - whichever is ready first wins.

### Node.js Comparison

```javascript
// Node.js - Promise.race()
const result = await Promise.race([
    fetch('https://api1.com/data'),
    fetch('https://api2.com/data'),
    new Promise((_, reject) =>
        setTimeout(() => reject(new Error('timeout')), 5000)
    )
]);
```

```go
// Go - select
select {
case data := <-api1Channel:
    fmt.Println("API 1 responded first:", data)
case data := <-api2Channel:
    fmt.Println("API 2 responded first:", data)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout!")
}
```

### Example 6: First Response Wins

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func fetchFromServer(name string, ch chan<- string) {
    // Simulate random response time (100-500ms)
    delay := time.Duration(100+rand.Intn(400)) * time.Millisecond
    time.Sleep(delay)
    ch <- fmt.Sprintf("%s responded in %v", name, delay)
}

func main() {
    rand.Seed(time.Now().UnixNano())

    server1 := make(chan string)
    server2 := make(chan string)

    go fetchFromServer("Server-1", server1)
    go fetchFromServer("Server-2", server2)

    // Wait for FIRST response only
    select {
    case msg := <-server1:
        fmt.Println("Winner:", msg)
    case msg := <-server2:
        fmt.Println("Winner:", msg)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout - no server responded!")
    }
}
```

**Output:** (varies each run)
```
Winner: Server-2 responded in 156ms
```
or
```
Winner: Server-1 responded in 203ms
```

### Example 7: Timeout Pattern

```go
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
    case <-time.After(2 * time.Second):
        fmt.Println("ERROR: Operation timed out after 2 seconds!")
    }

    fmt.Println("Moving on...")
}
```

**Output:**
```
Starting operation...
ERROR: Operation timed out after 2 seconds!
Moving on...
```

---

## Context

### What is Context?

Context handles **cancellation**, **timeouts**, and **passing values** across API boundaries.

### Visual: Context Propagation

```
HTTP Request arrives
       │
       ▼
┌──────────────────┐
│   Handler        │  ctx := c.Request.Context()
│   (timeout: 30s) │
└────────┬─────────┘
         │ ctx
         ▼
┌──────────────────┐
│   Service Layer  │  Uses same ctx
└────────┬─────────┘
         │ ctx
         ▼
┌──────────────────┐
│   Database       │  If ctx cancelled, stop query!
└──────────────────┘

If client disconnects → ctx cancelled → ALL operations stop
```

### Example 8: Context with Timeout

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func doWork(ctx context.Context, name string) {
    for i := 1; i <= 5; i++ {
        select {
        case <-ctx.Done():
            fmt.Printf("[%s] Cancelled! Reason: %v\n", name, ctx.Err())
            return
        default:
            fmt.Printf("[%s] Working... step %d\n", name, i)
            time.Sleep(500 * time.Millisecond)
        }
    }
    fmt.Printf("[%s] Completed!\n", name)
}

func main() {
    // Create context with 1.5 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
    defer cancel()  // Always call cancel!

    fmt.Println("Starting work with 1.5s timeout...")
    fmt.Println()

    doWork(ctx, "Worker")

    fmt.Println()
    fmt.Println("Main done")
}
```

**Output:**
```
Starting work with 1.5s timeout...

[Worker] Working... step 1
[Worker] Working... step 2
[Worker] Working... step 3
[Worker] Cancelled! Reason: context deadline exceeded

Main done
```

### Understanding Context Functions

#### `context.Background()`
The **root context** - empty, never cancelled. Starting point for all contexts.

```go
ctx := context.Background()  // Empty context, use as parent
```

#### `context.WithTimeout(parent, duration)`
Creates a context that **automatically cancels** after the duration.

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // Always call cancel to release resources!

// ctx.Done() will receive a value after 5 seconds
```

#### `context.WithCancel(parent)`
Creates a context you can **manually cancel**.

```go
ctx, cancel := context.WithCancel(context.Background())

// Later, when you want to stop:
cancel()  // ctx.Done() receives a value
```

#### `ctx.Done()` - The Cancellation Channel
Returns a channel that closes when context is cancelled/timed out.

```go
select {
case <-ctx.Done():
    // Context was cancelled or timed out
    fmt.Println("Stopped:", ctx.Err())
}
```

#### `ctx.Err()` - Why Was It Cancelled?
```go
ctx.Err()  // Returns:
           // - nil (not cancelled yet)
           // - context.Canceled (manually cancelled)
           // - context.DeadlineExceeded (timeout)
```

### Context vs time.Sleep vs select+time.After

| Approach | Can Cancel? | Propagates? | Use Case |
|----------|-------------|-------------|----------|
| `time.Sleep(5s)` | **No** | No | Simple delay, nothing else |
| `select + time.After` | **No** | No | Timeout for ONE operation |
| `context.WithTimeout` | **Yes** | **Yes** | Timeout across MULTIPLE operations |

#### time.Sleep - Just Waits

```go
time.Sleep(5 * time.Second)  // Blocks. Can't cancel. Can't escape.
```

#### select + time.After - Timeout ONE Operation

```go
select {
case result := <-doWork():
    fmt.Println("Got result")
case <-time.After(5 * time.Second):
    fmt.Println("Timeout")
}
// Problem: doWork() goroutine might still be running!
```

#### Context - Timeout + Signal to STOP

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := doWorkWithContext(ctx)  // Function checks ctx.Done()
if err != nil {
    fmt.Println("Failed or timed out")
}
// The function KNOWS to stop when ctx is cancelled
```

### Visual: Why Context is Different

```
select + time.After:
┌─────────────────────────────────────────────────────────────────┐
│  Main:     select { case <-doWork()... case <-time.After(5s) }  │
│                                                    ↓            │
│                                              Timeout! Move on   │
│                                                                 │
│  doWork:   Still running... doesn't know we gave up!            │
│            └── Wasting resources                                │
└─────────────────────────────────────────────────────────────────┘

Context:
┌─────────────────────────────────────────────────────────────────┐
│  Main:     ctx, cancel := context.WithTimeout(..., 5s)          │
│            doWorkWithContext(ctx)                               │
│                     │                                           │
│                     ↓  (after 5s, ctx.Done() fires)             │
│                                                                 │
│  doWork:   select { case <-ctx.Done(): return }                 │
│            └── Knows to stop! Cleans up properly.               │
└─────────────────────────────────────────────────────────────────┘
```

### When to Use What

| Situation | Use This |
|-----------|----------|
| Simple wait | `time.Sleep()` |
| Race two channels | `select` |
| Timeout one channel op | `select + time.After` |
| Pass timeout down call chain | `context.WithTimeout` |
| Cancel multiple goroutines | `context.WithCancel` |
| HTTP handlers | Context (already provided) |
| Database queries | Context |

### Example 9: Manual Cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: Stopping (reason: %v)\n", id, ctx.Err())
            return
        default:
            fmt.Printf("Worker %d: Working...\n", id)
            time.Sleep(300 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    // Start 3 workers
    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }

    // Let them work for 1 second
    time.Sleep(1 * time.Second)

    fmt.Println("\n>>> Cancelling all workers! <<<\n")
    cancel()  // This cancels ALL workers using this context

    // Give workers time to print their stop message
    time.Sleep(100 * time.Millisecond)
    fmt.Println("All workers stopped")
}
```

**Output:**
```
Worker 1: Working...
Worker 2: Working...
Worker 3: Working...
Worker 1: Working...
Worker 3: Working...
Worker 2: Working...
Worker 2: Working...
Worker 1: Working...
Worker 3: Working...

>>> Cancelling all workers! <<<

Worker 1: Stopping (reason: context canceled)
Worker 2: Stopping (reason: context canceled)
Worker 3: Stopping (reason: context canceled)
All workers stopped
```

### Node.js Comparison: AbortController

```javascript
// Node.js - AbortController (similar concept)
const controller = new AbortController();

async function worker(signal, id) {
    while (!signal.aborted) {
        console.log(`Worker ${id}: Working...`);
        await new Promise(r => setTimeout(r, 300));
    }
    console.log(`Worker ${id}: Stopped`);
}

// Start workers
worker(controller.signal, 1);
worker(controller.signal, 2);
worker(controller.signal, 3);

// Cancel after 1 second
setTimeout(() => {
    console.log(">>> Cancelling all workers! <<<");
    controller.abort();
}, 1000);
```

---

## sync Package

### WaitGroup: Wait for Multiple Goroutines

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()  // Decrements counter when done

    fmt.Printf("Worker %d: Starting\n", id)
    time.Sleep(time.Duration(id*100) * time.Millisecond)
    fmt.Printf("Worker %d: Done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)       // Increment counter
        go worker(i, &wg)
    }

    fmt.Println("Main: Waiting for all workers...")
    wg.Wait()  // Blocks until counter = 0
    fmt.Println("Main: All workers completed!")
}
```

**Output:**
```
Main: Waiting for all workers...
Worker 5: Starting
Worker 1: Starting
Worker 2: Starting
Worker 3: Starting
Worker 4: Starting
Worker 1: Done
Worker 2: Done
Worker 3: Done
Worker 4: Done
Worker 5: Done
Main: All workers completed!
```

### Node.js Comparison: Promise.all

```javascript
// Node.js equivalent
async function worker(id) {
    console.log(`Worker ${id}: Starting`);
    await new Promise(r => setTimeout(r, id * 100));
    console.log(`Worker ${id}: Done`);
}

async function main() {
    console.log("Main: Waiting for all workers...");
    await Promise.all([
        worker(1),
        worker(2),
        worker(3),
        worker(4),
        worker(5)
    ]);
    console.log("Main: All workers completed!");
}
```

### Mutex: Protect Shared Data

```go
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
```

**Output:**
```
Without mutex: 947
With mutex: 1000
```

---

## Real API Examples

### Example 10: Fetch Multiple APIs in Parallel (REAL CODE!)

```go
// File: learn/07_concurrency_examples/10_parallel_api.go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"
)

// Result holds API response
type Result struct {
    API      string
    Data     interface{}
    Duration time.Duration
    Error    error
}

// Fetch from a single API
func fetchAPI(url string, name string, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()

    start := time.Now()

    resp, err := http.Get(url)
    if err != nil {
        results <- Result{API: name, Error: err, Duration: time.Since(start)}
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        results <- Result{API: name, Error: err, Duration: time.Since(start)}
        return
    }

    var data interface{}
    json.Unmarshal(body, &data)

    results <- Result{
        API:      name,
        Data:     data,
        Duration: time.Since(start),
    }
}

func main() {
    // Real public APIs
    apis := map[string]string{
        "JSONPlaceholder": "https://jsonplaceholder.typicode.com/todos/1",
        "Cat Fact":        "https://catfact.ninja/fact",
        "Random User":     "https://randomuser.me/api/?results=1",
        "IP Info":         "https://ipapi.co/json/",
    }

    results := make(chan Result, len(apis))
    var wg sync.WaitGroup

    fmt.Println("Fetching from multiple APIs in parallel...")
    fmt.Println("==========================================")
    start := time.Now()

    // Launch all requests in parallel
    for name, url := range apis {
        wg.Add(1)
        go fetchAPI(url, name, results, &wg)
    }

    // Close channel when all done
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results as they arrive
    for result := range results {
        if result.Error != nil {
            fmt.Printf("❌ %s: Error - %v\n", result.API, result.Error)
        } else {
            fmt.Printf("✓ %s: Success (took %v)\n", result.API, result.Duration)
        }
    }

    fmt.Println("==========================================")
    fmt.Printf("Total time: %v\n", time.Since(start))
}
```

**Output:**
```
Fetching from multiple APIs in parallel...
==========================================
✓ JSONPlaceholder: Success (took 156.234ms)
✓ Cat Fact: Success (took 203.456ms)
✓ IP Info: Success (took 312.789ms)
✓ Random User: Success (took 523.123ms)
==========================================
Total time: 524.567ms
```

**Compare to Sequential (Node.js style):**
```
If done one-by-one: ~1195ms (156+203+312+524)
Done in parallel:   ~524ms  (only as slow as slowest)
Speedup:            ~2.3x faster!
```

### Example 11: Fetch with Timeout (Context)

```go
// File: learn/07_concurrency_examples/11_fetch_with_timeout.go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
)

func fetchWithTimeout(ctx context.Context, url string) (string, error) {
    // Create request with context
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

func main() {
    // Test 1: Fast API with reasonable timeout
    fmt.Println("Test 1: Fetching fast API with 5s timeout...")
    ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel1()

    start := time.Now()
    result, err := fetchWithTimeout(ctx1, "https://jsonplaceholder.typicode.com/todos/1")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Success! (took %v)\n", time.Since(start))
        fmt.Printf("Data: %.50s...\n", result)
    }

    fmt.Println()

    // Test 2: API with very short timeout (will fail)
    fmt.Println("Test 2: Fetching API with 1ms timeout (will fail)...")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Millisecond)
    defer cancel2()

    start = time.Now()
    result, err = fetchWithTimeout(ctx2, "https://jsonplaceholder.typicode.com/todos/1")
    if err != nil {
        fmt.Printf("Error (expected): %v\n", err)
    } else {
        fmt.Printf("Success: %.50s...\n", result)
    }
}
```

**Output:**
```
Test 1: Fetching fast API with 5s timeout...
Success! (took 134.567ms)
Data: {
  "userId": 1,
  "id": 1,
  "title": "delectu...

Test 2: Fetching API with 1ms timeout (will fail)...
Error (expected): Get "https://jsonplaceholder.typicode.com/todos/1": context deadline exceeded
```

### Example 12: Rate Limited API Calls

```go
// File: learn/07_concurrency_examples/12_rate_limited.go
package main

import (
    "fmt"
    "net/http"
    "sync"
    "time"
)

func main() {
    // URLs to fetch
    urls := []string{
        "https://jsonplaceholder.typicode.com/todos/1",
        "https://jsonplaceholder.typicode.com/todos/2",
        "https://jsonplaceholder.typicode.com/todos/3",
        "https://jsonplaceholder.typicode.com/todos/4",
        "https://jsonplaceholder.typicode.com/todos/5",
        "https://jsonplaceholder.typicode.com/todos/6",
        "https://jsonplaceholder.typicode.com/todos/7",
        "https://jsonplaceholder.typicode.com/todos/8",
        "https://jsonplaceholder.typicode.com/todos/9",
        "https://jsonplaceholder.typicode.com/todos/10",
    }

    // Semaphore: limit to 3 concurrent requests
    semaphore := make(chan struct{}, 3)
    var wg sync.WaitGroup

    fmt.Println("Fetching 10 URLs with max 3 concurrent...")
    fmt.Println("==========================================")
    start := time.Now()

    for i, url := range urls {
        wg.Add(1)
        go func(i int, url string) {
            defer wg.Done()

            semaphore <- struct{}{}        // Acquire slot (blocks if 3 running)
            defer func() { <-semaphore }() // Release slot

            reqStart := time.Now()
            resp, err := http.Get(url)
            if err != nil {
                fmt.Printf("[%2d] Error: %v\n", i+1, err)
                return
            }
            resp.Body.Close()

            fmt.Printf("[%2d] Done in %v (status: %d)\n", i+1, time.Since(reqStart), resp.StatusCode)
        }(i, url)
    }

    wg.Wait()
    fmt.Println("==========================================")
    fmt.Printf("Total time: %v\n", time.Since(start))
}
```

**Output:**
```
Fetching 10 URLs with max 3 concurrent...
==========================================
[ 1] Done in 145ms (status: 200)
[ 3] Done in 152ms (status: 200)
[ 2] Done in 156ms (status: 200)
[ 4] Done in 148ms (status: 200)
[ 6] Done in 151ms (status: 200)
[ 5] Done in 155ms (status: 200)
[ 7] Done in 142ms (status: 200)
[ 9] Done in 149ms (status: 200)
[ 8] Done in 153ms (status: 200)
[10] Done in 147ms (status: 200)
==========================================
Total time: 612ms
```

### Example 13: Worker Pool Pattern

```go
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
```

**Output:**
```
Starting 3 workers...
Sending 10 jobs...

Job  1 → Worker 1 completed in 110ms
Job  2 → Worker 2 completed in 120ms
Job  3 → Worker 3 completed in 130ms
Job  4 → Worker 1 completed in 140ms
Job  5 → Worker 2 completed in 150ms
Job  6 → Worker 3 completed in 160ms
Job  7 → Worker 1 completed in 170ms
Job  8 → Worker 2 completed in 180ms
Job  9 → Worker 3 completed in 190ms
Job 10 → Worker 1 completed in 200ms
```

### Visual: Worker Pool

```
                    ┌─────────────┐
                    │  Job Queue  │
                    │ [1][2][3]...│
                    └──────┬──────┘
                           │
           ┌───────────────┼───────────────┐
           ▼               ▼               ▼
    ┌──────────┐    ┌──────────┐    ┌──────────┐
    │ Worker 1 │    │ Worker 2 │    │ Worker 3 │
    └────┬─────┘    └────┬─────┘    └────┬─────┘
         │               │               │
         └───────────────┼───────────────┘
                         ▼
                  ┌─────────────┐
                  │  Results    │
                  │  Channel    │
                  └─────────────┘
```

---

## Common Patterns

### Pattern 1: Fan-Out, Fan-In

```
               Fan-Out                    Fan-In
                  │                          │
    ┌─────────────┼─────────────┐           │
    ▼             ▼             ▼           │
┌───────┐   ┌───────┐   ┌───────┐          │
│Worker1│   │Worker2│   │Worker3│          │
└───┬───┘   └───┬───┘   └───┬───┘          │
    │           │           │              │
    └───────────┼───────────┘              │
                ▼                          ▼
         ┌─────────────┐           ┌─────────────┐
         │  Merge All  │           │   Single    │
         │   Results   │──────────►│   Output    │
         └─────────────┘           └─────────────┘
```

### Pattern 2: Pipeline

```
Input ──► Stage 1 ──► Stage 2 ──► Stage 3 ──► Output

Each stage runs in its own goroutine, connected by channels.
Data flows through like an assembly line.
```

---

## Common Mistakes

### Mistake 1: Forgetting Loop Variable

```go
// WRONG - all goroutines see final value
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // Prints: 3, 3, 3
    }()
}

// RIGHT - pass as parameter
for i := 0; i < 3; i++ {
    go func(n int) {
        fmt.Println(n)  // Prints: 0, 1, 2 (some order)
    }(i)
}
```

### Mistake 2: Sending to Closed Channel

```go
ch := make(chan int)
close(ch)
ch <- 1  // PANIC: send on closed channel
```

### Mistake 3: No WaitGroup = Goroutine Dies

```go
// WRONG
func main() {
    go func() {
        fmt.Println("This might not print!")
    }()
    // main exits, goroutine killed
}

// RIGHT
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println("This will print!")
    }()
    wg.Wait()
}
```

### Mistake 4: Race Condition

```go
// WRONG - data race
counter := 0
for i := 0; i < 1000; i++ {
    go func() {
        counter++  // Multiple goroutines modifying!
    }()
}
// counter will be WRONG (less than 1000)

// RIGHT - use mutex
var mu sync.Mutex
for i := 0; i < 1000; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
```

---

## Memory Leaks vs Resource Leaks

### Understanding the Difference

```
┌─────────────────────────────────────────────────────────────────┐
│  RESOURCE LEAK (broader category)                               │
│  ├── Memory Leak                                                │
│  ├── Goroutine Leak                                             │
│  ├── File Handle Leak                                           │
│  ├── Network Connection Leak                                    │
│  ├── Database Connection Leak                                   │
│  └── Lock/Mutex Leak                                            │
│                                                                 │
│  Resource leak = Any system resource not properly freed         │
│  Memory leak = Specifically MEMORY not freed                    │
└─────────────────────────────────────────────────────────────────┘
```

### Memory Leak

**Definition:** Memory is allocated but never freed, even though it's no longer needed.

```go
var globalSlice []int

func leakMemory() {
    for i := 0; i < 1000000; i++ {
        // Allocating memory but never releasing it
        globalSlice = append(globalSlice, i)
    }
    // globalSlice keeps growing, memory never freed
}

func main() {
    for {
        leakMemory()  // Memory usage keeps increasing
        time.Sleep(1 * time.Second)
    }
}
```

**Result:**
```
0s:   Memory usage: 10 MB
10s:  Memory usage: 100 MB
20s:  Memory usage: 200 MB
30s:  Memory usage: 300 MB
...   Program crashes: Out of Memory (OOM)
```

### Resource Leaks (Different Types)

#### 1. File Handle Leak

```go
func leakFileHandles() {
    for i := 0; i < 1000; i++ {
        f, _ := os.Open("data.txt")
        // FORGOT to f.Close()!
        // File handle never released
    }
}
```

**Result:**
```
Error: too many open files (OS limit reached)
Even if memory is fine!
```

**Fix:**
```go
func correctFileHandling() {
    f, err := os.Open("data.txt")
    if err != nil {
        return
    }
    defer f.Close()  // Always close!
    // ... use file ...
}
```

#### 2. Goroutine Leak

```go
func leakGoroutines() {
    ch := make(chan int)  // Unbuffered channel

    for i := 0; i < 1000; i++ {
        go func() {
            ch <- i  // BLOCKS FOREVER (no receiver)
                     // Goroutine never exits
        }()
    }
}
```

**Result:**
```
1000 goroutines stuck forever
Each goroutine has a stack (~2-8 KB)
Total memory leak: ~2-8 MB
But more importantly: 1000 zombie goroutines!
```

**Fix:**
```go
func correctGoroutineHandling() {
    ch := make(chan int, 1000)  // Buffered! Won't block

    for i := 0; i < 1000; i++ {
        go func(val int) {
            ch <- val  // Doesn't block, goroutine can exit
        }(i)
    }
}
```

#### 3. Network Connection Leak

```go
func leakConnections() {
    for i := 0; i < 100; i++ {
        conn, _ := net.Dial("tcp", "example.com:80")
        // FORGOT to conn.Close()!
    }
}
```

**Result:**
```
Error: cannot create more connections
OS connection limit reached
Server can't accept new clients
```

**Fix:**
```go
func correctConnectionHandling() {
    conn, err := net.Dial("tcp", "example.com:80")
    if err != nil {
        return
    }
    defer conn.Close()  // Always close!
    // ... use connection ...
}
```

### Goroutine Leak: Both Memory AND Resource Leak

**A goroutine leak is BOTH a resource leak AND a memory leak.**

```
┌─────────────────────────────────────────────────────────────────┐
│  Each goroutine consumes:                                       │
│                                                                 │
│  1. MEMORY (for stack)           → Memory leak                  │
│     - Minimum: 2 KB                                             │
│     - Grows as needed: up to 1 GB per goroutine!                │
│                                                                 │
│  2. SCHEDULER SLOT               → Resource leak                │
│     - OS thread resources                                       │
│     - Scheduling overhead                                       │
│                                                                 │
│  3. CHANNEL/LOCK REFERENCES      → Prevents GC                  │
│     - Keeps other objects alive                                 │
└─────────────────────────────────────────────────────────────────┘
```

**Example:**
```go
func main() {
    for i := 0; i < 1000000; i++ {
        ch := make(chan int)
        go func() {
            <-ch  // Blocks forever, goroutine never exits
        }()
    }
    // 1,000,000 goroutines leaked!
    // Memory: ~2GB minimum (2KB × 1M)
    // Resource: Scheduler has 1M goroutines to manage
}
```

### Visual Comparison

```
MEMORY LEAK:
┌────────────────────────────────────────┐
│  RAM Usage Over Time                   │
│  ▲                                     │
│  │              ╱╱╱╱╱╱                 │
│  │         ╱╱╱╱╱                       │
│  │    ╱╱╱╱╱                            │
│  │ ╱╱╱                                 │
│  └────────────────────► Time           │
│                                        │
│  Eventually: Out of Memory (OOM Kill)  │
└────────────────────────────────────────┘

FILE HANDLE LEAK:
┌────────────────────────────────────────┐
│  Open Files                            │
│  ▲                                     │
│  │ ┌──┬──┬──┬──┬──┬──┐                 │
│  │ │  │  │  │  │  │  │                 │
│  │ │  │  │  │  │  │  │                 │
│  0 └──┴──┴──┴──┴──┴──┘                 │
│     Limit: 1024 (ulimit -n)            │
│                                        │
│  Eventually: "too many open files"     │
└────────────────────────────────────────┘

GOROUTINE LEAK:
┌────────────────────────────────────────┐
│  Goroutines                            │
│  ▲                                     │
│  │ ⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙             │
│  │ ⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙               │
│  │ ⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙⚙                       │
│  0 ⚙⚙⚙⚙                               │
│     All stuck, never exit              │
│                                        │
│  Eventually: OOM + scheduler overhead  │
└────────────────────────────────────────┘
```

### Real-World Example: HTTP Server Goroutine Leak

```go
// BAD: Leaks goroutines
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ch := make(chan Result)

    go func() {
        result := slowWork()
        ch <- result  // What if nobody reads this?
    }()

    select {
    case result := <-ch:
        json.NewEncoder(w).Encode(result)
    case <-time.After(5 * time.Second):
        http.Error(w, "timeout", 504)
        return  // ⚠️ Goroutine still running, blocked on ch!
    }
}
```

**What happens:**
```
Request 1: Timeout → 1 goroutine leaked
Request 2: Timeout → 2 goroutines leaked
Request 3: Timeout → 3 goroutines leaked
...
Request 10000: Timeout → 10000 goroutines leaked
Server: OOM killed!
```

**Fix with buffered channel:**
```go
// GOOD: No leak
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ch := make(chan Result, 1)  // Buffered! Goroutine can exit

    go func() {
        result := slowWork()
        ch <- result  // Doesn't block even if nobody reads
    }()

    select {
    case result := <-ch:
        json.NewEncoder(w).Encode(result)
    case <-time.After(5 * time.Second):
        http.Error(w, "timeout", 504)
        // Goroutine can send to buffer and exit cleanly
    }
}
```

**Better fix with context:**
```go
// BEST: Context-aware
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    ch := make(chan Result, 1)

    go func() {
        result := slowWorkWithContext(ctx)
        ch <- result
    }()

    select {
    case result := <-ch:
        json.NewEncoder(w).Encode(result)
    case <-ctx.Done():
        http.Error(w, "timeout", 504)
        // Goroutine sees ctx.Done(), exits cleanly
    }
}
```

### The Ping-Pong Example: Is It a Leak?

```go
func main() {
    ping := make(chan string)
    pong := make(chan string)

    go player("Ping", ping, pong)  // Infinite loop
    go player("Pong", pong, ping)  // Infinite loop

    ping <- "🏓"
    time.Sleep(2 * time.Second)
    fmt.Println("Game over!")
}  // main() exits → program ends
```

**Is this a leak?**

```
In this specific program: NO, not really
  → Program exits after 2 seconds
  → OS cleans up everything (goroutines, memory, channels)

If this were in a long-running server: YES, leak!
  → Goroutines accumulate over time
  → Memory usage grows
  → Eventually crash
```

**Better version with graceful shutdown:**

```go
func player(name string, receive <-chan string, send chan<- string, done <-chan struct{}) {
    for {
        select {
        case ball := <-receive:
            fmt.Printf("%s received: %s\n", name, ball)
            time.Sleep(100 * time.Millisecond)
            send <- ball
        case <-done:
            fmt.Printf("%s stopped\n", name)
            return  // Exit goroutine cleanly
        }
    }
}

func main() {
    ping := make(chan string)
    pong := make(chan string)
    done := make(chan struct{})

    go player("Ping", ping, pong, done)
    go player("Pong", pong, ping, done)

    ping <- "🏓"
    time.Sleep(2 * time.Second)

    close(done)  // Signal both to stop
    time.Sleep(100 * time.Millisecond)  // Let them cleanup
    fmt.Println("Game over!")
}
```

### Summary: Leak Types

| Type | What Leaks | Impact | Fix |
|------|------------|--------|-----|
| **Memory Leak** | Heap objects, slices, maps | OOM crash | Release references, avoid global accumulation |
| **Goroutine Leak** | Goroutines (+ their memory) | OOM + scheduler overhead | Use context, buffered channels, proper cleanup |
| **File Handle Leak** | File descriptors | "too many open files" | Always `defer f.Close()` |
| **Connection Leak** | Network/DB connections | Connection pool exhaustion | Always `defer conn.Close()` |
| **Channel Leak** | Blocked goroutines on channels | Goroutine leak | Use buffered channels or timeouts |

### Detection Tips

**1. Check goroutine count:**
```go
import "runtime"

func printGoroutines() {
    fmt.Println("Active goroutines:", runtime.NumGoroutine())
}
```

**2. Use Go's race detector:**
```bash
go run -race yourprogram.go
```

**3. Profile memory:**
```go
import _ "net/http/pprof"

go http.ListenAndServe("localhost:6060", nil)
// Visit http://localhost:6060/debug/pprof/
```

**4. Check for goroutine leaks:**
```bash
# Get goroutine dump
curl http://localhost:6060/debug/pprof/goroutine?debug=2
```

---

## Quick Reference

### Goroutines
```go
go functionName()         // Start goroutine
go func() { }()          // Anonymous goroutine
```

### Channels
```go
ch := make(chan int)      // Unbuffered
ch := make(chan int, 10)  // Buffered (capacity 10)
ch <- value               // Send (blocks if full)
value := <-ch             // Receive (blocks if empty)
close(ch)                 // Close (sender only!)
```

### Select
```go
select {
case v := <-ch1:          // First ready wins
case ch2 <- value:        // Can also send
case <-time.After(5*time.Second):  // Timeout
default:                  // Non-blocking
}
```

### Context
```go
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
ctx, cancel := context.WithCancel(ctx)
defer cancel()            // Always call!
<-ctx.Done()              // Blocks until cancelled
ctx.Err()                 // context.Canceled or context.DeadlineExceeded
```

### sync Package
```go
var wg sync.WaitGroup
wg.Add(1)                 // Before goroutine
wg.Done()                 // Inside goroutine (defer)
wg.Wait()                 // Wait for all

var mu sync.Mutex
mu.Lock()
mu.Unlock()
```

---

## More Simple Examples (Copy & Run!)

### Example 14: Simple Counter with Channel

```go
// File: learn/07_concurrency_examples/14_counter_channel.go
package main

import "fmt"

func main() {
    // Channel to receive count result
    result := make(chan int)

    // Goroutine counts 1 to 100
    go func() {
        sum := 0
        for i := 1; i <= 100; i++ {
            sum += i
        }
        result <- sum  // Send result back
    }()

    // Wait and receive
    total := <-result
    fmt.Println("Sum of 1-100:", total)
}
```

**Output:**
```
Sum of 1-100: 5050
```

**How it works:**
```
Main:       [Create ch]──[Wait for result]─────────[Got 5050!]──[Print]
                │               │ (blocked)            ▲
Goroutine:      └──►[Count 1+2+3...+100]──[Send 5050]──┘
```

---

### Example 15: Ping Pong with Channels

```go
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
        time.Sleep(300 * time.Millisecond)
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
```

**Output:**
```
Ping received: 🏓
Pong received: 🏓
Ping received: 🏓
Pong received: 🏓
Ping received: 🏓
Pong received: 🏓
Game over!
```

**How it works:**
```
        ping channel          pong channel
             │                     │
   ┌─────────▼─────────┐          │
   │  Ping Player      │          │
   │  receive from ping│──────────┼───► send to pong
   └───────────────────┘          │
                                  │
                       ┌──────────▼──────────┐
                       │  Pong Player        │
   send to ping ◄──────│  receive from pong  │
                       └─────────────────────┘
```

---

### Example 16: Multiple Return Values with Channels

```go
// File: learn/07_concurrency_examples/16_multi_return.go
package main

import (
    "fmt"
    "time"
)

type Result struct {
    Value int
    Error error
}

func calculate(x, y int, ch chan<- Result) {
    time.Sleep(100 * time.Millisecond)  // Simulate work

    if y == 0 {
        ch <- Result{Error: fmt.Errorf("division by zero")}
        return
    }
    ch <- Result{Value: x / y}
}

func main() {
    results := make(chan Result, 3)

    go calculate(10, 2, results)   // 10 / 2 = 5
    go calculate(20, 4, results)   // 20 / 4 = 5
    go calculate(30, 0, results)   // Error!

    for i := 0; i < 3; i++ {
        r := <-results
        if r.Error != nil {
            fmt.Println("Error:", r.Error)
        } else {
            fmt.Println("Result:", r.Value)
        }
    }
}
```

**Output:**
```
Result: 5
Result: 5
Error: division by zero
```

---

### Example 17: Done Channel Pattern (Graceful Stop)

```go
// File: learn/07_concurrency_examples/17_done_channel.go
package main

import (
    "fmt"
    "time"
)

func worker(done <-chan bool) {
    for {
        select {
        case <-done:
            fmt.Println("Worker: Got stop signal, cleaning up...")
            time.Sleep(200 * time.Millisecond)
            fmt.Println("Worker: Cleanup done, exiting")
            return
        default:
            fmt.Println("Worker: Working...")
            time.Sleep(300 * time.Millisecond)
        }
    }
}

func main() {
    done := make(chan bool)

    go worker(done)

    // Let worker run for 1 second
    time.Sleep(1 * time.Second)

    fmt.Println("Main: Sending stop signal...")
    done <- true

    // Wait for cleanup
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Main: Exiting")
}
```

**Output:**
```
Worker: Working...
Worker: Working...
Worker: Working...
Main: Sending stop signal...
Worker: Got stop signal, cleaning up...
Worker: Cleanup done, exiting
Main: Exiting
```

---

### Example 18: Generator Pattern

```go
// File: learn/07_concurrency_examples/18_generator.go
package main

import "fmt"

// Generator returns a channel that produces values
func fibonacci(n int) <-chan int {
    ch := make(chan int)

    go func() {
        a, b := 0, 1
        for i := 0; i < n; i++ {
            ch <- a
            a, b = b, a+b
        }
        close(ch)  // Signal no more values
    }()

    return ch
}

func main() {
    fmt.Println("First 10 Fibonacci numbers:")

    for num := range fibonacci(10) {
        fmt.Print(num, " ")
    }
    fmt.Println()
}
```

**Output:**
```
First 10 Fibonacci numbers:
0 1 1 2 3 5 8 13 21 34
```

**How it works:**
```
                     ┌──────────────┐
fibonacci(10) ──►    │  Generator   │  ──► ch (returns channel)
                     │  goroutine   │
                     └──────────────┘
                            │
                   Produces values: 0, 1, 1, 2, 3, 5, 8, 13, 21, 34
                            │
                            ▼
                     for num := range ch { ... }
```

---

### Example 19: Ticker (Periodic Tasks)

```go
// File: learn/07_concurrency_examples/19_ticker.go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Tick every 500ms
    ticker := time.NewTicker(500 * time.Millisecond)
    done := make(chan bool)

    go func() {
        for {
            select {
            case <-done:
                return
            case t := <-ticker.C:
                fmt.Println("Tick at", t.Format("15:04:05.000"))
            }
        }
    }()

    // Run for 2.5 seconds
    time.Sleep(2500 * time.Millisecond)
    ticker.Stop()
    done <- true
    fmt.Println("Ticker stopped")
}
```

**Output:**
```
Tick at 14:30:22.500
Tick at 14:30:23.000
Tick at 14:30:23.500
Tick at 14:30:24.000
Tick at 14:30:24.500
Ticker stopped
```

---

### Example 20: Timeout with Context (Most Common Pattern)

```go
// File: learn/07_concurrency_examples/20_context_timeout.go
package main

import (
    "context"
    "fmt"
    "time"
)

func slowDatabase(ctx context.Context) (string, error) {
    // Simulates slow database query
    select {
    case <-time.After(3 * time.Second):
        return "data from database", nil
    case <-ctx.Done():
        return "", ctx.Err()  // Return why we stopped
    }
}

func main() {
    // Scenario 1: Query completes in time
    fmt.Println("=== Test 1: 5 second timeout (query takes 3s) ===")
    ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel1()

    start := time.Now()
    result, err := slowDatabase(ctx1)
    if err != nil {
        fmt.Printf("Failed: %v (took %v)\n", err, time.Since(start))
    } else {
        fmt.Printf("Success: %s (took %v)\n", result, time.Since(start))
    }

    // Scenario 2: Timeout before query completes
    fmt.Println("\n=== Test 2: 1 second timeout (query takes 3s) ===")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel2()

    start = time.Now()
    result, err = slowDatabase(ctx2)
    if err != nil {
        fmt.Printf("Failed: %v (took %v)\n", err, time.Since(start))
    } else {
        fmt.Printf("Success: %s (took %v)\n", result, time.Since(start))
    }
}
```

**Output:**
```
=== Test 1: 5 second timeout (query takes 3s) ===
Success: data from database (took 3.001s)

=== Test 2: 1 second timeout (query takes 3s) ===
Failed: context deadline exceeded (took 1.001s)
```

---

## Cheat Sheet: When to Use What

| Situation | Use This |
|-----------|----------|
| Run something in background | `go func() { }()` |
| Wait for goroutine to finish | `sync.WaitGroup` |
| Send data between goroutines | `chan` |
| First result wins | `select` with multiple channels |
| Timeout for operation | `context.WithTimeout()` |
| Cancel ongoing work | `context.WithCancel()` |
| Protect shared variable | `sync.Mutex` |
| Run code only once | `sync.Once` |
| Limit concurrent operations | Buffered channel as semaphore |
| Periodic task | `time.Ticker` |
| One-time delay | `time.After()` |

---

## Node.js to Go Translation Table

| Node.js | Go |
|---------|-----|
| `async function` | Regular function |
| `await promise` | Blocking call (or `<-channel`) |
| `Promise.all([...])` | Goroutines + WaitGroup |
| `Promise.race([...])` | `select` statement |
| `setTimeout(fn, ms)` | `time.After(duration)` |
| `setInterval(fn, ms)` | `time.NewTicker(duration)` |
| `AbortController` | `context.WithCancel()` |
| `signal: abortController.signal` | Pass `ctx` to functions |
| EventEmitter | Channels |
| Worker threads | Goroutines (much lighter!) |

---

## Run the Examples!

Create a folder and try these examples:

```bash
mkdir -p learn/07_concurrency_examples
cd learn/07_concurrency_examples
go mod init concurrency_examples

# Create files and run:
go run 10_parallel_api.go
go run 11_fetch_with_timeout.go
go run 12_rate_limited.go
go run 13_worker_pool.go
go run 14_counter_channel.go
go run 15_ping_pong.go
go run 16_multi_return.go
go run 17_done_channel.go
go run 18_generator.go
go run 19_ticker.go
go run 20_context_timeout.go
```
