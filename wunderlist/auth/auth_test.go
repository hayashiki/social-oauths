package auth

import (
	"net/http"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	req, err := New("id", "http://localhost")

	if err != nil {
		t.Fatal(err)
	}

	if len(req.State) != 36 {
		t.Errorf("Expect state length 36, got %v", len(req.State))
	}
}

func TestClient_RequestURL(t *testing.T) {
	req, err := New("id", "http://localhost/callback")

	if err != nil {
		t.Fatal(err)
	}

	urlStr, err := req.RequestURL()
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		key         string
		expectValue string
	}{
		{"client_id", "id"},
		{"redirect_uri", "http://localhost/callback"},
		{"state", req.State},
	}

	for _, test := range tests {
		v := u.Query().Get(test.key)
		if test.expectValue != v {
			t.Errorf("Expect %v:%v, got %v", test.key, test.expectValue, v)
		}
	}
}

func Test_ParseAuthorize(t *testing.T) {
	tests := []struct {
		urlValues url.Values
	}{
		{
			url.Values{
				"code":  []string{"code_value"},
				"state": []string{"state_value"},
			},
		},
	}
	for _, test := range tests {
		req := &http.Request{Form: test.urlValues}
		resp, err := ParseRequest(req)

		if err != nil {
			t.Fatal(err)
		}

		if resp.Code != test.urlValues.Get("code") {
			t.Errorf("Expect code_value, got:%#v", test.urlValues.Get("code"))
		}

		if resp.State != test.urlValues.Get("state") {
			t.Errorf("Expect state_value, got:%#v", test.urlValues.Get("state"))
		}
	}
}
