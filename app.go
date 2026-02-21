package gopherstatus

import (
	"fmt"
	"net/http"

	"gopherstatus/status"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
)

type GSApp struct {
	Config   *GSConfig
	Template *pongo2.TemplateSet
}

// Attempt to get the GSApp pointer from the context.
func GSAppFromGinContext(c *gin.Context) (*GSApp, error) {
	app, ok := c.Get("gs_app")
	if !ok {
		return nil, fmt.Errorf("Could not find GSApp instance in context. Is the GSMiddleware being used?")
	}
	gs_app, ok := app.(*GSApp)
	if !ok {
		return nil, fmt.Errorf("The `gs_app` key in the context is not an instance of *GSApp.")
	}
	return gs_app, nil
}

// Render a template with the specified context as the response using the gin.Context.
func (gs_app *GSApp) RenderTemplate(c *gin.Context, template_name string, context pongo2.Context) {
	template, err := gs_app.Template.FromFile(template_name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Header("Content-Type", "text/html")
	err = template.ExecuteWriter(context, c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

// Uses the GSConfig and the StatusManagerFactory to initialize all the
// StatusManger instance to check status.
func InitializeStatusManagers(config *GSConfig, factory *status.StatusManagerFactory) {
	for k, v := range config.Status {
		status_config, ok := v.(map[string]any)
		if !ok {
			// TODO: log something here
			continue
		}
		status_manager_type, ok := status_config["type"].(string)
		if !ok {
			// TODO: log something her
			continue
		}
		manager, err := factory.Create(status_manager_type, status_config)
		if err != nil {
			// TODO: log something here.
			continue
		}
		config.Managers[k] = manager
	}
}
