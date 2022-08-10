package asekopoolliveapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	authCookieName        = "connect.sid"
	loginEndpoint         = "/api/login"
	statusEndpoint        = "/api/units/%d/simple"
	graphEndpoint         = "/api/chart/%d"
	currentValuesEndpoint = "/api/units/%d"
)

type Client struct {
	conn     *resty.Client
	deviceID int32
}

func NewClient(debug bool) *Client {
	r := resty.New()

	r.SetBaseURL("https://pool.aseko.com")

	if debug {
		r.SetDebug(true)
	}

	return &Client{
		conn: r,
	}
}

func (c *Client) Login(username, password string, deviceID int32) error {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Agree    string `json:"agree"`
	}

	loginPayload := loginRequest{
		Username: username,
		Password: password,
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

	c.deviceID = deviceID

	return nil
}

type StatusData struct {
	SerialNumber string    `json:"serialNumber"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Timezone     string    `json:"timezone"`
	IsOnline     bool      `json:"isOnline"`
	DateLastData time.Time `json:"dateLastData"`
	HasError     bool      `json:"hasError"`
}

func (c *Client) Status() (*StatusData, error) {
	endpoint := fmt.Sprintf(statusEndpoint, c.deviceID)

	var status StatusData
	resp, err := c.conn.R().
		SetResult(&status).
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("non success response code receieved from api: %+v", err)
	}

	return &status, nil
}

type GraphData struct {
	IsLastest bool
	Variables []GraphVariable
	Items     []GraphItem
}

type graphResponse struct {
	IsLastest bool          `json:"isLastest"`
	Variables []interface{} `json:"variables"`
	Items     []interface{} `json:"items"`
}

type GraphVariable struct {
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Unit     string  `json:"unit"`
	Icon     string  `json:"icon"`
	Color    string  `json:"color"`
	Required float64 `json:"required,omitempty"`
}

type GraphItem struct {
	Timestamp      time.Time
	UnixTimestamp  int64   `json:"timestamp"`
	WaterTemp      float64 `json:"waterTemp,omitempty"`
	Ph             float64 `json:"ph,omitempty"`
	Rx             float64 `json:"rx,omitempty"`
	Offline        int     `json:"offline,omitempty"`
	PhNoCal        float64 `json:"phNoCal,omitempty"`
	WaterTempNoCal float64 `json:"waterTempNoCal,omitempty"`
	PhMinus        float64 `json:"phMinus,omitempty"`
	PhMinusValue   float64 `json:"phMinus-value,omitempty"`
	Cl             float64 `json:"cl,omitempty"`
	ClValue        float64 `json:"cl-value,omitempty"`
	NoWaterFlow    int     `json:"noWaterFlow,omitempty"`
}

func (c *Client) GetGraphData(startDate, endDate time.Time) (*GraphData, error) {
	endpoint := fmt.Sprintf(graphEndpoint, c.deviceID)

	if startDate.UnixMilli() == 0 || endDate.UnixMilli() == 0 {
		return nil, errors.New("invalid start date or end date provided")
	}

	resp, err := c.conn.R().
		SetQueryParams(map[string]string{
			"begin": fmt.Sprintf("%d", startDate.UnixMilli()),
			"end":   fmt.Sprintf("%d", endDate.UnixMilli()),
		}).Get(endpoint)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("non success response code receieved from api: %+v", err)
	}

	return parseGraphData(resp.String())
}

func parseGraphData(data string) (*GraphData, error) {
	var gr graphResponse
	err := json.Unmarshal([]byte(data), &gr)
	if err != nil {
		return nil, err
	}

	// Removing broken nodes in response
	variables := make([]interface{}, 0)
	for _, v := range gr.Variables {
		s, ok := v.(string)
		if ok {
			if s == "" {
				continue
			}
		}

		variables = append(variables, v)
	}

	gr.Variables = variables

	newDataString, err := json.Marshal(gr)
	if err != nil {
		return nil, err
	}

	var gp GraphData
	err = json.Unmarshal(newDataString, &gp)
	if err != nil {
		return nil, err
	}

	for index, item := range gp.Items {
		gp.Items[index].Timestamp = time.UnixMilli(item.UnixTimestamp)
	}

	return &gp, nil
}

type CurrentValues struct {
	Variables []struct {
		Type         string  `json:"type"`
		Name         string  `json:"name"`
		Unit         string  `json:"unit"`
		Icon         string  `json:"icon"`
		Color        string  `json:"color"`
		HasError     bool    `json:"hasError"`
		CurrentValue float64 `json:"currentValue"`
		Required     float64 `json:"required,omitempty"`
		Alarm        struct {
			Active   bool    `json:"active"`
			MinValue float64 `json:"minValue"`
			MaxValue float64 `json:"maxValue"`
		} `json:"alarm"`
	} `json:"variables"`
	Errors      []interface{} `json:"errors"`
	Info        []interface{} `json:"info"`
	ErrorsAlarm struct {
		Active   bool        `json:"active"`
		MinValue interface{} `json:"minValue"`
		MaxValue interface{} `json:"maxValue"`
	} `json:"errorsAlarm"`
}

func (c *Client) CurrentValues() (*CurrentValues, error) {
	endpoint := fmt.Sprintf(currentValuesEndpoint, c.deviceID)

	var currentValues CurrentValues

	resp, err := c.conn.R().
		SetResult(&currentValues).
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("non success response code receieved from api: %+v", err)
	}

	return &currentValues, nil
}
