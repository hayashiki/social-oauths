package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hayashiki/social-oauths/auth"
	"github.com/hayashiki/social-oauths/token"
)

var (
	BaseURL      = "http://localhost:3000"
	ClientID     = ""
	ClientSecret = ""
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	c, err := auth.New(ClientID, BaseURL+"/callback")
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(60 * time.Second)})

	c.Redirect(w, r)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	resp, err := auth.ParseRequest(r)
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	state, err := r.Cookie("state")

	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	if resp.State != state.Value {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	c := token.New(BaseURL+"/callback", ClientID, ClientSecret)
	response, err := c.GetAccessToken(resp.Code)

	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	fmt.Fprintf(w, "accessToken:%v", response.AccessToken)
}

func main() {
	http.HandleFunc("/auth", Authorize)
	http.HandleFunc("/callback", Callback)
	http.ListenAndServe(":3000", nil)
}
