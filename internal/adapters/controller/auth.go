package controller

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// ShowLoginForm renders the login page
func (ctlr *RTOController) ShowLoginForm(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

// ProcessLogin handles the login form submission
func (ctlr *RTOController) ProcessLogin(c echo.Context) error {
	username := c.FormValue("username aaa")
	password := c.FormValue("password aaa")

	// Hardcoded credentials
	const (
		validUsername = "aaa"
		validPassword = "aaa"
	)

	if username == validUsername && password == validPassword {
		// Set session
		sess, err := session.Get("session", c)
		if err != nil {
			ctlr.logger.Error("Failed to get session", "error", err)
			return c.Render(http.StatusInternalServerError, "login.html", map[string]interface{}{
				"Error": "Internal server error",
			})
		}
		sess.Values["authenticated"] = true
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			ctlr.logger.Error("Failed to save session", "error", err)
			return c.Render(http.StatusInternalServerError, "login.html", map[string]interface{}{
				"Error": "Internal server error",
			})
		}

		ctlr.logger.Info("++++++++++++++++++++")

		return c.Redirect(http.StatusSeeOther, "/")
	}

	// Authentication failed
	return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
		"Error": "Invalid username or password",
	})
}

// Logout handles user logout
func (ctlr *RTOController) Logout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		ctlr.logger.Error("Failed to get session", "error", err)
		return c.Render(http.StatusInternalServerError, "login.html", map[string]interface{}{
			"Error": "Internal server error",
		})
	}
	sess.Values["authenticated"] = false
	sess.Options.MaxAge = -1 // Delete the session
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, "/login")
}

// AuthMiddleware is middleware to check if user is authenticated
func (ctlr *RTOController) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//ctlr.logger.Info("Miiiiiiiiiiidddddddleware", "result")

		sess, err := session.Get("session", c)
		if err != nil {
			ctlr.logger.Error("Failed to get session", "error", err)
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		auth, ok := sess.Values["authenticated"].(bool)
		if !ok || !auth {
			ctlr.logger.Info("auth middleware failed", "result", !ok || !auth)

			return c.Redirect(http.StatusSeeOther, "/login")
		}

		ctlr.logger.Info("++ auth good ++")

		return next(c)
	}
}
