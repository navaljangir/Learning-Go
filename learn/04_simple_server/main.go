package main

import (
	"fmt"
	"net/http"
)

func main() {
	// ============================================
	// LESSON 7: Your First HTTP Server
	// ============================================

	// Step 1: Create a handler function
	// This function runs when someone visits "/"
	// w = where we write the response (to send back to browser)
	// r = the incoming request (contains URL, headers, etc.)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	// Step 2: Create another route
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is the about page")
	})

	// Step 3: Start the server
	fmt.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// ============================================
// How it works:
// ============================================
//
// 1. http.HandleFunc("/", handler)
//    - Tells Go: "When someone visits /, run this function"
//
// 2. func(w http.ResponseWriter, r *http.Request)
//    - w (ResponseWriter) = use this to send data back to the user
//    - r (Request) = contains info about the incoming request
//
// 3. fmt.Fprintf(w, "Hello")
//    - Writes "Hello" to the response (what user sees in browser)
//
// 4. http.ListenAndServe(":8080", nil)
//    - Starts server on port 8080
//    - nil means "use the default router we set up with HandleFunc"
//
// To test:
//   go run main.go
//   Then open: http://localhost:8080 and http://localhost:8080/about
