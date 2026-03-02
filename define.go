package main

const (
	TblAlpsTokens = "haancare_alps_tokens"
	TblDelivery   = "delivery"
)

type OrderStat int

const (
	OrderStatPrint   OrderStat = iota + 1
	OrderStatPending OrderStat = iota + 4
)
