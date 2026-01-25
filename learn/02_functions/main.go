package main

import "fmt"

func main() {
	// ============================================
	// LESSON 5: Functions
	// ============================================

	// Calling a simple function
	sayHello()

	// Function with parameters
	greet("Tejas")

	// Function with return value
	result := add(5, 3)
	fmt.Println("5 + 3 =", result)

	// Function with multiple return values (very common in Go!)
	sum, product := calculate(4, 5)
	fmt.Printf("Sum: %d, Product: %d\n", sum, product)

	// Function that returns error (Go pattern)
	value, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10 / 2 =", value)
	}

	// Try dividing by zero
	value, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// Simple function - no parameters, no return
func sayHello() {
	fmt.Println("Hello!")
}

// Function with parameter
func greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// Function with return value
func add(a int, b int) int {
	return a + b
}

// Shorter: when parameters have same type
func multiply(a, b int) int {
	return a * b
}

// Function with multiple return values
func calculate(a, b int) (int, int) {
	sum := a + b
	product := a * b
	return sum, product
}

// Function with named return values
func divide(a, b int) (result int, err error) {
	if b == 0 {
		err = fmt.Errorf("cannot divide by zero")
		return // returns zero value for result and the error
	}
	result = a / b
	return // returns result and nil error
}
