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


## Examples

The `examples` folder has some full running examples. Note that you should be familiar with the mechanics of `github.com/gorilla/websockets` as this library essentially just provides helper functions around the websockets library


# TODO

* Enable message history buffer usage
* Enable connection drop handling
* Enable better error handling
* Enable optional payload and subscription parameters
* Enable configuration usage