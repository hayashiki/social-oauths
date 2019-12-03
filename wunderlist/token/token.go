package token

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	HTTPClient   *http.Client
	RedirectURL  string
	ClientID     string
	ClientSecret string
}

type request struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type response struct {
	AccessToken string `json:"access_token"`
}

type Option func(*Client)

func New(redirectURL string, clientID string, clientSecret string, opts ...Option) *Client {
	c := &Client{
		HTTPClient:   http.DefaultClient,
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	for _, o := range opts {
		o(c)
	}
	return c
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// POST https://www.wunderlist.com/oauth/access_token
// client_id	string	Required. The client ID you received from Wunderlist when you registered.
// client_secret	string	Required. The client secret you received from Wunderlist when you registered.
// code string	Required. The code you received as a response to Step 1.
func (c *Client) GetAccessToken(code string) (*response, error) {
	j, _ := json.Marshal(request{
		Code:         code,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
	})

	req, err := http.NewRequest("POST", "https://www.wunderlist.com/oauth/access_token", bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)

	// resp, err := c.HTTPClient.PostForm("https://www.wunderlist.com/oauth/access_token", v)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		response := &response{}
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, errors.Errorf("cannot get a access token:%v", resp.Status)
}
