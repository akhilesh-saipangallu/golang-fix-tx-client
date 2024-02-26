package tx_client

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
