package main

import (
	"context"
	"log"
	"os"

	"github.com/bonedaddy/go-blocknative/client"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// create the base client struct
	cl, err := client.New(context.Background(), client.Opts{
		Scheme: "wss",
		Host:   "api.blocknative.com",
		Path:   "/v0",
		// derive the api key from an environment variable
		// this sets the Client::apiKey field allowing you to retrieve the api key using
		// Client::APIKey
		APIKey: os.Getenv("BLOCKNATIVE_DAPP_ID"),
	})
	if err != nil {
		panic(err)
	}
	// this defers closure of connection and uses proper websockets connection closing semantics
	defer cl.Close()
	// send the initialization message to blocknatives api
	if err := cl.Initialize(client.NewBaseMessageMainnet(cl.APIKey())); err != nil {
		panic(err)
	}
	// subscribe to events by address
	if err := cl.WriteJSON(client.NewAddressSubscribe(
		client.NewBaseMessageMainnet(
			cl.APIKey(),
		),
		"0x88dF592F8eb5D7Bd38bFeF7dEb0fBc02cf3778a0",
	)); err != nil {
		panic(err)
	}
	// read messages in a loop
	for {
		var msg client.EthTxPayload
		if err := cl.ReadJSON(&msg); err != nil {
			if err := cl.ReadJSON(msg); err != nil {
				if e, ok := err.(*websocket.CloseError); ok {
					if e.Code != 1000 {
						log.Fatal("mempMon read", err)
					}
				}
				return
			} else {
				log.Println("receive expected close message: ", err)
				continue
			}
		}
		log.Printf("receive message:\n%+v\n", msg)
	}
}
