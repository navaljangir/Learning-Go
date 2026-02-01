package server

import (
	"gin_server/constants"
	"gin_server/middlewares"
	"gin_server/routes"

	"github.com/gin-gonic/gin"
)




// Server represents the HTTP server
type Server struct {
	router *gin.Engine
}

// NewServer creates and configures a new server instance
func NewServer() *Server {
	// Set Gin mode (debug, release, test)
	// gin.SetMode(gin.ReleaseMode) // Uncomment for production

	router := gin.New() // Create router without default middleware

	// Apply global middlewares
	router.Use(gin.Recovery())              // Recover from panics
	router.Use(middlewares.LoggerMiddleware()) // Custom logger
	router.Use(middlewares.CORSMiddleware())   // CORS support

	// Setup routes
	routes.SetupRoutes(router)

	return &Server{
		router: router,
	}
}

// Run starts the HTTP server
func (s *Server) Run() error {
	return s.router.Run(constants.ServerPort)
}

// Router returns the gin engine (useful for testing)
func (s *Server) Router() *gin.Engine {
	return s.router
}
