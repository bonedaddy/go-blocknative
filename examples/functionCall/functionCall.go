package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/bonedaddy/go-blocknative/client"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/oklog/run"
	"github.com/pkg/errors"
)

func main() {
	ExitOnErr(godotenv.Load(), "loading .env file")

	logger := log.New(os.Stdout, "", log.Ltime|log.Lshortfile)

	var g run.Group

	// Run groups.
	{
		g.Add(run.SignalHandler(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM))

		mempMon, err := client.New(context.Background(), client.Opts{
			Scheme: "wss",
			Host:   "api.blocknative.com",
			Path:   "/v0",
			APIKey: os.Getenv("BLOCKNATIVE_DAPP_ID"),
		})

		ExitOnErr(err, "create blocknative client")

		baseMsg := client.BaseMessage{
			Timestamp: time.Now(),
			DappID:    os.Getenv("BLOCKNATIVE_DAPP_ID"),
			Version:   "v0",
			Blockchain: client.Blockchain{
				System:  "ethereum",
				Network: "rinkeby",
			},
		}

		ExitOnErr(mempMon.Initialize(baseMsg), "initialize subs")

		parsed, err := abi.JSON(strings.NewReader(TellorABI))
		ExitOnErr(err, "parsing contract ABI")

		cfgMsg := client.NewConfig(
			"0x88dF592F8eb5D7Bd38bFeF7dEb0fBc02cf3778a0",
			true,
			parsed,
			[]map[string]string{
				{
					"contractCall.methodName": "submitMiningSolution",
					"_propertySearch":         "true",
				},
			},
		)

		cfgMsgWithBase := client.NewConfiguration(baseMsg, cfgMsg)

		ExitOnErr(mempMon.EventSub(cfgMsgWithBase), "config subs")
		log.Print("subscription created")

		g.Add(func() error {
			for {
				msg := &client.EthTxPayload{}
				if err := mempMon.ReadJSON(msg); err != nil {
					if e, ok := err.(*websocket.CloseError); ok {
						if e.Code != 1000 {
							log.Fatal("mempMon read", err)
						}
					}
					return nil
				}
				log.Printf("msg: %+v \n", msg)
			}
		}, func(error) {
			mempMon.Close()
			log.Printf("closed")
		})

		if err := g.Run(); err != nil {
			logger.Println(fmt.Sprintf("%+v", errors.Wrapf(err, "run group stacktrace")))
		}

	}
}

func ExitOnErr(err error, msg string) {
	logger := log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
	if err != nil {
		logger.Output(2, fmt.Sprintf("root execution error:%+v msg:%+v", err, msg))
		os.Exit(1)
	}
}

const TellorABI = `[
    {
        "inputs": [
            {
                "internalType": "string",
                "name": "_nonce",
                "type": "string"
            },
            {
                "internalType": "uint256[5]",
                "name": "_requestId",
                "type": "uint256[5]"
            },
            {
                "internalType": "uint256[5]",
                "name": "_value",
                "type": "uint256[5]"
            }
        ],
        "name": "submitMiningSolution",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    }
]`
