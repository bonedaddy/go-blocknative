package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	ctx := context.TODO()
	client, err := New(ctx, Opts{Scheme: "wss", Host: "api.blocknative.com", Path: "/v0", PrintConnectResponse: true})
	require.NoError(t, err)

	require.NoError(t, client.Initialize(BaseMessage{
		Timestamp: time.Now().String(),
		DappID:    os.Getenv("BLOCKNATIVE_API"),
		Version:   "v0",
		Blockchain: Blockchain{
			System:  "ethereum",
			Network: "main",
		},
	}))

	require.NoError(t, client.Close())
}
