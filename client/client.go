package client

import (
	"context"
	"errors"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Opts provides configuration over the websocket connection
type Opts struct {
	Scheme               string
	Host                 string
	Path                 string
	PrintConnectResponse bool
}

// ConnectResponse is the message we receive when opening a connection to the API
type ConnectResponse struct {
	ConnectionID  string `json:"connectionId"`
	ServerVersion string `json:"serverVersion"`
	ShowUX        bool   `json:"showUX"`
	Status        string `json:"status"`
	Version       int    `json:"version"`
}

// Client wraps gorilla websocket connections
type Client struct {
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

// New returns a new blocknative websocket client
func New(ctx context.Context, opts Opts) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)
	u := url.URL{
		Scheme: opts.Scheme,
		Host:   opts.Host,
		Path:   opts.Path,
	}
	c, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		cancel()
		return nil, err
	}
	// this checks out connection to blocknative's api and makes sure that we connected properly
	var out ConnectResponse
	if err := c.ReadJSON(&out); err != nil {
		cancel()
		return nil, err
	}
	if out.Status != "ok" {
		cancel()
		return nil, errors.New("failed to initialize websockets api connection")
	}
	if opts.PrintConnectResponse {
		log.Printf("%+v\n", out)
	}
	return &Client{conn: c, ctx: ctx, cancel: cancel}, nil
}

// Initialize is used to handle blocknative websockets api initialization
// note we set CategoryCode and EventCode ourselves
func (c *Client) Initialize(msg BaseMessage) error {
	msg.CategoryCode = "initialize"
	msg.EventCode = "checkDappId"
	return c.conn.WriteJSON(&msg)
}

// Close is used to terminate our websocket client
func (c *Client) Close() error {
	c.cancel()
	return c.conn.Close()
}
