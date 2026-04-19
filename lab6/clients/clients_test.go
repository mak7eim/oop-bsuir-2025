package clients

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

type mockWeatherClient struct {
	currentTemp decimal.Decimal
	currentErr  error
	forecast    []ForecastEntry
	forecastErr error
}

func (m *mockWeatherClient) LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (decimal.Decimal, error) {
	return m.currentTemp, m.currentErr
}

func (m *mockWeatherClient) LocationForecast(lat decimal.Decimal, lon decimal.Decimal, hours int) ([]ForecastEntry, error) {
	return m.forecast, m.forecastErr
}

func TestParseProvider(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Provider
		wantErr bool
	}{
		{name: "default openweather", input: "", want: ProviderOpenWeather},
		{name: "openweather explicit", input: "openweather", want: ProviderOpenWeather},
		{name: "google", input: "google", want: ProviderGoogle},
		{name: "unknown provider", input: "yahoo", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProvider(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWeatherClientFactory_Get(t *testing.T) {
	openWeather := &mockWeatherClient{currentTemp: decimal.NewFromInt(10)}
	google := &mockWeatherClient{currentTemp: decimal.NewFromInt(20)}
	factory := NewWeatherClientFactory(openWeather, google)

	client, err := factory.Get(ProviderGoogle)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	temp, err := client.LocationCurrentTemperature(decimal.Zero, decimal.Zero)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !temp.Equal(decimal.NewFromInt(20)) {
		t.Fatalf("got %s, want 20", temp)
	}
}

func TestWeatherClientFactory_GetUnknownProvider(t *testing.T) {
	factory := NewWeatherClientFactory(&mockWeatherClient{}, &mockWeatherClient{})

	_, err := factory.Get(Provider("unknown"))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOpenWeatherClient_LocationCurrentTemperature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("units") != "metric" {
			t.Fatalf("expected metric units")
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"main":{"temp":12.5}}`))
	}))
	defer server.Close()

	client := NewOpenWeatherClient("test-key", server.URL, server.URL)
	client.httpClient = server.Client()

	temp, err := client.LocationCurrentTemperature(decimal.NewFromFloat(53.9), decimal.NewFromFloat(27.5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !temp.Equal(decimal.NewFromFloat(12.5)) {
		t.Fatalf("got %s, want 12.5", temp)
	}
}

func TestOpenWeatherClient_LocationCurrentTemperatureBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewOpenWeatherClient("test-key", server.URL, server.URL)
	client.httpClient = server.Client()

	_, err := client.LocationCurrentTemperature(decimal.NewFromInt(1), decimal.NewFromInt(1))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOpenWeatherClient_LocationForecast(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("cnt") != "1" {
			t.Fatalf("expected cnt=1, got %s", r.URL.Query().Get("cnt"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"list":[{"dt":1700000000,"main":{"temp":5}},{"dt":1700010800,"main":{"temp":7}}]}`))
	}))
	defer server.Close()

	client := NewOpenWeatherClient("test-key", server.URL, server.URL)
	client.httpClient = server.Client()

	entries, err := client.LocationForecast(decimal.NewFromInt(1), decimal.NewFromInt(1), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}

	if !entries[0].Temperature.Equal(decimal.NewFromInt(5)) {
		t.Fatalf("unexpected first temperature: %s", entries[0].Temperature)
	}
}

func TestGoogleWeatherClient_LocationCurrentTemperature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("unitsSystem") != "METRIC" {
			t.Fatalf("expected METRIC units")
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"temperature":{"degrees":15.2,"unit":"CELSIUS"}}`))
	}))
	defer server.Close()

	client := NewGoogleWeatherClient("test-key", server.URL, server.URL)
	client.httpClient = server.Client()

	temp, err := client.LocationCurrentTemperature(decimal.NewFromFloat(51.5), decimal.NewFromFloat(-0.12))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !temp.Equal(decimal.NewFromFloat(15.2)) {
		t.Fatalf("got %s, want 15.2", temp)
	}
}

func TestGoogleWeatherClient_LocationForecast(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("hours") != "2" {
			t.Fatalf("expected hours=2, got %s", r.URL.Query().Get("hours"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"forecastHours":[{"interval":{"startTime":"2025-02-05T23:00:00Z"},"temperature":{"degrees":11.4,"unit":"CELSIUS"}}]}`))
	}))
	defer server.Close()

	client := NewGoogleWeatherClient("test-key", server.URL, server.URL)
	client.httpClient = server.Client()

	entries, err := client.LocationForecast(decimal.NewFromInt(1), decimal.NewFromInt(1), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}

	expectedTime := time.Date(2025, 2, 5, 23, 0, 0, 0, time.UTC)
	if !entries[0].Time.Equal(expectedTime) {
		t.Fatalf("unexpected time: %v", entries[0].Time)
	}
}

func TestGoogleWeatherClient_LocationCurrentTemperatureNetworkError(t *testing.T) {
	client := NewGoogleWeatherClient("test-key", "http://127.0.0.1:1", "http://127.0.0.1:1")

	_, err := client.LocationCurrentTemperature(decimal.NewFromInt(1), decimal.NewFromInt(1))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOpenWeatherClient_LocationForecastUsesDefaultHours(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"list":[]}`))
	}))
	defer server.Close()

	client := NewOpenWeatherClient("test-key", server.URL, server.URL)
	client.httpClient = server.Client()

	_, err := client.LocationForecast(decimal.NewFromInt(1), decimal.NewFromInt(1), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMockWeatherClientErrors(t *testing.T) {
	client := &mockWeatherClient{
		currentErr:  errors.New("current failed"),
		forecastErr: errors.New("forecast failed"),
	}

	if _, err := client.LocationCurrentTemperature(decimal.Zero, decimal.Zero); err == nil {
		t.Fatalf("expected current error")
	}

	if _, err := client.LocationForecast(decimal.Zero, decimal.Zero, 1); err == nil {
		t.Fatalf("expected forecast error")
	}
}
