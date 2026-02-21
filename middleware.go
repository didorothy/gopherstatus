package gopherstatus

import (
	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
)

// Middleware to add the config and the template set to the gin Context.
func GSMiddleWare(config *GSConfig) gin.HandlerFunc {
	loader := pongo2.MustNewLocalFileSystemLoader(config.TemplatePath)
	template_set := pongo2.NewSet("default", loader)
	app := GSApp{
		Config:   config,
		Template: template_set,
	}
	return func(c *gin.Context) {
		// Set up values we want in context.
		c.Set("gs_app", &app)

		// handle the request
		c.Next()
	}
}
