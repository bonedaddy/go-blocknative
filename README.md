# go-blocknative

`go-blocknative` provides an api client for blocknative's websocket api. It allows subscribing to events by address or by transaction id and handles correct initialization as required by the documentation. It also includes the ability to store commands sent to the api in a history buffer, such that in the event of a connection failure we can restablish the current session as blocknative does not handle this on their end.

# Usage

This library provides functionality to subscribe to events from blocknatives API using either transaction hashes or addresses. In addition it provides a number of different message types as indicated below

## Message Types

## BaseMessage

The `BaseMessage` struct contains all required fields that need to be sent in messages to blocknative's API. To easily construct new base messages for the mainnet you can use `NewBaseMessageMainnet("yourApiKey")`.

## TxSubscribe

The `TxSubscribe` struct is used when subscribing/unsubscribing to events by transaction hash. If you want to send a message to subscribe to events use `NewTxSubscribe` supplying a base message along with the transaction hash to subscribe to. If you want to send a message to unsubscribe from events use `NewTxUnsubscribe` supplying a base message along with the transaction hash to unsubscribe from
## AddressSubscribe

The `AddressSubscribe` struct is like `TxSubscribe` but allows subscribing/unsubscribing to events by ethereum account addresses. If you want to send a message to subscribe to events use `NewAddressSubscribe` supplying a base message along with the address to subscribe to. If you want to send a message to unsubscribe from events use `NewAddressUnsubcribe`.

## EthTxPayload

When subscribe to events the `EthTxPayload` will be returned anytime an event is received for a transaction or address we are subscribed to. It is suitable for generalized processing of events, however you will likely want to use a use-case specific structure for better processing. Depending on the contract events being emitted they may have more information that what can be captured by this structure.

## Example

The following example show cases how to subscribe to events by an address, and reading a response. Note that you should be familiar with the mechanics of `github.com/gorilla/websockets` as this library essentially just provides helper functions around the websockets library


```Golang
package main

import (
    "log"
    "context"
    "github.com/bonedaddy/go-blocknative/client"
)

func main() {
    // create the base client struct
    cl, err := client.New(context.Background(), Opts{
        Scheme: "wss", 
        Host: "api.blocknative.com", 
        Path: "/v0",
        // derive the api key from an environment variable  
        // this sets the Client::apiKey field allowing you to retrieve the api key using
        // Client::APIKey
        APIKey: os.Getenv("BLOCKNATIVE_API"),   
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
		"someEthereumAddress",
	)); err != nil {
        panic(err)
    }
    // read messages in a loop
    for {
        var out client.EthTxPayload
        if err := cl.ReadJSON(&out); err != nil {
            // used to ignore the following event
            // websocket: close 1005 (no status)
            // this may not be necessary however it appears that
            // if we timeout on a read because no messages were received
            // this is the error emitted so we should ignore this
            if websocket.IsUnexpectedCloseError(err, 1005) {
                log.Println("receive unexpected close, exiting: ", err)
                panic(err)
            } else {
                log.Println("receive expected close message: ", err)
                continue
            }
        }
        log.Printf("receive message:\n%+v\n", out)
    }
}
```

# TODO

* Enable message history buffer usage
* Enable connection drop handling
* Enable better error handling
* Enable optional payload and subscription parameters
* Enable configuration usage