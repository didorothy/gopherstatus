package main

import (
	"errors"
	"fmt"
	"gopherstatus/status"
	"gopherstatus/views"
	"os"
	"path/filepath"

	"gopherstatus"

	"github.com/gin-gonic/gin"
)

// If there is a static folder at the same level as the templates, serve the static files.
func include_static_route(config *gopherstatus.GSConfig, router *gin.Engine) {
	static_directory := filepath.Join(filepath.Dir(config.TemplatePath), "static")
	_, err := os.Stat(static_directory)
	if !errors.Is(err, os.ErrNotExist) {
		router.Static("/static", static_directory)
	}
}

// All routes should be defined in this function.
func build_routes(config *gopherstatus.GSConfig) *gin.Engine {
	router := gin.Default()
	// Tell Gin to use our middlware
	router.Use(gopherstatus.GSMiddleWare(config))

	include_static_route(config, router)

	// Define the routes.
	router.GET("/", views.ListAllStatuses)
	router.POST("/status/:service_name/", views.ServiceStatus)
	router.POST("/start/:service_name/", views.StartService)
	router.POST("/stop/:service_name/", views.StopService)
	return router
}

func main() {
	factory := status.NewStatusManagerFactory()
	factory.Register("docker", status.CreateDockerManager)
	factory.Register("systemctl", status.CreateSystemctlManager)

	config, err := gopherstatus.ParseConfig("config.toml")
	gopherstatus.InitializeStatusManagers(config, factory)
	if err != nil {
		fmt.Println("Failed to parse config.toml file. (config.toml file must be in the same directory as the executable.)", err)
	}

	router := build_routes(config)

	fmt.Println(*config)

	// Start server
	serve_on := fmt.Sprintf("%s:%d", config.IP, config.Port)
	fmt.Printf("Serving on %s\n", serve_on)
	router.Run(serve_on)
}
