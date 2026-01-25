package main

import "fmt"

func main() {
	// ============================================
	// LESSON 1: Printing Output
	// ============================================

	// fmt.Println - prints and adds newline automatically
	fmt.Println("Hello, Go!")
	fmt.Println("This is on a new line")

	// fmt.Print - prints without newline
	fmt.Print("Hello ")
	fmt.Print("World")
	fmt.Println() // just adds a newline

	// fmt.Printf - formatted printing (like C's printf)
	// %s = string, %d = integer, %f = float, %v = any value
	name := "Tejas"
	age := 25
	fmt.Printf("Name: %s, Age: %d\n", name, age)

	// ============================================
	// LESSON 2: Variables
	// ============================================

	// Method 1: var keyword (explicit type)
	var message string = "Hello"
	var count int = 10
	var price float64 = 99.99
	var isActive bool = true

	fmt.Println(message, count, price, isActive)

	// Method 2: var with type inference (Go figures out the type)
	var city = "Mumbai" // Go knows it's a string

	// Method 3: Short declaration := (most common, only inside functions)
	country := "India" // Go infers type automatically

	fmt.Println(city, country)

	// ============================================
	// LESSON 3: Basic Types
	// ============================================

	var myInt int = 42           // whole numbers
	var myFloat float64 = 3.14   // decimal numbers
	var myString string = "text" // text
	var myBool bool = true       // true or false

	fmt.Printf("int: %d, float: %f, string: %s, bool: %t\n",
		myInt, myFloat, myString, myBool)

	// ============================================
	// LESSON 4: Constants
	// ============================================

	const PI = 3.14159
	const AppName = "MyApp"
	// PI = 3.14  // ERROR! Constants cannot be changed

	fmt.Println("PI:", PI, "App:", AppName)
}
