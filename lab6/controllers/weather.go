package controllers

import (
	"lab6/clients"
	"lab6/locations"
	weather "lab6/models/weather"

	"github.com/shopspring/decimal"
)

type WeatherController struct {
	Factory *clients.WeatherClientFactory
}

func NewWeatherController(factory *clients.WeatherClientFactory) *WeatherController {
	return &WeatherController{Factory: factory}
}

func (c *WeatherController) GetCurrentWeather(provider clients.Provider, location locations.Location) (weather.CurrentWeather, error) {
	client, err := c.Factory.Get(provider)
	if err != nil {
		return weather.CurrentWeather{}, err
	}

	temperature, err := client.LocationCurrentTemperature(location.Lat, location.Lon)
	if err != nil {
		return weather.CurrentWeather{}, err
	}

	return weather.CurrentWeather{
		Location:    locationDisplayName(location),
		Temperature: temperature,
		Provider:    string(provider),
	}, nil
}

func (c *WeatherController) GetForecast(provider clients.Provider, location locations.Location, hours int) (weather.WeatherForecast, error) {
	client, err := c.Factory.Get(provider)
	if err != nil {
		return weather.WeatherForecast{}, err
	}

	entries, err := client.LocationForecast(location.Lat, location.Lon, hours)
	if err != nil {
		return weather.WeatherForecast{}, err
	}

	items := make([]weather.ForecastItem, 0, len(entries))
	for _, entry := range entries {
		items = append(items, weather.ForecastItem{
			Time:        entry.Time.Format("2006-01-02T15:04:05Z07:00"),
			Temperature: entry.Temperature,
		})
	}

	return weather.WeatherForecast{
		Location: locationDisplayName(location),
		Provider: string(provider),
		Items:    items,
	}, nil
}

func (c *WeatherController) GetBatchCurrentWeather(provider clients.Provider, locationsList []locations.Location) (weather.BatchCurrentWeather, error) {
	client, err := c.Factory.Get(provider)
	if err != nil {
		return weather.BatchCurrentWeather{}, err
	}

	result := weather.BatchCurrentWeather{
		Provider:  string(provider),
		Locations: make([]weather.LocationWeather, 0, len(locationsList)),
	}

	for _, location := range locationsList {
		item := weather.LocationWeather{Location: locationDisplayName(location)}

		temperature, err := client.LocationCurrentTemperature(location.Lat, location.Lon)
		if err != nil {
			item.Error = err.Error()
		} else {
			item.Temperature = temperature
		}

		result.Locations = append(result.Locations, item)
	}

	return result, nil
}

func locationDisplayName(location locations.Location) string {
	if location.Name != "" {
		return location.Name
	}

	return location.Lat.String() + "," + location.Lon.String()
}

// Deprecated: kept for backward compatibility in tests during migration.
type CurrentWeatherController[T clients.WeatherDataClient] struct {
	Client T
}

func NewCurrentWeatherController[T clients.WeatherDataClient](client T) *CurrentWeatherController[T] {
	return &CurrentWeatherController[T]{Client: client}
}

func (c *CurrentWeatherController[T]) GetCurrentWeather(lat decimal.Decimal, lon decimal.Decimal) (weather.CurrentWeather, error) {
	temperature, err := c.Client.LocationCurrentTemperature(lat, lon)
	if err != nil {
		return weather.CurrentWeather{}, err
	}

	return weather.CurrentWeather{
		Temperature: temperature,
	}, nil
}
