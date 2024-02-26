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

	initiator, err := quickfix.NewInitiator(txApp, quickfix.NewMemoryStoreFactory(), settings, fileLogFactory)
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

	var orderPlacementErr error
	order1 := tx_client.TradeRequest{
		ClientOrderId: tx_client.GenerateRandomString(11),
		OrderType:     tx_client.LIMIT_ORDER,
		Price:         3300.412345678923456,
		Quantity:      0.00012,
		Side:          tx_client.SELL_SIDE,
		Symbol:        "BTC/USD",
		TimeInForce:   tx_client.FOK_TIME_IN_FORCE,
	}

	if orderPlacementErr = txApp.PlaceNewOrder(initiator, order1); orderPlacementErr != nil {
		log.Println("orderPlacementErr: ", orderPlacementErr)
		initiator.Stop()
		os.Exit(1)
	}
	time.Sleep(5 * time.Second)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit
}
