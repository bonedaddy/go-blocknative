package client

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
