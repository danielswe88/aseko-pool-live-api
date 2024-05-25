package asekopoolliveapi

import (
	"errors"
	"fmt"
	"time"
)

type Chart struct {
	Variables struct {
		Actuals []struct {
			Type     string  `json:"type"`
			Name     string  `json:"name"`
			Unit     string  `json:"unit"`
			Icon     string  `json:"icon"`
			Color    string  `json:"color"`
			Required float64 `json:"required"`
		} `json:"actuals"`
		CommonEvents []struct {
			Type  string `json:"type"`
			Name  string `json:"name"`
			Unit  string `json:"unit"`
			Icon  string `json:"icon"`
			Color string `json:"color"`
		} `json:"commonEvents"`
		ModeEvents   []interface{} `json:"modeEvents"`
		DosingEvents []struct {
			Type  string `json:"type"`
			Name  string `json:"name"`
			Unit  string `json:"unit"`
			Icon  string `json:"icon"`
			Color string `json:"color"`
		} `json:"dosingEvents"`
		OtherEvents []interface{} `json:"otherEvents"`
		RelayEvents []interface{} `json:"relayEvents"`
	} `json:"variables"`
	IsLastest bool `json:"isLastest"`
	Items     []struct {
		Timestamp      int64   `json:"timestamp"`
		PH             float64 `json:"ph"`
		WaterTemp      float64 `json:"waterTemp"`
		Redox          int     `json:"rx"`
		WaterTempNoCal float64 `json:"waterTempNoCal,omitempty"`
		PhNoCal        float64 `json:"phNoCal,omitempty"`
	} `json:"items"`
}

func (c *Client) Chart(startDate, endDate time.Time) (*Chart, error) {
	endpoint := fmt.Sprintf(chartEndpoint, c.deviceID)

	if startDate.UnixMilli() == 0 || endDate.UnixMilli() == 0 {
		return nil, errors.New("invalid start date or end date provided")
	}

	var data Chart
	resp, err := c.conn.R().
		SetQueryParams(map[string]string{
			"begin": fmt.Sprintf("%d", startDate.UnixMilli()),
			"end":   fmt.Sprintf("%d", endDate.UnixMilli()),
		}).SetResult(&data).Get(endpoint)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("non success response code receieved from api: %+v", err)
	}

	return &data, nil
}
