package controllers

import (
	"errors"
	"testing"
	"time"

	"lab6/clients"
	"lab6/locations"

	"github.com/shopspring/decimal"
)

type controllerMockClient struct {
	currentTemp decimal.Decimal
	currentErr  error
	forecast    []clients.ForecastEntry
	forecastErr error
}

func (m *controllerMockClient) LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (decimal.Decimal, error) {
	return m.currentTemp, m.currentErr
}

func (m *controllerMockClient) LocationForecast(lat decimal.Decimal, lon decimal.Decimal, hours int) ([]clients.ForecastEntry, error) {
	return m.forecast, m.forecastErr
}

func newTestController(openWeather, google clients.WeatherDataClient) *WeatherController {
	return NewWeatherController(clients.NewWeatherClientFactory(openWeather, google))
}

func TestWeatherController_GetCurrentWeather(t *testing.T) {
	controller := newTestController(
		&controllerMockClient{currentTemp: decimal.NewFromInt(10)},
		&controllerMockClient{currentTemp: decimal.NewFromInt(20)},
	)

	result, err := controller.GetCurrentWeather(clients.ProviderGoogle, locations.Location{
		Name: "London",
		Lat:  decimal.NewFromFloat(51.5074),
		Lon:  decimal.NewFromFloat(-0.1278),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Temperature.Equal(decimal.NewFromInt(20)) {
		t.Fatalf("got %s, want 20", result.Temperature)
	}

	if result.Provider != "google" || result.Location != "London" {
		t.Fatalf("unexpected response: %+v", result)
	}
}

func TestWeatherController_GetCurrentWeatherProviderError(t *testing.T) {
	controller := NewWeatherController(clients.NewWeatherClientFactory(&controllerMockClient{}, nil))

	_, err := controller.GetCurrentWeather(clients.ProviderGoogle, locations.Location{
		Lat: decimal.NewFromInt(1),
		Lon: decimal.NewFromInt(1),
	})
	if err == nil {
		t.Fatalf("expected provider error")
	}
}

func TestWeatherController_GetCurrentWeatherClientError(t *testing.T) {
	controller := newTestController(
		&controllerMockClient{currentErr: errors.New("provider failed")},
		&controllerMockClient{},
	)

	_, err := controller.GetCurrentWeather(clients.ProviderOpenWeather, locations.Location{
		Lat: decimal.NewFromInt(1),
		Lon: decimal.NewFromInt(1),
	})
	if err == nil {
		t.Fatalf("expected client error")
	}
}

func TestWeatherController_GetForecast(t *testing.T) {
	forecastTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	controller := newTestController(
		&controllerMockClient{},
		&controllerMockClient{
			forecast: []clients.ForecastEntry{
				{Time: forecastTime, Temperature: decimal.NewFromFloat(3.5)},
			},
		},
	)

	result, err := controller.GetForecast(clients.ProviderGoogle, locations.Location{Name: "Tokyo"}, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 1 || !result.Items[0].Temperature.Equal(decimal.NewFromFloat(3.5)) {
		t.Fatalf("unexpected forecast: %+v", result)
	}
}

func TestWeatherController_GetForecastError(t *testing.T) {
	controller := newTestController(
		&controllerMockClient{forecastErr: errors.New("forecast failed")},
		&controllerMockClient{},
	)

	_, err := controller.GetForecast(clients.ProviderOpenWeather, locations.Location{}, 1)
	if err == nil {
		t.Fatalf("expected forecast error")
	}
}

func TestWeatherController_GetBatchCurrentWeather(t *testing.T) {
	controller := newTestController(
		&controllerMockClient{currentTemp: decimal.NewFromInt(5)},
		&controllerMockClient{},
	)

	result, err := controller.GetBatchCurrentWeather(clients.ProviderOpenWeather, []locations.Location{
		{Name: "Minsk"},
		{Name: "Warsaw", Lat: decimal.NewFromInt(1), Lon: decimal.NewFromInt(2)},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Locations) != 2 {
		t.Fatalf("got %d locations, want 2", len(result.Locations))
	}

	if !result.Locations[0].Temperature.Equal(decimal.NewFromInt(5)) {
		t.Fatalf("unexpected first temperature: %s", result.Locations[0].Temperature)
	}
}

func TestWeatherController_GetBatchCurrentWeatherPartialFailure(t *testing.T) {
	callCount := 0
	client := &controllerMockClient{
		currentTemp: decimal.NewFromInt(8),
		currentErr:  nil,
	}

	controller := NewWeatherController(clients.NewWeatherClientFactory(&batchFailClient{
		successClient: client,
		failOnSecond:  true,
		callCount:     &callCount,
	}, &controllerMockClient{}))

	result, err := controller.GetBatchCurrentWeather(clients.ProviderOpenWeather, []locations.Location{
		{Name: "Minsk"},
		{Name: "London"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Locations[0].Error != "" || result.Locations[1].Error == "" {
		t.Fatalf("expected partial failure, got %+v", result.Locations)
	}
}

type batchFailClient struct {
	successClient clients.WeatherDataClient
	failOnSecond  bool
	callCount     *int
}

func (c *batchFailClient) LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (decimal.Decimal, error) {
	*c.callCount++
	if c.failOnSecond && *c.callCount == 2 {
		return decimal.Zero, errors.New("second location failed")
	}

	return c.successClient.LocationCurrentTemperature(lat, lon)
}

func (c *batchFailClient) LocationForecast(lat decimal.Decimal, lon decimal.Decimal, hours int) ([]clients.ForecastEntry, error) {
	return c.successClient.LocationForecast(lat, lon, hours)
}

func TestCurrentWeatherController_GetCurrentWeather(t *testing.T) {
	controller := NewCurrentWeatherController(&controllerMockClient{currentTemp: decimal.NewFromInt(12)})

	result, err := controller.GetCurrentWeather(decimal.NewFromInt(1), decimal.NewFromInt(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.Temperature.Equal(decimal.NewFromInt(12)) {
		t.Fatalf("got %s, want 12", result.Temperature)
	}
}

func TestLocationDisplayNameUsesCoordinatesWhenNameMissing(t *testing.T) {
	name := locationDisplayName(locations.Location{
		Lat: decimal.NewFromInt(1),
		Lon: decimal.NewFromInt(2),
	})

	if name != "1,2" {
		t.Fatalf("got %q, want coordinate label", name)
	}
}
