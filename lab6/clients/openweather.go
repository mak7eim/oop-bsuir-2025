package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type openWeatherResponse struct {
	Main struct {
		Temp decimal.Decimal `json:"temp"`
	} `json:"main"`
}

type openWeatherForecastResponse struct {
	List []struct {
		Dt   int64 `json:"dt"`
		Main struct {
			Temp decimal.Decimal `json:"temp"`
		} `json:"main"`
	} `json:"list"`
}

type OpenWeatherClient struct {
	apiKey            string
	currentURL        string
	forecastURL       string
	httpClient        *http.Client
}

func NewOpenWeatherClient(apiKey, currentURL, forecastURL string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey:      apiKey,
		currentURL:  currentURL,
		forecastURL: forecastURL,
		httpClient:  http.DefaultClient,
	}
}

func (c *OpenWeatherClient) LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (decimal.Decimal, error) {
	url := fmt.Sprintf("%s?lat=%s&lon=%s&appid=%s&units=metric",
		c.currentURL, lat.String(), lon.String(), c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to call openweather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("openweather returned bad status: %d", resp.StatusCode)
	}

	var data openWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return decimal.Zero, fmt.Errorf("failed to decode response: %w", err)
	}

	return data.Main.Temp, nil
}

func (c *OpenWeatherClient) LocationForecast(lat decimal.Decimal, lon decimal.Decimal, hours int) ([]ForecastEntry, error) {
	if hours <= 0 {
		hours = 24
	}

	limit := (hours + 2) / 3
	if limit < 1 {
		limit = 1
	}

	url := fmt.Sprintf("%s?lat=%s&lon=%s&appid=%s&units=metric&cnt=%d",
		c.forecastURL, lat.String(), lon.String(), c.apiKey, limit)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call openweather forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openweather forecast returned bad status: %d", resp.StatusCode)
	}

	var data openWeatherForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	entries := make([]ForecastEntry, 0, len(data.List))
	for _, item := range data.List {
		entries = append(entries, ForecastEntry{
			Time:        time.Unix(item.Dt, 0).UTC(),
			Temperature: item.Main.Temp,
		})
	}

	return entries, nil
}
