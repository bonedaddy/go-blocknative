package client

import (
	"context"
	"log"
	"net/url"
	"sync"

	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
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
	Version       int    `json:"version"`
}

// Client wraps gorilla websocket connections
type Client struct {
	mx      sync.RWMutex
	conn    *websocket.Conn
	ctx     context.Context
	cancel  context.CancelFunc
	initMsg BaseMessage // used to resend the initialization msg if connection drops
	opts    Opts
	history *MsgHistory
}

// New returns a new blocknative websocket client caller must make sure to initialize afterwards
func New(ctx context.Context, opts Opts) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)
	client := &Client{
		ctx:     ctx,
		cancel:  cancel,
		opts:    opts,
		history: &MsgHistory{},
	}
	if err := client.doConnect(ctx, url.URL{
		Scheme: opts.Scheme,
		Host:   opts.Host,
		Path:   opts.Path,
	}); err != nil {
		cancel()
		return nil, err
	}
	return client, nil
}

// Initialize is used to handle blocknative websockets api initialization
// note we set CategoryCode and EventCode ourselves
func (c *Client) Initialize(msg BaseMessage) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	msg.CategoryCode = "initialize"
	msg.EventCode = "checkDappId"
	c.initMsg = msg
	if err := c.conn.WriteJSON(&msg); err != nil {
		return err
	}
	var out interface{}
	_ = c.conn.ReadJSON(&out)
	return nil
}

// ReadJSON is a wrapper around Conn:ReadJSON
// You should provide a pointer otherwise you will likely
// encounter a nil interface type as the returned value
func (c *Client) ReadJSON(out interface{}) error {
	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.conn.ReadJSON(out)
}

// WriteJSON is a wrapper around Conn:WriteJSON
// Do not provide a pointer as this could cause problems
// with the message history buffer if the provided value
// becomes garbage collected
func (c *Client) WriteJSON(out interface{}) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	// push the message into the history buffer
	c.history.Push(out)
	return c.conn.WriteJSON(out)
}

// APIKey returns the api key being used by the client
func (c *Client) APIKey() string {
	return c.opts.APIKey
}

// Close is used to terminate our websocket client
func (c *Client) Close() error {
	c.mx.Lock()
	defer c.mx.Unlock()
	err := c.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	)
	if err != nil {
		log.Println("failed to send close message: ", err)
	}
	// close the underlying connection
	c.conn.Close()
	c.cancel()
	return err
}

// ReInit should only be used in the event that we receive an unexpected
// error and allows us to replay previous messages
func (c *Client) ReInit() error {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.doConnect(c.ctx, url.URL{
		Scheme: c.opts.Scheme,
		Host:   c.opts.Host,
		Path:   c.opts.Path,
	})
	// dont empty the buffer such that future errors
	// can reuse the message history
	msgs := c.history.CopyAll()
	// send the initialize messsage
	// we do not store this in the message history buffer
	if err := c.conn.WriteJSON(c.initMsg); err != nil {
		return errors.Wrap(err, "fatal error received")
	}
	// drain
	_ = c.conn.ReadJSON(nil)
	for _, msg := range msgs {
		if err := c.conn.WriteJSON(&msg); err != nil {
			// TODO(bonedaddy): figure out how to properly handle
			log.Println("receive error during reinitialization: ", err)
			return err
		}
		// drain
		_ = c.conn.ReadJSON(nil)
	}
	return nil
}

// ShouldReInit is used to check the given error
// and return whether or not we should reinitialize the connection
func (c *Client) ShouldReInit(err error) bool {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
		return false
	}
	return true
}

// doConnect should only be used during creation of the initial client object or during reinitialization
// caller must take care of locking
func (c *Client) doConnect(ctx context.Context, u url.URL) error {
	// close the previous connection if it exists
	if c.conn != nil {
		c.conn.Close()
	}
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return err
	}
	c.conn = conn
	// this checks out connection to blocknative's api and makes sure that we connected properly
	var out ConnectResponse
	if err := c.conn.ReadJSON(&out); err != nil {
		return err
	}
	if out.Status != "ok" {
		return errors.New("failed to initialize websockets api connection")
	}
	if c.opts.PrintConnectResponse {
		log.Printf("%+v\n", out)
	}
	return nil
}
