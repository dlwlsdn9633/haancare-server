package main

import (
	"fmt"
	"net/http"

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
	var body DeliveryGetBody
	var err error
	if err = c.Bind(&body); err != nil {
		log.Error(fmt.Sprintf("failed to bind body: %+v", err))
		return ResInternalServerError(c, fmt.Sprintf("failed to bind body: %v", err))
	}

	var deliveries []Delivery
	deliveries, err = GetDeliverys(OrderStat(body.OrderStat))
	if err != nil {
		log.Error(fmt.Sprintf("failed to get deliveries: %+v (orderState: %d)", err, body.OrderStat))
		return ResInternalServerError(c, fmt.Sprintf("failed to get deliveries: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]any{"deliveries": deliveries})
}
