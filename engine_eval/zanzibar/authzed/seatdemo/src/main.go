package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"seatdemo/handler"
)

func main() {
	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handler.GetInfo)
	//userAccess
	e.GET("/tenant/:tenant/user", handler.GetTenantUsers)

	//product
	e.GET("/tenant/:tenant/product/:pinstance/license", handler.GetTenantUsers)
	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
