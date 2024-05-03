package main

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/falconxio/fix_tx_client/configs"
	"github.com/falconxio/fix_tx_client/tx_client"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/store/file"
)

const (
	API_KEY    string = ""
	PASSPHRASE string = ""
	SECRET_KEY string = ""

	FIX_HOST            string = ""
	SENDER_COMP_ID      string = ""
	TARGET_COMP_ID      string = ""
	SOCKET_CONNECT_PORT string = ""
)

func GetConfig(sessionConfig configs.FixSessionConfig) (io.Reader, error) {
	configTemplateFileName := path.Join("configs", "tx_client.cfg")
	templateContent, err := ioutil.ReadFile(configTemplateFileName)
	if err != nil {
		return nil, err
	}

	templateString := string(templateContent)
	tmpl, err := template.New("config").Parse(templateString)
	if err != nil {
		return nil, err
	}

	var configString bytes.Buffer
	err = tmpl.Execute(&configString, sessionConfig)
	if err != nil {
		return nil, err
	}
	rawBytes := configString.Bytes()

	return bytes.NewReader(rawBytes), nil
}

func getTxClient() (*tx_client.TransactionClient, *quickfix.Initiator) {
	config, configErr := GetConfig(configs.FixSessionConfig{
		FixHost:           FIX_HOST,
		FileLogPath:       path.Join("fix_logs"),
		FileStorePath:     path.Join("fix_logs", "store"),
		SenderCompID:      SENDER_COMP_ID,
		TargetCompID:      TARGET_COMP_ID,
		SocketConnectPort: SOCKET_CONNECT_PORT,
	})

	if configErr != nil {
		log.Println("configErr: ", configErr)
	}

	settings, err := quickfix.ParseSettings(config)
	if err != nil {
		log.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	txApp := &tx_client.TransactionClient{
		ApiKey:     API_KEY,
		Passphrase: PASSPHRASE,
		SecretKey:  SECRET_KEY,
	}
	fileLogFactory, err := quickfix.NewFileLogFactory(settings)
	if err != nil {
		log.Printf("Error creating FileLogFactory: %s\n", err)
		os.Exit(1)
	}

	initiator, err := quickfix.NewInitiator(txApp, file.NewStoreFactory(settings), settings, fileLogFactory)
	if err != nil {
		log.Printf("Unable to create Initiator: %s\n", err)
		os.Exit(1)
	}
	return txApp, initiator
}

func main() {
	txApp, initiator := getTxClient()
	initiator.Start()
	defer initiator.Stop()
	time.Sleep(5 * time.Second)

	// place order
	var err error
	order1 := tx_client.OrderRequest{
		ClientOrderId: tx_client.GenerateRandomString(11),
		OrderType:     tx_client.LIMIT_ORDER,
		Price:         59490.48,
		Quantity:      0.00012,
		Side:          tx_client.SELL_SIDE,
		Symbol:        "BTC/USD",
		TimeInForce:   tx_client.FOK_TIME_IN_FORCE,
	}

	if err = txApp.PlaceNewOrder(initiator, order1); err != nil {
		log.Println("orderPlacementErr: ", err)
		initiator.Stop()
		os.Exit(1)
	}
	time.Sleep(5 * time.Second)

	// cancel order
	cancelRequest := tx_client.CancelOrderRequest{
		OrigClientOrderId: order1.ClientOrderId,
		ClientOrderId:     tx_client.GenerateRandomString(11),
		Side:              tx_client.SELL_SIDE,
		OrderId:           "orderId",
	}
	if err = txApp.CancelOrder(initiator, cancelRequest); err != nil {
		log.Println("orderCancelErr: ", err)
		initiator.Stop()
		os.Exit(1)
	}

	// cancel replace
	cancelReplaceRequest := tx_client.OrderReplaceRequest{
		OrigClientOrderId: order1.ClientOrderId,
		ClientOrderId:     tx_client.GenerateRandomString(11),
		OrderId:           "orderId",
		OrderType:         tx_client.LIMIT_ORDER,
		Price:             59490.48,
		Quantity:          0.00012,
		Side:              tx_client.SELL_SIDE,
		Symbol:            "BTC/USD",
		TimeInForce:       tx_client.GTX_TIME_IN_FORCE,
		ExpireTime:        time.Now().Add(10 * time.Minute),
	}
	if err = txApp.ReplaceOrder(initiator, cancelReplaceRequest); err != nil {
		log.Println("orderCancelReplaceErr: ", err)
		initiator.Stop()
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	initiator.Stop()
}
