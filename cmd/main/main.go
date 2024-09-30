package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/robstave/rto/handlers"
	api "github.com/robstave/rto/internal"

	slogecho "github.com/samber/slog-echo"
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

	e := api.GetEcho()

	handlers.SetLogger(handlers.InitializeLogger()) // Optional: If you prefer setting a package-level logger

	// not working
	lg := handlers.GetLogger()
	mw := slogecho.New(lg)
	e.Use(mw)

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

	// Initialize events
	log.Println("init events")
	holidaysPath := filepath.Join("data", "holidays.json")
	eventsPath := filepath.Join("data", "events.json")
	handlers.InitializeEvents(holidaysPath, eventsPath)
	// Initialize preferences
	preferencesPath := filepath.Join("data", "preferences.json")
	handlers.InitializePreferences(preferencesPath)

	log.Println("starting")
	// Start the server on port 8761
	if err := e.Start(":8761"); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}

}

/*



func main() {
	e := echo.New()
	log.Println("start")

	// Middleware (optional)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	log.Println("set logger")

	handlers.SetLogger(handlers.InitializeLogger()) // Optional: If you prefer setting a package-level logger

	// not working
	lg := handlers.GetLogger()
	mw := slogecho.New(lg)
	e.Use(mw)

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

	// Static files
	e.Static("/static", "static")

	e.GET("/add-event", handlers.ShowAddEventForm) // New route to show add event form
	e.POST("/add-event", handlers.AddEvent)        // Existing POST route to handle form submission

	e.GET("/events", handlers.EventsList)
	e.GET("/prefs", handlers.ShowPrefs)
	e.POST("/prefs/update", handlers.UpdatePreferences) // New route for updating preferences

	// Routes
	e.GET("/", handlers.Home)

	// Initialize events
	log.Println("init events")
	holidaysPath := filepath.Join("data", "holidays.json")
	eventsPath := filepath.Join("data", "events.json")
	handlers.InitializeEvents(holidaysPath, eventsPath)
	// Initialize preferences
	preferencesPath := filepath.Join("data", "preferences.json")
	handlers.InitializePreferences(preferencesPath)

	// Register the new route for toggling attendance
	e.POST("/toggle-attendance", handlers.ToggleAttendance)

	// Register the new route for adding default days
	e.POST("/prefs/add-default-days", handlers.AddDefaultDays)

	log.Println("starting")
	// Start the server on port 8761
	if err := e.Start(":8761"); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}

}


*/
