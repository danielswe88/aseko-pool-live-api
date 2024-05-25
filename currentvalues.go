package asekopoolliveapi

import "fmt"

type CurrentValues struct {
	Variables []struct {
		Type         string `json:"type"`
		Name         string `json:"name"`
		Unit         string `json:"unit"`
		Icon         string `json:"icon"`
		Color        string `json:"color"`
		HasError     bool   `json:"hasError"`
		CurrentValue int    `json:"currentValue"`
		Required     int    `json:"required,omitempty"`
		Alarm        struct {
			Active   bool `json:"active"`
			MinValue int  `json:"minValue"`
			MaxValue int  `json:"maxValue"`
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
