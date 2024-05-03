package tx_client

import (
	"log"

	"github.com/quickfixgo/quickfix"
)

type TransactionClient struct {
	ApiKey     string
	Passphrase string
	SecretKey  string
	SessionId  quickfix.SessionID
}

func (tc *TransactionClient) OnCreate(sessionID quickfix.SessionID) {
	tc.SessionId = sessionID
}

// Upon login subscribe to all symbols
func (tc *TransactionClient) OnLogon(sessionID quickfix.SessionID) {
	log.Println("logon successful, session id: ", sessionID)
}

func (tc TransactionClient) OnLogout(sessionID quickfix.SessionID) {}

func (tc TransactionClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	if msg.IsMsgTypeOf("A") {
		log.Println("Sending logon")
		LogonProcessor{}.process(msg, tc.ApiKey, tc.Passphrase, tc.SecretKey)
	}
	log.Println("ToAdmin msg: ", msg)
}

func (tc TransactionClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	log.Println("ToApp msg: ", msg)
	return nil
}

func (tc TransactionClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	log.Println("FromAdmin msg: ", msg.ToMessage())
	return nil
}

func (tc TransactionClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	log.Println("FromApp msg: ", msg)
	return
}

func (tc *TransactionClient) PlaceNewOrder(initiator *quickfix.Initiator, tradeRequest TradeRequest) error {
	nosMsg := NewOrderMsgGenerator{}.generate(tradeRequest)
	return quickfix.SendToTarget(nosMsg, tc.SessionId)
}

func (tc *TransactionClient) ReplaceOrder(initiator *quickfix.Initiator, orderReplaceRequest OrderReplaceRequest) error {
	nosMsg := OrderCancelReplaceMsgGenerator{}.generate(orderReplaceRequest)
	return quickfix.SendToTarget(nosMsg, tc.SessionId)
}

func (tc *TransactionClient) CancelOrder(initiator *quickfix.Initiator, cancelRequest CancelOrderRequest) error {
	msg := OrderCancelMsgGenerator{}.generate(cancelRequest)
	return quickfix.SendToTarget(msg, tc.SessionId)
}
