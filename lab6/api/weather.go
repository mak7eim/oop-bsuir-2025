package api

import (
	"lab6/clients"
	"lab6/controllers"
	"lab6/locations"
	weather "lab6/models/weather"
	"lab6/shared/responses"
	"lab6/shared/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	Controller *controllers.WeatherController
}

func NewWeatherHandler(controller *controllers.WeatherController) *WeatherHandler {
	return &WeatherHandler{Controller: controller}
}

func NewCurrentWeatherHandler() *WeatherHandler {
	openWeather := clients.NewOpenWeatherClient(
		utils.GetEnv("OPENWEATHER_API_KEY", ""),
		utils.GetEnv("OPENWEATHER_BASE_URL", ""),
		utils.GetEnv("OPENWEATHER_FORECAST_URL", ""),
	)
	googleWeather := clients.NewGoogleWeatherClient(
		utils.GetEnv("GOOGLE_WEATHER_API_KEY", ""),
		utils.GetEnv("GOOGLE_WEATHER_CURRENT_URL", ""),
		utils.GetEnv("GOOGLE_WEATHER_FORECAST_URL", ""),
	)

	return NewWeatherHandler(controllers.NewWeatherController(
		clients.NewWeatherClientFactory(openWeather, googleWeather),
	))
}

// HandleGetCurrentWeather godoc
// @Summary      Get current weather
// @Description  Returns current weather for coordinates or a supported city
// @Tags         weather
// @Produce      json
// @Param        lat       query     string  false  "Latitude"
// @Param        lon       query     string  false  "Longitude"
// @Param        city      query     string  false  "City name (Minsk, London, Tokyo, Shanghai, Warsaw)"
// @Param        provider  query     string  false  "Weather provider"  Enums(openweather,google)  default(openweather)
// @Success      200       {object}  responses.SuccessResponse[weather.CurrentWeather]
// @Failure      400       {object}  responses.StatusResponse
// @Failure      500       {object}  responses.StatusResponse
// @Router       /weather [get]
func (h *WeatherHandler) HandleGetCurrentWeather(c *gin.Context) {
	provider, location, ok := h.parseProviderAndLocation(c)
	if !ok {
		return
	}

	result, err := h.Controller.GetCurrentWeather(provider, location)
	if err != nil {
		c.JSON(500, responses.StatusResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(200, responses.SuccessResponse[weather.CurrentWeather]{Code: 200, Message: "Success", Data: result})
}

// HandleGetForecast godoc
// @Summary      Get weather forecast
// @Description  Returns weather forecast for coordinates or a supported city
// @Tags         weather
// @Produce      json
// @Param        lat       query     string  false  "Latitude"
// @Param        lon       query     string  false  "Longitude"
// @Param        city      query     string  false  "City name (Minsk, London, Tokyo, Shanghai, Warsaw)"
// @Param        hours     query     int     false  "Forecast horizon in hours"  default(24)
// @Param        provider  query     string  false  "Weather provider"  Enums(openweather,google)  default(openweather)
// @Success      200       {object}  responses.SuccessResponse[weather.WeatherForecast]
// @Failure      400       {object}  responses.StatusResponse
// @Failure      500       {object}  responses.StatusResponse
// @Router       /weather/forecast [get]
func (h *WeatherHandler) HandleGetForecast(c *gin.Context) {
	provider, location, ok := h.parseProviderAndLocation(c)
	if !ok {
		return
	}

	hours, err := strconv.Atoi(c.DefaultQuery("hours", "24"))
	if err != nil || hours <= 0 {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: "invalid hours parameter"})
		return
	}

	result, err := h.Controller.GetForecast(provider, location, hours)
	if err != nil {
		c.JSON(500, responses.StatusResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(200, responses.SuccessResponse[weather.WeatherForecast]{Code: 200, Message: "Success", Data: result})
}

// HandleGetBatchCurrentWeather godoc
// @Summary      Get current weather for multiple cities
// @Description  Returns current temperature for several supported cities
// @Tags         weather
// @Produce      json
// @Param        cities    query     string  true   "Comma-separated city names"
// @Param        provider  query     string  false  "Weather provider"  Enums(openweather,google)  default(openweather)
// @Success      200       {object}  responses.SuccessResponse[weather.BatchCurrentWeather]
// @Failure      400       {object}  responses.StatusResponse
// @Failure      500       {object}  responses.StatusResponse
// @Router       /weather/batch [get]
func (h *WeatherHandler) HandleGetBatchCurrentWeather(c *gin.Context) {
	provider, err := clients.ParseProvider(c.DefaultQuery("provider", string(clients.ProviderOpenWeather)))
	if err != nil {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: err.Error()})
		return
	}

	citiesRaw := strings.TrimSpace(c.Query("cities"))
	if citiesRaw == "" {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: "cities parameter is required"})
		return
	}

	cityNames := splitCSV(citiesRaw)
	locationsList, err := locations.ResolveCities(cityNames)
	if err != nil {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: err.Error()})
		return
	}

	result, err := h.Controller.GetBatchCurrentWeather(provider, locationsList)
	if err != nil {
		c.JSON(500, responses.StatusResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(200, responses.SuccessResponse[weather.BatchCurrentWeather]{Code: 200, Message: "Success", Data: result})
}

func (h *WeatherHandler) parseProviderAndLocation(c *gin.Context) (clients.Provider, locations.Location, bool) {
	provider, err := clients.ParseProvider(c.DefaultQuery("provider", string(clients.ProviderOpenWeather)))
	if err != nil {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: err.Error()})
		return "", locations.Location{}, false
	}

	location, err := locations.ResolveLocation(c.Query("city"), c.Query("lat"), c.Query("lon"))
	if err != nil {
		c.JSON(400, responses.StatusResponse{Code: 400, Message: err.Error()})
		return "", locations.Location{}, false
	}

	return provider, location, true
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
