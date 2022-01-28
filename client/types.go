package client

import (
	"os"
	"time"
)

// BaseMessage is the base message required for all interactions with the websockets api
type BaseMessage struct {
	CategoryCode string    `json:"categoryCode"`
	EventCode    string    `json:"eventCode"`
	Timestamp    time.Time `json:"timeStamp"`
	DappID       string    `json:"dappId"` // api key
	Version      string    `json:"version"`
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

// EthTxPayload is payload returned from a subscription to blocknative api
type EthTxPayload struct {
	Version       int       `json:"version"`
	ServerVersion string    `json:"serverVersion"`
	TimeStamp     time.Time `json:"timeStamp"`
	ConnectionID  string    `json:"connectionId"`
	Status        string    `json:"status"`
	Event         struct {
		BaseMessage
		Transaction struct {
			Type                 int       `json:"type"`
			MaxFeePerGas         string    `json:"maxFeePerGas"`
			MaxPriorityFeePerGas string    `json:"maxPriorityFeePerGas"`
			BaseFeePerGas        string    `json:"baseFeePerGas"`
			TimeStamp            time.Time `json:"timeStamp"`
			Status               string    `json:"status"`
			MonitorID            string    `json:"monitorId"`
			MonitorVersion       string    `json:"monitorVersion"`
			TimePending          string    `json:"timePending"`
			PendingTimeStamp     time.Time `json:"pendingTimeStamp"`
			BlocksPending        int       `json:"blocksPending"`
			Hash                 string    `json:"hash"`
			From                 string    `json:"from"`
			To                   string    `json:"to"`
			Value                string    `json:"value"`
			Gas                  int       `json:"gas"`
			GasPrice             string    `json:"gasPrice"`
			GasPriceGwei         float64   `json:"gasPriceGwei"`
			Nonce                int       `json:"nonce"`
			BlockHash            string    `json:"blockHash"`
			BlockNumber          int       `json:"blockNumber"`
			TransactionIndex     int       `json:"transactionIndex"`
			Input                string    `json:"input"`
			GasUsed              string    `json:"gasUsed"`
			Asset                string    `json:"asset"`
			WatchedAddress       string    `json:"watchedAddress"`
			Direction            string    `json:"direction"`
			Counterparty         string    `json:"counterparty"`
		} `json:"transaction"`
	} `json:"event"`
}

// Configuration enables configuration of the blocknative websockets api
// and wraps the Config type
type Configuration struct {
	BaseMessage
	Config `json:"config"`
}

// Config provides a specific config instance
type Config struct {
	//  valid Ethereum address or 'global'
	Scope string `json:"scope"`
	// A slice of valid filters (jsql: https://github.com/deitch/searchjs)
	Filters []map[string]string `json:"filters,omitempty"`
	// JSON abis
	ABI interface{} `json:"abi,omitempty"`
	// defines whether the service should automatically watch the address as defined in
	WatchAddress bool `json:"watchAddress,omitempty"`
}

// NewConfig returns a new config instance
func NewConfig(scope string, watchAddress bool, abis interface{}) Config {
	cfg := Config{
		Scope:        scope,
		WatchAddress: watchAddress,
	}
	if abis != nil {
		cfg.ABI = abis
	}

	return cfg
}

// NewConfiguration constructs a new configuration message
func NewConfiguration(msg BaseMessage, config Config) Configuration {
	msg.CategoryCode = "configs"
	msg.EventCode = "put"
	return Configuration{
		BaseMessage: msg,
		Config:      config,
	}
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
		apiKey = os.Getenv("BLOCKNATIVE_DAPP_ID")
	}
	return BaseMessage{
		Timestamp: time.Now(),
		DappID:    apiKey,
		Blockchain: Blockchain{
			System:  "ethereum",
			Network: "main",
		},
	}
}
