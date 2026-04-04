package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type DeliveryGetBody struct {
	OrderStat int `json:"orderStat"`
}

type DeliveryPutBody struct {
	OrderNum  string `json:"orderNumber"`
	OrderStat int    `json:"orderStat"`
}

func iApiVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"version": "1.0.0"})
}

func iApiDeliveries(c echo.Context) error {
	orderStat := c.Param("orderStat")
	if IsNullOrEmpty(&orderStat) {
		errMsg := "orderStat is required"
		log.Error(errMsg)
		return ResBadRequest(c, errMsg)
	}

	orderStatInt, err := strconv.Atoi(orderStat)
	if err != nil {
		log.Error(fmt.Sprintf("error while casting string to int: %+v - %s", err, orderStat))
		return ResInternalServerError(c, fmt.Sprintf("error while casting string to int: %v - %s", err, orderStat))
	}

	deliveries, err := GetDeliverys(OrderStat(orderStatInt))
	if err != nil {
		log.Error(fmt.Sprintf("failed to get deliveries: %+v (orderState: %s)", err, orderStat))
		return ResInternalServerError(c, fmt.Sprintf("failed to get deliveries: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]any{"deliveries": deliveries})
}

func iApiInvoice(c echo.Context) error {
	orderNum := c.Param("orderNum")
	if IsNullOrEmpty(&orderNum) {
		errMsg := "orderNum is required"
		log.Error(errMsg)
		return ResBadRequest(c, errMsg)
	}

	token := c.Request().Header.Get("token")
	if IsNullOrEmpty(&token) {
		errMsg := "token is required"
		log.Error(errMsg)
		return ResBadRequest(c, errMsg)
	}

	baseUrls := []string{
		"https://pid.alps.llogis.com:18210/pid/ftr/pacltrc/inner/popinvrgstinfo",
		"https://pid.alps.llogis.com:18210/pid/ftr/pacltrc/inner/popinvrsrvinfo",
	}

	var allResults []OrderResult
	for _, url := range baseUrls {
		results, err := GetAlpsOrders(url, orderNum, token)
		if err != nil {
			log.Error(fmt.Sprintf("error while get alps orders: %+v", err))
			return ResInternalServerError(c, fmt.Sprintf("error while get alps orders: %v", err))
		}
		allResults = append(allResults, results...)
	}
	return c.JSON(http.StatusOK, allResults)
}

func iApiToken(c echo.Context) error {
	token, err := GetToken()
	if err != nil {
		log.Error(fmt.Sprintf("error while getting token: %+v", err))
		return ResInternalServerError(c, fmt.Sprintf("error while getting token: %v", err))
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
