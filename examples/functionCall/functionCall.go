package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/bonedaddy/go-blocknative/client"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gorilla/websocket"
	"github.com/oklog/run"
	"github.com/pkg/errors"
)

const (
	netName      = "main"
	contractAddr = "0x361cd36de2ffe3167904c58a5b0b22cf9217e466"
	methodName   = "submitMiningSolution"
)

func main() {
	// ExitOnErr(godotenv.Load(), "loading .env file")

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
			Blockchain: client.Blockchain{
				System:  "ethereum",
				Network: netName,
			},
		}

		ExitOnErr(mempMon.Initialize(baseMsg), "initialize subs")

		var abi interface{}
		ExitOnErr(json.Unmarshal([]byte(TellorABI), &abi), "marshal abi")

		cfgMsg := client.NewConfig(
			contractAddr,
			true,
			abi,
		)
		cfgMsg.Filters = []map[string]string{
			{
				"contractCall.methodName": methodName,
				"_propertySearch":         "true",
			},
		}

		cfgMsgWithBase := client.NewConfiguration(baseMsg, cfgMsg)

		msg, err := json.Marshal(cfgMsgWithBase)
		ExitOnErr(err, "config message marshal")
		log.Println("cfgMsgWithBase", string(msg))

		ExitOnErr(mempMon.EventSub(cfgMsgWithBase), "config subs")
		log.Print("subscription created   ", "network:", netName, "   contract:", contractAddr, "    method:", methodName)

		g.Add(func() error {
			for {
				msg := &client.EthTxPayload{}
				if err := mempMon.ReadJSON(msg); err != nil {
					if e, ok := err.(*websocket.CloseError); ok {
						if e.Code != 1000 {
							log.Fatal("mempMon read", err)
						}
					}
					return err
				}
				log.Printf("msg: %+v \n", msg)
				log.Printf("func args: %+v \n", parseInput(msg.Event.Transaction.Input))

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

func parseInput(input string) interface{} {
	abiT, err := abi.JSON(strings.NewReader(TellorABI))
	ExitOnErr(err, "loading the abi")

	inputData, err := hex.DecodeString(input[10:])
	ExitOnErr(err, "input decode")

	method, exist := abiT.Methods[methodName]
	if !exist {
		ExitOnErr(errors.New("method doesn't exists in the abi"), "")
	}

	output, err := method.Inputs.Unpack(inputData)
	ExitOnErr(err, "args unpack")

	return output
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
            },
            {
                "internalType": "uint256",
                "name": "_pass",
                "type": "uint256"
            }
        ],
        "name": "submitMiningSolution",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    }
]`
