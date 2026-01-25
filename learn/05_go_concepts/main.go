package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	// ============================================
	// 1. STRUCTS (like JS objects, but typed)
	// ============================================
	fmt.Println("=== STRUCTS ===")

	// JS: const user = { name: "Tejas", age: 25 }
	// Go: must define structure first

	user1 := User{Name: "Tejas", Age: 25}
	fmt.Println(user1)
	fmt.Println("Name:", user1.Name)

	// Struct with method
	user1.Birthday() // Age becomes 26
	fmt.Println("After birthday:", user1.Age)

	// ============================================
	// 2. POINTERS (& and *)
	// ============================================
	fmt.Println("\n=== POINTERS ===")

	// Why pointers? Two reasons:
	// 1. Modify original value in functions
	// 2. Avoid copying large data

	x := 10
	fmt.Println("x before:", x)

	// & = "address of" - gives memory location
	// * = "value at" - gives value at that address

	p := &x          // p points to x's memory address
	fmt.Println("address of x:", p)
	fmt.Println("value at p:", *p)

	*p = 20          // change value through pointer
	fmt.Println("x after:", x) // x is now 20!

	// Without pointer - function gets a COPY
	num := 100
	doubleValue(num)
	fmt.Println("without pointer:", num) // still 100

	// With pointer - function modifies ORIGINAL
	doublePointer(&num)
	fmt.Println("with pointer:", num) // now 200

	// ============================================
	// 3. ERROR HANDLING (no try/catch!)
	// ============================================
	fmt.Println("\n=== ERROR HANDLING ===")

	// Go pattern: functions return (result, error)
	// You MUST check error before using result

	// Success case
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10 / 2 =", result)
	}

	// Error case
	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// Creating custom errors
	err = validateAge(-5)
	if err != nil {
		fmt.Println("Validation failed:", err)
	}

	// ============================================
	// 4. INTERFACES (polymorphism in Go)
	// ============================================
	fmt.Println("\n=== INTERFACES ===")

	// Interface = contract defining behavior
	// Any type that has the methods = implements interface
	// NO "implements" keyword needed (implicit)

	dog := Dog{Name: "Bruno"}
	cat := Cat{Name: "Whiskers"}

	// Both implement Animal interface
	makeSound(dog) // Bruno says: Woof!
	makeSound(cat) // Whiskers says: Meow!

	// Slice of different types through interface
	animals := []Animal{dog, cat}
	for _, animal := range animals {
		animal.Speak()
	}

	// ============================================
	// 5. GOROUTINES (lightweight threads)
	// ============================================
	fmt.Println("\n=== GOROUTINES ===")

	// go keyword = run function in background
	// Much lighter than OS threads (can have millions)

	// This runs in background
	go printNumbers("goroutine")

	// This runs in main thread
	printNumbers("main")

	// Problem: main might exit before goroutine finishes
	// Solution: channels (next section) or sync.WaitGroup

	time.Sleep(100 * time.Millisecond) // wait for goroutine (bad practice, use channels)

	// ============================================
	// 6. CHANNELS (goroutine communication)
	// ============================================
	fmt.Println("\n=== CHANNELS ===")

	// Channels = pipes for sending data between goroutines
	// make(chan Type) creates a channel

	// Create channel that carries strings
	messages := make(chan string)

	// Start goroutine that sends to channel
	go func() {
		messages <- "Hello from goroutine!" // send
	}()

	// Receive from channel (blocks until data arrives)
	msg := <-messages // receive
	fmt.Println(msg)

	// Practical example: parallel tasks
	ch := make(chan int)

	go calculateSquare(5, ch)
	go calculateSquare(10, ch)

	// Receive both results
	result1 := <-ch
	result2 := <-ch
	fmt.Println("Squares:", result1, result2)

	// Buffered channel (doesn't block until buffer full)
	buffered := make(chan int, 2) // buffer size 2
	buffered <- 1
	buffered <- 2
	// buffered <- 3  // this would block (buffer full)
	fmt.Println(<-buffered, <-buffered)

	// ============================================
	// 7. SELECT (handle multiple channels)
	// ============================================
	fmt.Println("\n=== SELECT ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch1 <- "from channel 1"
	}()

	go func() {
		time.Sleep(20 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	// Select waits on multiple channels
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		}
	}

	fmt.Println("\n=== DONE ===")
}

// ============================================
// STRUCT DEFINITION
// ============================================

type User struct {
	Name string
	Age  int
}

// Method on struct (like class method in JS)
// (u *User) = receiver - which struct this method belongs to
// * = pointer receiver (can modify the struct)
func (u *User) Birthday() {
	u.Age++ // modifies original because pointer receiver
}

// ============================================
// POINTER FUNCTIONS
// ============================================

func doubleValue(n int) {
	n = n * 2 // only modifies local copy
}

func doublePointer(n *int) {
	*n = *n * 2 // modifies original through pointer
}

// ============================================
// ERROR HANDLING FUNCTIONS
// ============================================

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil // nil = no error
}

func validateAge(age int) error {
	if age < 0 {
		return fmt.Errorf("age cannot be negative: got %d", age)
	}
	return nil
}

// ============================================
// INTERFACE DEFINITION
// ============================================

type Animal interface {
	Speak() // any type with Speak() method implements Animal
}

type Dog struct {
	Name string
}

func (d Dog) Speak() {
	fmt.Printf("%s says: Woof!\n", d.Name)
}

type Cat struct {
	Name string
}

func (c Cat) Speak() {
	fmt.Printf("%s says: Meow!\n", c.Name)
}

// Function accepting interface = accepts any type implementing it
func makeSound(a Animal) {
	a.Speak()
}

// ============================================
// GOROUTINE/CHANNEL FUNCTIONS
// ============================================

func printNumbers(label string) {
	for i := 1; i <= 3; i++ {
		fmt.Printf("%s: %d\n", label, i)
		time.Sleep(10 * time.Millisecond)
	}
}

func calculateSquare(n int, ch chan int) {
	ch <- n * n // send result to channel
}
