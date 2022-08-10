package main

import (
	"context"
	"log"
	"os"

	"github.com/bonedaddy/go-blocknative/client"
	"github.com/gorilla/websocket"
)

func main() {

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
		"0x68b3465833fb72a70ecdf485e0e4c7bd8665fc45", //Uniswap router address with high volume from Blocknative website.
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
					log.Fatal("websocket closed", err)
				}
				log.Println("failed to read msg", err)
				continue
			} else {
				log.Println("receive expected close message: ", err)
				continue
			}
			log.Println("failed to read msg memory", err)
			continue
		}
		log.Printf("receive message:\n%+v\n", msg)
	}
}
