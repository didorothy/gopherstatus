package views

import (
	"gopherstatus"
	"gopherstatus/status"
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
)

func ListAllStatuses(c *gin.Context) {
	gs_app, err := gopherstatus.GSAppFromGinContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	statuses := make(map[string][]string)
	for k, v := range gs_app.Config.Managers {
		status := v.Status()
		statuses[k] = []string{status.Status.String(), status.Description}
	}

	gs_app.RenderTemplate(c, "statuses.html", pongo2.Context{"statuses": statuses})
}

func ServiceStatus(c *gin.Context) {
	gs_app, err := gopherstatus.GSAppFromGinContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	service_name := c.Param("service_name")

	manager, ok := gs_app.Config.Managers[service_name]
	if !ok {
		c.JSON(404, "No such service")
	}

	c.JSON(200, manager.Status())
}

func StartService(c *gin.Context) {
	gs_app, err := gopherstatus.GSAppFromGinContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	service_name := c.Param("service_name")

	manager, ok := gs_app.Config.Managers[service_name]
	if !ok {
		c.JSON(404, "No such service")
	}

	current_status := manager.Status()
	if current_status.Status == status.StatusRunning {
		c.JSON(500, "Service already running.")
	}

	result := manager.Start()
	c.JSON(200, result)
}

func StopService(c *gin.Context) {
	gs_app, err := gopherstatus.GSAppFromGinContext(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	service_name := c.Param("service_name")

	manager, ok := gs_app.Config.Managers[service_name]
	if !ok {
		c.JSON(404, "No such service")
	}

	current_status := manager.Status()
	if current_status.Status != status.StatusRunning {
		c.JSON(500, "Service already stopped.")
	}

	result := manager.Stop()
	c.JSON(200, result)
}
