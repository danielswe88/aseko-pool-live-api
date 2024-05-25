package asekopoolliveapi

import (
	"fmt"
	"time"
)

type Status struct {
	SerialNumber string    `json:"serialNumber"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Timezone     string    `json:"timezone"`
	IsOnline     bool      `json:"isOnline"`
	DateLastData time.Time `json:"dateLastData"`
	HasError     bool      `json:"hasError"`
}

func (c *Client) Status() (*Status, error) {
	endpoint := fmt.Sprintf(statusEndpoint, c.deviceID)

	var status Status
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
