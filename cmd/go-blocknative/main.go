package main

import (
	"log"
	"os"

	"github.com/bonedaddy/go-blocknative/client"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v2"
)

var (
	apiClient *client.Client
)

func main() {
	app := cli.NewApp()
	app.Name = "go-blocknative"
	app.Usage = "cli for interacting with blocknative api"
	app.Before = func(c *cli.Context) (err error) {
		apiClient, err = client.New(c.Context, client.Opts{
			Scheme: c.String("scheme"),
			Host:   c.String("host"),
			Path:   c.String("api.path"),
			APIKey: c.String("api.key"),
		})
		if err != nil {
			return
		}
		err = apiClient.Initialize(client.NewBaseMessageMainnet(c.String("api.key")))
		return
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "api.key",
			EnvVars: []string{"BLOCKNATIVE_API"},
			Usage:   "blocknative api key",
		},
		&cli.StringFlag{
			Name:  "address",
			Usage: "address to use when subscribing to events",
			Value: "0xfa6de2697D59E88Ed7Fc4dFE5A33daC43565ea41",
		},
		&cli.StringFlag{
			Name:  "tx.hash",
			Usage: "transaction hash to use when subscribing to events",
		},
		&cli.StringFlag{
			Name:  "scheme",
			Usage: "connection scheme to use",
			Value: "wss",
		},
		&cli.StringFlag{
			Name:  "host",
			Usage: "host to connect to",
			Value: "api.blocknative.com",
		},
		&cli.StringFlag{
			Name:  "api.path",
			Usage: "api path to use",
			Value: "/v0",
		},
	}
	app.Commands = cli.Commands{
		&cli.Command{
			Name:    "subscribe",
			Aliases: []string{"sub"},
			Usage:   "event subscription commands",
			Subcommands: cli.Commands{
				&cli.Command{
					Name:  "address",
					Usage: "subscribe to events based on addresse",
					Action: func(c *cli.Context) error {
						if err := apiClient.WriteJSON(client.NewAddressSubscribe(
							client.NewBaseMessageMainnet(
								c.String("api.key"),
							),
							c.String("address"),
						)); err != nil {
							return err
						}
						for {
							var out interface{}
							if err := apiClient.ReadJSON(&out); err != nil {
								// used to ignore the following event
								// websocket: close 1005 (no status)
								if websocket.IsUnexpectedCloseError(err, 1005) {
									log.Println("receive unexpected close, exiting: ", err)
									break
								} else {
									log.Println("receive expected close message: ", err)
									continue
								}
							}
							log.Printf("receive message:\n%+v\n", out)
						}
						defer apiClient.Close()
						return nil
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
