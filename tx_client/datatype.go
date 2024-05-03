package tx_client

import "time"

type SideEnum string

const (
	BUY_SIDE  SideEnum = "1"
	SELL_SIDE SideEnum = "2"
)

type OrderType string

const (
	MARKET_ORDER OrderType = "1"
	LIMIT_ORDER  OrderType = "2"
)

type TimeInForce string

const (
	GTC_TIME_IN_FORCE TimeInForce = "1"
	FOK_TIME_IN_FORCE TimeInForce = "4"
	GTX_TIME_IN_FORCE TimeInForce = "5"
)

type TradeRequest struct {
	ClientOrderId string
	OrderType     OrderType
	Price         float64
	Quantity      float64
	Side          SideEnum
	Symbol        string
	TimeInForce   TimeInForce
}

type OrderReplaceRequest struct {
	OrigClientOrderId string
	ClientOrderId     string
	OrderId           string
	OrderType         OrderType
	Price             float64
	Quantity          float64
	Side              SideEnum
	Symbol            string
	TimeInForce       TimeInForce
	ExpireTime        time.Time
}

type CancelOrderRequest struct {
	OrigClientOrderId string
	ClientOrderId     string
	Side              SideEnum
	OrderId           string
}
