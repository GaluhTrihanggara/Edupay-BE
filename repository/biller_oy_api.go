package repository

import (
	"Edupay/config"
	"Edupay/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type BillerOyApiRepository interface {
	BillInquryRepository(payload *model.OyBillerApi) (*model.OyBillerApiResponse, error)
	PayBillRepository(payload *model.OyBillerApi) (*model.OyBillerApiResponse, error)
	BillPaymentStatusRepository(payload *model.OyBillerApi) (*model.OyBillerApiResponse, error)
}

type billerOyApiRepository struct{}

func NewBillerOyApiOyApiRepository() BillerOyApiRepository {
	return &billerOyApiRepository{}
}

func (*billerOyApiRepository) BillInquryRepository(payload *model.OyBillerApi) (*model.OyBillerApiResponse, error) {
	url := fmt.Sprintf("%s/v2/bill", config.AppConfig.BaseUrl)
	resp, err := doRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}

	var response model.OyBillerApiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	return &response, nil
}

func (*billerOyApiRepository) PayBillRepository(payload *model.OyBillerApi) (*model.OyBillerApiResponse, error) {
	url := fmt.Sprintf("%s/v2/bill/payment", config.AppConfig.BaseUrl)
	resp, err := doRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}

	var response model.OyBillerApiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	return &response, nil
}

func (*billerOyApiRepository) BillPaymentStatusRepository(payload *model.OyBillerApi) (*model.OyBillerApiResponse, error) {
	url := fmt.Sprintf("%s/v2/bill/status", config.AppConfig.BaseUrl)
	resp, err := doRequest(http.MethodGet, url, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}

	var response model.OyBillerApiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	return &response, nil
}
