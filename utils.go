package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Assert(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("ASSERTION FAILED: %s", msg))
	}
}

func ResBadRequest(c echo.Context, err string) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{"err": err})
}

func ResInternalServerError(c echo.Context, err string) error {
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{"err": err})
}
