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
		DosingChlor            float64 `json:"cl-value"`
		DosingPHMinus          float64 `json:"phMinus-value"`
		PH                     float64 `json:"ph"`
		PHNotCalibrated        float64 `json:"phNoCal"`
		Redox                  int     `json:"rx"`
		Timestamp              int64   `json:"timestamp"`
		WaterTemp              float64 `json:"waterTemp"`
		WaterTempNotCalibrated float64 `json:"waterTempNoCal"`
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
