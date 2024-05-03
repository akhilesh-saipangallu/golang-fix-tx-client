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

func (NewOrderMsgGenerator) generate(orderRequest OrderRequest) fix44nos.NewOrderSingle {
	request := fix44nos.New(
		field.NewClOrdID(orderRequest.ClientOrderId),
		field.NewSide(enum.Side(orderRequest.Side)),
		field.NewTransactTime(time.Now()),
		field.OrdTypeField{FIXString: quickfix.FIXString(orderRequest.OrderType)},
	)
	request.Set(field.NewSymbol(orderRequest.Symbol))

	quantityDecimal := decimal.NewFromFloat(orderRequest.Quantity)
	request.Set(field.NewOrderQty(quantityDecimal, ROUNDING_DECIMALS))

	priceDecimal := decimal.NewFromFloat(orderRequest.Price)

	if orderRequest.OrderType == LIMIT_ORDER {
		request.Set(field.NewPrice(priceDecimal, ROUNDING_DECIMALS))
	}
	request.Set(field.NewTimeInForce(enum.TimeInForce(orderRequest.TimeInForce)))

	return request
}
