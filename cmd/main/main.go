package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	api "github.com/robstave/rto/internal"

	"github.com/labstack/echo/v4"
	"github.com/robstave/rto/internal/adapters/controller"
	"github.com/robstave/rto/logger"
	slogecho "github.com/samber/slog-echo"
)

// Quarter configuration – centralizing start/end dates for the quarter
var (
	QuarterStart = time.Date(2024, time.December, 30, 0, 0, 0, 0, time.UTC)
	QuarterEnd   = time.Date(2025, time.March, 31, 0, 0, 0, 0, time.UTC)
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

	// Read DB_PATH from environment variable, set a default if not provided
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/db.sqlite3" // Default path
	}

	slogger := logger.InitializeLogger()
	logger.SetLogger(slogger) // Optional: If you prefer setting a package-level logger
	rtoClt := controller.NewRTOController(dbPath, slogger, QuarterStart, QuarterEnd)

	// Initialize session middleware with a cookie store

	e := api.GetEcho(rtoClt)
	mw := slogecho.New(slogger)
	e.Use(mw)

	//e.Use(session.Middleware(sessions.NewCookieStore([]byte())))

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

	log.Println("templates loaded")

	log.Println("starting")
	// Start the server on port 8761
	if err := e.Start(":8761"); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}

}
