package tx_client

import (
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	fix44nos "github.com/quickfixgo/fix44/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
)

const ROUNDING_DECIMALS int32 = 8

type NewOrderMsgGenerator struct{}

func (NewOrderMsgGenerator) generate(tradeRequest TradeRequest) fix44nos.NewOrderSingle {
	request := fix44nos.New(
		field.NewClOrdID(tradeRequest.ClientOrderId),
		field.NewSide(enum.Side(tradeRequest.Side)),
		field.NewTransactTime(time.Now()),
		field.OrdTypeField{FIXString: quickfix.FIXString(tradeRequest.OrderType)},
	)
	request.Set(field.NewSymbol(tradeRequest.Symbol))

	quantityDecimal := decimal.NewFromFloat(tradeRequest.Quantity)
	request.Set(field.NewOrderQty(quantityDecimal, ROUNDING_DECIMALS))

	priceDecimal := decimal.NewFromFloat(tradeRequest.Price)

	if tradeRequest.OrderType == LIMIT_ORDER {
		request.Set(field.NewPrice(priceDecimal, ROUNDING_DECIMALS))
	}
	request.Set(field.NewTimeInForce(enum.TimeInForce(tradeRequest.TimeInForce)))

	return request
}
