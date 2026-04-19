package weather

import "github.com/shopspring/decimal"

type CurrentWeather struct {
	Location    string          `json:"location,omitempty"`
	Temperature decimal.Decimal `json:"temperature"`
	Provider    string          `json:"provider"`
}

type ForecastItem struct {
	Time        string          `json:"time"`
	Temperature decimal.Decimal `json:"temperature"`
}

type WeatherForecast struct {
	Location string         `json:"location,omitempty"`
	Provider string         `json:"provider"`
	Items    []ForecastItem `json:"items"`
}

type LocationWeather struct {
	Location    string          `json:"location"`
	Temperature decimal.Decimal `json:"temperature,omitempty"`
	Error       string          `json:"error,omitempty"`
}

type BatchCurrentWeather struct {
	Provider  string            `json:"provider"`
	Locations []LocationWeather `json:"locations"`
}
