package main

import "fmt"

func main() {
	// ============================================
	// LESSON 6: Printf Format Specifiers
	// ============================================

	// These are placeholders that get replaced by values

	name := "Tejas"
	age := 25
	height := 5.9
	isStudent := true

	// %s - String
	fmt.Printf("Name: %s\n", name)

	// %d - Integer (decimal)
	fmt.Printf("Age: %d\n", age)

	// %f - Float (decimal number)
	fmt.Printf("Height: %f\n", height)    // 5.900000
	fmt.Printf("Height: %.1f\n", height)  // 5.9 (1 decimal place)
	fmt.Printf("Height: %.2f\n", height)  // 5.90 (2 decimal places)

	// %t - Boolean (true/false)
	fmt.Printf("Is Student: %t\n", isStudent)

	// %v - Any value (Go figures out the format)
	fmt.Printf("Name: %v, Age: %v, Height: %v\n", name, age, height)

	// %T - Type of the variable
	fmt.Printf("Type of name: %T\n", name)   // string
	fmt.Printf("Type of age: %T\n", age)     // int
	fmt.Printf("Type of height: %T\n", height) // float64

	// ============================================
	// Escape Characters
	// ============================================

	// \n - New line
	fmt.Printf("Line 1\nLine 2\nLine 3\n")

	// \t - Tab
	fmt.Printf("Column1\tColumn2\tColumn3\n")

	// \\ - Backslash
	fmt.Printf("Path: C:\\Users\\Tejas\n")

	// \" - Quote inside string
	fmt.Printf("He said \"Hello\"\n")

	// ============================================
	// Putting it together
	// ============================================

	port := ":8080"
	// This: "http://localhost%s\n" with port=":8080"
	// Becomes: "http://localhost:8080" + newline
	fmt.Printf("Server running at http://localhost%s\n", port)
}
