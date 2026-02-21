package gopherstatus

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
)

func TestGSMiddlewareSetsAppInContext(t *testing.T) {
	// --------------------------------------------------------------------
	// Setup: create a temporary directory that contains a simple template
	// --------------------------------------------------------------------
	tmpDir := t.TempDir()

	// Create a dummy template file that will be rendered by pongo2
	const tmplName = "hello.tmpl"
	const tmplContent = "Hello {{ name }}!"
	if err := os.WriteFile(filepath.Join(tmpDir, tmplName), []byte(tmplContent), 0644); err != nil {
		t.Fatalf("unable to write temporary template: %v", err)
	}

	// --------------------------------------------------------------------
	// Create a configuration that points to the temp template dir
	// --------------------------------------------------------------------
	cfg := &GSConfig{
		TemplatePath: tmpDir,
		// The other fields are zero‑values – they are not needed for this test.
	}

	// --------------------------------------------------------------------
	// Prepare a Gin engine that uses the middleware
	// --------------------------------------------------------------------
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(GSMiddleWare(cfg))

	// The route will read the GSApp from the context, load the template
	// and return some data back so we can assert on it.
	engine.GET("/test", func(c *gin.Context) {
		app, err := GSAppFromGinContext(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Load the template we created above
		tmpl, err := app.Template.FromFile(tmplName)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("cannot load template: %v", err)})
			return
		}

		// Render the template
		out, err := tmpl.Execute(pongo2.Context{"name": "world"})
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("cannot render template: %v", err)})
			return
		}

		// Return everything we need to assert in the test
		c.JSON(200, gin.H{
			"rendered": out,
		})
	})

	// --------------------------------------------------------------------
	// Exercise: make a request and assert the response
	// --------------------------------------------------------------------
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// The response body is JSON – unmarshal it
	var resp struct {
		TemplateName string `json:"template_name"`
		Rendered     string `json:"rendered"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("cannot unmarshal response JSON: %v", err)
	}

	expectedRendered := "Hello world!"
	if resp.Rendered != expectedRendered {
		t.Fatalf("unexpected rendered content – want %q, got %q", expectedRendered, resp.Rendered)
	}
}

func TestGSMiddleware_PanicsWithBadTemplatePath(t *testing.T) {
	// A path that definitely does not exist
	badPath := filepath.Join(os.TempDir(), "nonexistent", "template_dir")

	cfg := &GSConfig{
		TemplatePath: badPath,
	}

	// MustNewLocalFileSystemLoader will panic if the path is invalid.
	// We want to make sure the middleware does indeed panic in that case.
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected GSMiddleWare to panic with bad template path, but it did not")
		}
	}()

	// The middleware is created here – it will panic during its
	// initialization phase, before the request handling.
	_ = GSMiddleWare(cfg)
}
