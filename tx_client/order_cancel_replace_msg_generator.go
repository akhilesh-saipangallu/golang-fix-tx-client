package tx_client

import (
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	fix44ocrr "github.com/quickfixgo/fix44/ordercancelreplacerequest"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
)

type OrderCancelReplaceMsgGenerator struct{}

func (OrderCancelReplaceMsgGenerator) generate(orderReplaceRequest OrderReplaceRequest) fix44ocrr.OrderCancelReplaceRequest {
	request := fix44ocrr.New(
		field.NewOrigClOrdID(orderReplaceRequest.OrigClientOrderId),
		field.NewClOrdID(orderReplaceRequest.ClientOrderId),
		field.NewSide(enum.Side(orderReplaceRequest.Side)),
		field.NewTransactTime(time.Now()),
		field.OrdTypeField{FIXString: quickfix.FIXString(orderReplaceRequest.OrderType)},
	)
	request.Set(field.NewOrderID(orderReplaceRequest.OrderId))
	request.Set(field.NewTimeInForce(enum.TimeInForce(orderReplaceRequest.TimeInForce)))
	request.Set(field.NewSymbol(orderReplaceRequest.Symbol))
	quantityDecimal := decimal.NewFromFloat(orderReplaceRequest.Quantity)
	request.Set(field.NewOrderQty(quantityDecimal, ROUNDING_DECIMALS))

	if orderReplaceRequest.OrderType == LIMIT_ORDER {
		priceDecimal := decimal.NewFromFloat(orderReplaceRequest.Price)
		request.Set(field.NewPrice(priceDecimal, ROUNDING_DECIMALS))
	}

	if orderReplaceRequest.TimeInForce == GTX_TIME_IN_FORCE {
		request.Set(field.NewExpireTime(orderReplaceRequest.ExpireTime))
	}

	return request
}
