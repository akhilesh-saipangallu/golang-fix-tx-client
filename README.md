# Overview
This is the official Golang client for the FalconX transaction FIX gateway.

# Quick start

Under the `main.go` update the following constants
1. API_KEY
2. PASSPHRASE
3. SECRET_KEY
4. FIX_HOST
5. SENDER_COMP_ID
6. TARGET_COMP_ID
7. SOCKET_CONNECT_PORT

By default the client will try to place a `limit` order for `BTC/USD` token pair at `1 USD` price. You can modify this as needed.

# About FalconX
FalconX is an institutional digital asset brokerage. 