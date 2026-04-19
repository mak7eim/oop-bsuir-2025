package main

import (
	"lab6/api"
	_ "lab6/api/docs"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Weather Example API
// @version         1.0
// @BasePath        /api/v1

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	weatherHandler := api.NewCurrentWeatherHandler()

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/weather", weatherHandler.HandleGetCurrentWeather)
		v1.GET("/weather/forecast", weatherHandler.HandleGetForecast)
		v1.GET("/weather/batch", weatherHandler.HandleGetBatchCurrentWeather)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
