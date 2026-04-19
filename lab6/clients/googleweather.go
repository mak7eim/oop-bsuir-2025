package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type googleTemperature struct {
	Degrees decimal.Decimal `json:"degrees"`
	Unit    string          `json:"unit"`
}

type googleCurrentResponse struct {
	Temperature googleTemperature `json:"temperature"`
}

type googleForecastResponse struct {
	ForecastHours []struct {
		Interval struct {
			StartTime string `json:"startTime"`
		} `json:"interval"`
		Temperature googleTemperature `json:"temperature"`
	} `json:"forecastHours"`
}

type GoogleWeatherClient struct {
	apiKey      string
	currentURL  string
	forecastURL string
	httpClient  *http.Client
}

func NewGoogleWeatherClient(apiKey, currentURL, forecastURL string) *GoogleWeatherClient {
	return &GoogleWeatherClient{
		apiKey:      apiKey,
		currentURL:  currentURL,
		forecastURL: forecastURL,
		httpClient:  http.DefaultClient,
	}
}

func (c *GoogleWeatherClient) LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (decimal.Decimal, error) {
	url := fmt.Sprintf("%s?key=%s&location.latitude=%s&location.longitude=%s&unitsSystem=METRIC",
		c.currentURL, c.apiKey, lat.String(), lon.String())

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to call google weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("google weather returned bad status: %d", resp.StatusCode)
	}

	var data googleCurrentResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return decimal.Zero, fmt.Errorf("failed to decode response: %w", err)
	}

	return data.Temperature.Degrees, nil
}

func (c *GoogleWeatherClient) LocationForecast(lat decimal.Decimal, lon decimal.Decimal, hours int) ([]ForecastEntry, error) {
	if hours <= 0 {
		hours = 24
	}
	if hours > 240 {
		hours = 240
	}

	url := fmt.Sprintf("%s?key=%s&location.latitude=%s&location.longitude=%s&unitsSystem=METRIC&hours=%d",
		c.forecastURL, c.apiKey, lat.String(), lon.String(), hours)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call google weather forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google weather forecast returned bad status: %d", resp.StatusCode)
	}

	var data googleForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	entries := make([]ForecastEntry, 0, len(data.ForecastHours))
	for _, item := range data.ForecastHours {
		parsedTime, err := time.Parse(time.RFC3339, item.Interval.StartTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse forecast time: %w", err)
		}

		entries = append(entries, ForecastEntry{
			Time:        parsedTime.UTC(),
			Temperature: item.Temperature.Degrees,
		})
	}

	return entries, nil
}
