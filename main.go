package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robstave/rto/handlers"
)

// TemplateRenderer is a custom renderer for Echo
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	log.Println("start")

	// Middleware (optional)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	funcMap := template.FuncMap{
		"formatDate": func(t time.Time, layout string) string {
			return t.Format(layout)
		},
	}

	// Parse the templates with custom functions
	renderer := &TemplateRenderer{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer

	log.Println("templates load")

	// Static files
	e.Static("/static", "static")

	e.GET("/auth", func(c echo.Context) error {
		log.Println("test")

		return c.String(http.StatusOK, "auth")
	})

	// Routes
	e.GET("/", handlers.Home)

	log.Println("Route '/' registered with handlers.Home")

	/*
		e.GET("/", func(c echo.Context) error {
			data := struct {
				CurrentDate time.Time
			}{
				CurrentDate: time.Now(),
			}

			err := c.Render(http.StatusOK, "home.html", data)
			if err != nil {
				log.Printf("Template rendering error: %v", err)
				return c.String(http.StatusInternalServerError, "Internal Server Error")
			}

			return nil
		})
	*/

	log.Println("starting")
	// Start the server on port 8761
	if err := e.Start(":8761"); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}
}
