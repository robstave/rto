package api

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robstave/rto/handlers"
	"github.com/robstave/rto/internal/adapters/controller"
)

func GetEcho(rtoCtl *controller.RTOController) *echo.Echo {

	e := echo.New()

	// Middleware (optional)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	log.Println("set logger")

	//handlers.SetLogger(handlers.InitializeLogger()) // Optional: If you prefer setting a package-level logger

	// not working
	//lg := handlers.GetLogger()
	//mw := slogecho.New(lg)
	//e.Use(mw)

	// Static files
	e.Static("/static", "static")

	e.GET("/add-event", handlers.ShowAddEventForm) // New route to show add event form
	e.POST("/add-event", handlers.AddEvent)        // Existing POST route to handle form submission

	e.GET("/events", handlers.EventsList)
	e.GET("/prefs", handlers.ShowPrefs)
	e.POST("/prefs/update", handlers.UpdatePreferences) // New route for updating preferences

	// Routes
	e.GET("/", handlers.Home)

	// Register the new route for toggling attendance
	e.POST("/toggle-attendance", handlers.ToggleAttendance)

	// Register the new route for adding default days
	e.POST("/prefs/add-default-days", handlers.AddDefaultDays)

	return e
}
