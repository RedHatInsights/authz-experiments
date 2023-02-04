package main

import (
	"seatdemo/seatdemo/src/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handler.HelloHandler)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
