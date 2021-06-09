package client

import (
	"context"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Opts provides configuration over the websocket connection
type Opts struct {
	Scheme               string
	Host                 string
	Path                 string
	APIKey               string
	PrintConnectResponse bool
}

// ConnectResponse is the message we receive when opening a connection to the API
type ConnectResponse struct {
	ConnectionID  string `json:"connectionId"`
	ServerVersion string `json:"serverVersion"`
	ShowUX        bool   `json:"showUX"`
	Status        string `json:"status"`
	Reason        string `json:"reason"`
	Version       int    `json:"version"`
}

// Client wraps gorilla websocket connections
type Client struct {
	conn    *websocket.Conn
	ctx     context.Context
	cancel  context.CancelFunc
	initMsg BaseMessage // used to resend the initialization msg if connection drops
	apiKey  string
	mtx     sync.RWMutex
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
		return nil, errors.Errorf("failed to initialize websockets connection reason:", out.Reason)
	}
	if opts.PrintConnectResponse {
		log.Printf("%+v\n", out)
	}
	return &Client{conn: c, ctx: ctx, cancel: cancel, apiKey: opts.APIKey}, nil
}

// Initialize is used to handle blocknative websockets api initialization
// note we set CategoryCode and EventCode ourselves.
func (c *Client) Initialize(msg BaseMessage) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	msg.CategoryCode = "initialize"
	msg.EventCode = "checkDappId"
	c.initMsg = msg
	if err := c.conn.WriteJSON(&msg); err != nil {
		return err
	}
	var out ConnectResponse
	err := c.conn.ReadJSON(&out)
	if err != nil {
		return err
	}
	if out.Status != "ok" {
		return errors.Errorf("failed to initialize api connection reason:%v", out.Reason)
	}
	return nil
}

// EventSub creates an event subscription.
func (c *Client) EventSub(msg Configuration) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if err := c.conn.WriteJSON(&msg); err != nil {
		return err
	}

	var out ConnectResponse
	err := c.conn.ReadJSON(&out)
	if err != nil {
		return err
	}
	if out.Status != "ok" {
		return errors.Errorf("failed to create subscription reason:%v", out.Reason)
	}

	return nil
}

// ReadJSON is a wrapper around Conn:ReadJSON
func (c *Client) ReadJSON(out interface{}) error {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.conn.ReadJSON(out)
}

// WriteJSON is a wrapper around Conn:WriteJSON
func (c *Client) WriteJSON(out interface{}) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.conn.WriteJSON(out)
}

// APIKey returns the api key being used by the client
func (c *Client) APIKey() string {
	return c.apiKey
}

// Close is used to terminate our websocket client
func (c *Client) Close() error {
	err := c.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	c.cancel()
	return err
}
