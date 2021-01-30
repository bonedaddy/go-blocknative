package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	ctx := context.TODO()
	client, err := New(ctx, Opts{Scheme: "wss", Host: "api.blocknative.com", Path: "/v0", PrintConnectResponse: true})
	require.NoError(t, err)

	// test base message creation deriving the api key from an environment variable
	require.NoError(t, client.Initialize(NewBaseMessageMainnet("")))

	t.Log("sending subscribe message")
	// test base message creation using the api key embedded into the client struct
	addrSub := NewAddressSubscribe(NewBaseMessageMainnet(client.APIKey()), "0xfa6de2697D59E88Ed7Fc4dFE5A33daC43565ea41")
	require.NoError(t, client.WriteJSON(addrSub))
	t.Log("reading message...")
	var out interface{}
	require.NoError(t, client.ReadJSON(&out))
	t.Log(out)
	require.NoError(t, client.WriteJSON(
		NewConfiguration(
			NewBaseMessageMainnet(
				client.APIKey(),
			),
			NewConfig(
				"0xfa6de2697D59E88Ed7Fc4dFE5A33daC43565ea41",
				false,
				[]string{logSwapABI},
				nil,
			),
		)),
	)
	require.NoError(t, client.ReadJSON(&out))
	t.Log(out)
	require.NoError(t, client.Close())
}

var (
	logSwapABI = `{
		"anonymous": false,
		"inputs": [
		  {
			"indexed": true,
			"internalType": "address",
			"name": "caller",
			"type": "address"
		  },
		  {
			"indexed": true,
			"internalType": "address",
			"name": "tokenIn",
			"type": "address"
		  },
		  {
			"indexed": true,
			"internalType": "address",
			"name": "tokenOut",
			"type": "address"
		  },
		  {
			"indexed": false,
			"internalType": "uint256",
			"name": "tokenAmountIn",
			"type": "uint256"
		  },
		  {
			"indexed": false,
			"internalType": "uint256",
			"name": "tokenAmountOut",
			"type": "uint256"
		  }
		],
		"name": "LOG_SWAP",
		"type": "event"
	}`
)
