package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoHandle *echo.Echo

func initializeRoutes() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover()) //TODO: eval real necessary middlewares, this is just added as per the docs

	// Routes
	e.GET("/", GetInfo)
	//userAccess
	e.GET("/tenant/:tenant/user", GetTenantUsers)

	//product
	e.GET("/tenant/:tenant/product/:pinstance/license", GetLicenseInfoForProductInstance)
	e.POST("/tenant/:tenant/product/:pinstance/license", GrantLicenseIfNotFull)

	return e
}

func GetEcho() *echo.Echo {
	if echoHandle == nil {
		echoHandle = initializeRoutes()
	}

	return echoHandle
}
