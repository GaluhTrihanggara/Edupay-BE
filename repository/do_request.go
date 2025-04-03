package repository

import (
	"Edupay/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func doRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", echo.MIMEApplicationJSON)
	req.Header.Add("Accept", echo.MIMEApplicationJSON)
	req.Header.Add("x-oy-username", config.AppConfig.Username)
	req.Header.Add("x-api-key", config.AppConfig.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
