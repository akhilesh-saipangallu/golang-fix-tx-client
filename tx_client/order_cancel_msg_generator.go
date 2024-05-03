package tx_client

import (
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	fix44ocr "github.com/quickfixgo/fix44/ordercancelrequest"
)

type OrderCancelMsgGenerator struct{}

func (OrderCancelMsgGenerator) generate(cancelRequest CancelOrderRequest) fix44ocr.OrderCancelRequest {
	request := fix44ocr.New(
		field.NewOrigClOrdID(cancelRequest.OrigClientOrderId),
		field.NewClOrdID(cancelRequest.ClientOrderId),
		field.NewSide(enum.Side(cancelRequest.Side)),
		field.NewTransactTime(time.Now()),
	)

	request.Set(field.NewOrderID(cancelRequest.OrderId))
	return request
}
