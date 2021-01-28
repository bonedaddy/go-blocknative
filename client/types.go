package client

import (
	"os"
	"time"
)

// BaseMessage is the base message required for all interactions with the websockets api
type BaseMessage struct {
	CategoryCode string `json:"categoryCode"`
	EventCode    string `json:"eventCode"`
	Timestamp    string `json:"timeStamp"`
	DappID       string `json:"dappId"` // api key
	Version      string `json:"version"`
	Blockchain   `json:"blockchain"`
}

// Blockchain is a type fulfilling the blockchain params
type Blockchain struct {
	System  string `json:"system"`
	Network string `json:"network"`
}

// TxSubscribe is used to subscribe/unsubscribe to transaction id events
type TxSubscribe struct {
	BaseMessage
	Transaction `json:"transaction"`
}

// Transaction bundles a single tx
type Transaction struct {
	Hash string `json:"hash"`
}

// AddressSubscribe is used to subscribe/unsubscribe to address events
type AddressSubscribe struct {
	BaseMessage
	Account `json:"account"`
}

// Account bundles a single account address
type Account struct {
	Address string `json:"address"`
}

// NewTxSubscribe constructs a Transaction subscription message
func NewTxSubscribe(msg BaseMessage, txHash string) TxSubscribe {
	msg.CategoryCode = "activeTransaction"
	msg.EventCode = "txSent"
	return TxSubscribe{
		BaseMessage: msg,
		Transaction: Transaction{Hash: txHash},
	}
}

// NewTxUnsubscribe constructs a Transaciton unsubscribe message
func NewTxUnsubscribe(msg BaseMessage, txHash string) TxSubscribe {
	msg.CategoryCode = "activeTransaction"
	msg.EventCode = "unwatch"
	return TxSubscribe{
		BaseMessage: msg,
		Transaction: Transaction{Hash: txHash},
	}
}

// NewAddressSubscribe constructs a address subscription message
func NewAddressSubscribe(msg BaseMessage, address string) AddressSubscribe {
	msg.CategoryCode = "accountAddress"
	msg.EventCode = "watch"
	return AddressSubscribe{
		BaseMessage: msg,
		Account:     Account{Address: address},
	}
}

// NewAddressUnsubscribe constructs a address unsubscribe message
func NewAddressUnsubscribe(msg BaseMessage, address string) AddressSubscribe {
	msg.CategoryCode = "accountAddress"
	msg.EventCode = "unwatch"
	return AddressSubscribe{
		BaseMessage: msg,
		Account:     Account{Address: address},
	}
}

// NewBaseMessageMainnet returns a base message suitable for mainnet usage
func NewBaseMessageMainnet(apiKey string) BaseMessage {
	if apiKey == "" {
		apiKey = os.Getenv("BLOCKNATIVE_API")
	}
	return BaseMessage{
		Timestamp: time.Now().String(),
		DappID:    apiKey,
		Version:   "v0",
		Blockchain: Blockchain{
			System:  "ethereum",
			Network: "main",
		},
	}
}
