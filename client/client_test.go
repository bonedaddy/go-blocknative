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

	require.NoError(t, client.Initialize(NewBaseMessageMainnet()))

	t.Log("sending subscribe message")
	addrSub := NewAddressSubscribe(NewBaseMessageMainnet(), "0xfa6de2697D59E88Ed7Fc4dFE5A33daC43565ea41")
	require.NoError(t, client.WriteJSON(addrSub))

	t.Log("reading message...")
	var out interface{}
	require.NoError(t, client.ReadJSON(&out))
	t.Log("message: ", out)
	require.NoError(t, client.Close())
}
