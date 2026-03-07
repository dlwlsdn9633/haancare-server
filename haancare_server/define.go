package main

const (
	TblAlpsTokens = "haancare_alps_tokens"
	TblDelivery   = "delivery"
)

type OrderStat int

const (
	OrderStatAll     OrderStat = 0
	OrderStatPrint   OrderStat = 1
	OrderStatPending OrderStat = 5
)
