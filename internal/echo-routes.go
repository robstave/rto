package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gorilla/sessions"

	"github.com/labstack/echo-contrib/session"
	"github.com/robstave/rto/internal/adapters/controller"
)

func GetEcho(rtoCtl *controller.RTOController) *echo.Echo {

	e := echo.New()

	// Middleware (optional)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Pre(middleware.RemoveTrailingSlash())
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	store := sessions.NewCookieStore([]byte("    ffffff"))
	store.Options = &sessions.Options{
		Path: "/",
		//Domain:   "192.168.86.176", // Update as per your domain
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false,                // Set to true in production (requires HTTPS)
		SameSite: http.SameSiteLaxMode, // Adjust as needed
	}
	e.Use(session.Middleware(store))

	// Static files
	e.Static("/static", "static")

	// Public Routes
	e.GET("/login", rtoCtl.ShowLoginForm)
	e.POST("/login", rtoCtl.ProcessLogin)
	e.GET("/logout", rtoCtl.Logout)

	// Group protected routes
	r := e.Group("")
	r.Use(rtoCtl.AuthMiddleware)

	r.GET("/add-event", rtoCtl.ShowAddEventForm) //  show add event form
	r.POST("/add-event", rtoCtl.AddEvent)        //  handle form submission

	r.GET("/events", rtoCtl.EventsList)
	r.GET("/prefs", rtoCtl.ShowPrefs)
	r.POST("/prefs/update", rtoCtl.UpdatePreferences) // New route for updating preferences

	// Routes
	r.GET("/", rtoCtl.Home)
	r.GET("", rtoCtl.Home)

	r.POST("/toggle-attendance", rtoCtl.ToggleAttendance)

	r.POST("/prefs/add-default-days", rtoCtl.AddDefaultDays)
	r.DELETE("/events/delete/:id", rtoCtl.DeleteEvent)
	r.POST("/add-events-json", rtoCtl.BulkAddEventsJSON)

	r.DELETE("/events/clear/:date", rtoCtl.ClearEventsForDate)

	r.GET("/export/markdown", rtoCtl.ExportEventsMarkdown)

	r.GET("/chart-data", rtoCtl.GetChartData)

	return e
}
