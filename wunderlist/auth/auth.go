package auth

import (
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

// ref. https://developer.wunderlist.com/documentation/concepts/authorization
// example https://www.wunderlist.com/oauth/authorize?client_id=ID&redirect_uri=URL&state=RANDOM

// Client is wunderlist authorization
type Client struct {
	ClientID    string
	RedirectURI string
	State       string
}

// AuthorizeResponse is authorization response
type AuthorizeResponse struct {
	State string
	Code  string
}

const authorizationURL = "https://www.wunderlist.com/oauth/authorize"

// ?client_id=ID&redirect_uri=URL&state=RANDOM

// New Client
func New(clientID string, redirectURI string) (*Client, error) {
	randomID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Client{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		State:       randomID.String(),
	}, nil
}

func (c *Client) Redirect(w http.ResponseWriter, r *http.Request) error {
	urlStr, err := c.RequestURL()
	if err != nil {
		return nil
	}
	http.Redirect(w, r, urlStr, http.StatusFound)
	return nil
}

func (c *Client) RequestURL() (string, error) {
	u, err := url.Parse(authorizationURL)
	if err != nil {
		return "", nil
	}

	v := url.Values{}
	v.Add("client_id", c.ClientID)
	v.Add("redirect_uri", c.RedirectURI)
	v.Add("state", c.State)

	u.RawQuery = v.Encode()

	return u.String(), nil
}

func ParseRequest(r *http.Request) (*AuthorizeResponse, error) {
	return &AuthorizeResponse{
		Code:  r.FormValue("code"),
		State: r.FormValue("state"),
	}, nil
}
