package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"lab6/clients"
	"lab6/controllers"
	weather "lab6/models/weather"
	"lab6/shared/responses"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type handlerMockClient struct {
	currentTemp decimal.Decimal
	currentErr  error
	forecast    []clients.ForecastEntry
	forecastErr error
}

func (m *handlerMockClient) LocationCurrentTemperature(lat decimal.Decimal, lon decimal.Decimal) (decimal.Decimal, error) {
	return m.currentTemp, m.currentErr
}

func (m *handlerMockClient) LocationForecast(lat decimal.Decimal, lon decimal.Decimal, hours int) ([]clients.ForecastEntry, error) {
	return m.forecast, m.forecastErr
}

func setupTestRouter(client clients.WeatherDataClient) *gin.Engine {
	gin.SetMode(gin.TestMode)

	handler := NewWeatherHandler(controllers.NewWeatherController(
		clients.NewWeatherClientFactory(client, client),
	))

	r := gin.New()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/weather", handler.HandleGetCurrentWeather)
		v1.GET("/weather/forecast", handler.HandleGetForecast)
		v1.GET("/weather/batch", handler.HandleGetBatchCurrentWeather)
	}

	return r
}

func TestHandleGetCurrentWeatherByCity(t *testing.T) {
	router := setupTestRouter(&handlerMockClient{currentTemp: decimal.NewFromInt(7)})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather?city=London&provider=openweather", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, body %s", rec.Code, rec.Body.String())
	}

	var response responses.SuccessResponse[weather.CurrentWeather]
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Data.Location != "London" || !response.Data.Temperature.Equal(decimal.NewFromInt(7)) {
		t.Fatalf("unexpected data: %+v", response.Data)
	}
}

func TestHandleGetCurrentWeatherInvalidLocation(t *testing.T) {
	router := setupTestRouter(&handlerMockClient{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather?city=Paris", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d, want 400", rec.Code)
	}
}

func TestHandleGetCurrentWeatherInvalidProvider(t *testing.T) {
	router := setupTestRouter(&handlerMockClient{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather?city=London&provider=yahoo", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d, want 400", rec.Code)
	}
}

func TestHandleGetForecast(t *testing.T) {
	router := setupTestRouter(&handlerMockClient{
		forecast: []clients.ForecastEntry{{Temperature: decimal.NewFromInt(4)}},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather/forecast?city=Tokyo&hours=3", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, body %s", rec.Code, rec.Body.String())
	}
}

func TestHandleGetBatchCurrentWeather(t *testing.T) {
	router := setupTestRouter(&handlerMockClient{currentTemp: decimal.NewFromInt(9)})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather/batch?cities=Minsk,Warsaw", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, body %s", rec.Code, rec.Body.String())
	}

	var response responses.SuccessResponse[weather.BatchCurrentWeather]
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response.Data.Locations) != 2 {
		t.Fatalf("got %d locations, want 2", len(response.Data.Locations))
	}
}

func TestHandleGetBatchCurrentWeatherMissingCities(t *testing.T) {
	router := setupTestRouter(&handlerMockClient{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/weather/batch", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d, want 400", rec.Code)
	}
}
