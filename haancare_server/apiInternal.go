package main

import (
	"fmt"
	"net/http"
	"strconv" // strconv 패키지 추가

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
	// URL 쿼리 파라미터 "orderStat" 값을 문자열로 가져옵니다.
	orderStatStr := c.QueryParam("orderStat")

	// 문자열을 정수(int)로 변환합니다.
	orderStat, err := strconv.Atoi(orderStatStr)
	if err != nil {
		// orderStat이 없거나 숫자가 아니면 기본값(예: 0 또는 특정 값)을 사용하거나 에러 처리할 수 있습니다.
		// 여기서는 0을 기본값으로 사용하여 모든 데이터를 조회하도록 합니다.
		orderStat = 0 
	}

	deliveries, err := GetDeliverys(OrderStat(orderStat))
	if err != nil {
		log.Error(fmt.Sprintf("failed to get deliveries: %+v (orderState: %d)", err, orderStat))
		return ResInternalServerError(c, fmt.Sprintf("failed to get deliveries: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]any{"deliveries": deliveries})
}
