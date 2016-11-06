package wlcli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"time"
)

var version = "0.0.1"
var endPoint = "https://a.wunderlist.com/api/v1"
var userAgent = fmt.Sprintf("wscli/%s (%s)", version, runtime.Version())

// Client for wunderList
type Client struct {
	URL        *url.URL
	HTTPClient *http.Client

	ClientID, ClientSecret string

	Logger *log.Logger
}

// User for wunderList
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Revision  int       `json:"revision"`
	Type      string    `json:"type"`
}

// NewClient construct client
func NewClient(clientid, clientsecret string) (*Client, error) {
	if len(clientid) == 0 {
		return nil, errors.New("missing clientid")
	}

	if len(clientsecret) == 0 {
		return nil, errors.New("missing clientsecret")
	}

	client := &Client{}
	client.URL, _ = url.Parse(endPoint)
	client.ClientID = clientid
	client.ClientSecret = clientsecret
	client.HTTPClient = new(http.Client)
	client.Logger = log.New(ioutil.Discard, "", log.LstdFlags)

	return client, nil

}

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("X-Client-ID", c.ClientID)
	req.Header.Set("X-Access-Token", c.ClientSecret)
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(out)
}

// GetUser returns User
func (c *Client) GetUser(ctx context.Context) (*User, error) {
	spath := "/user"
	req, err := c.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Check status code here…

	var user User
	if err := decodeBody(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
