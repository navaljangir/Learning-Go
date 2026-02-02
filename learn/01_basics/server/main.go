package main

import "fmt"

type Server struct {
    name string
    port int
}

// Constructor - creates the instance
func NewServer(name string, port int) *Server {
    return &Server{
        name: name,
        port: port,
    }
}

// Method with pointer receiver
func (s *Server) Start() {
    // "s" is the Server instance
    // It's automatically passed when you call server.Start()
    fmt.Printf("Server '%s' starting on port %d\n", s.name, s.port)
}

func main() {
    // Step 1: Create instance
    server := NewServer("API", 8080)
    //  └─ This is a *Server instance
    
    // Step 2: Call method
    server.Start()
    //  │      └─ Regular method call
    //  └─ This becomes "s" inside Start()
    
    // What Go does internally:
    // Start(server)  // server is passed as "s"
}