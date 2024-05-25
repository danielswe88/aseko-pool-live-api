package asekopoolliveapi

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const (
	authCookieName        = "connect.sid"
	loginEndpoint         = "/api/login"
	statusEndpoint        = "/api/units/%s/simple"
	chartEndpoint         = "/api/chart/%s"
	currentValuesEndpoint = "/api/units/%s"
)

type Client struct {
	conn *resty.Client

	username string
	password string
	deviceID string
}

func NewClient(username, password, deviceID string, debug bool) (*Client, error) {
	r := resty.New()

	r.SetBaseURL("https://pool.aseko.com")
	r.SetDebug(debug)
	r.SetHeader("User-Agent", "Mozilla/5.0 (iPhone; U; CPU iPhone OS 11.4.1; en_US; ) AppleWebKit/0.0 (KHTML, like Gecko) Version/0.0; GmmClient:google_ios/com.google.Maps/4.54.8/Mobile/ios:iPhone10,6/iOS-AppStore")

	c := Client{
		username: username,
		password: password,
		deviceID: deviceID,
		conn:     r,
	}

	err := c.login()
	if err != nil {
		return nil, fmt.Errorf("logging in: %v", err)
	}

	return &c, nil
}

func (c *Client) login() error {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Agree    string `json:"agree"`
	}

	loginPayload := loginRequest{
		Username: c.username,
		Password: c.password,
		Agree:    "on",
	}

	resp, err := c.conn.R().
		SetBody(loginPayload).
		Post(loginEndpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("status code 'unauthorized' receieved from api: %+v", err)
	}

	if resp.StatusCode() == http.StatusForbidden {
		return fmt.Errorf("status code 'forbidden' receieved from api: %+v", err)
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("non success response code receieved from api: %+v", err)
	}

	cookieFound := false
	for _, cookie := range resp.Cookies() {
		if cookie.Name == authCookieName {
			c.conn.SetCookie(cookie)
			cookieFound = true
		}
	}

	if !cookieFound {
		return fmt.Errorf("auth cookie not found in response")
	}

	return nil
}
