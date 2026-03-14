package main

import "github.com/labstack/echo/v4"

func InitRouter(e *echo.Echo) {
	e.GET("/version", iApiVersion)
	e.GET("/deliveries", iApiDeliveries)
}