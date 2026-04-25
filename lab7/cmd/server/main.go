package main

import (
	"context"
	"log"
	"os"
	"time"

	"lab7-oop/api"
	"lab7-oop/api/handlers"
	"lab7-oop/api/middleware"
	"lab7-oop/clients/auth"
	"lab7-oop/clients/sqlite"
	"lab7-oop/controllers"
)

// @title           Food Delivery API
// @version         1.0
// @description     Backend for food delivery system (Clean Architecture lab).
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description     JWT token. Format: Bearer {token}
func main() {
	dbPath := envOr("DB_PATH", "food_delivery.db")
	jwtSecret := envOr("JWT_SECRET", "lab7-dev-secret-change-me")
	addr := envOr("ADDR", ":8080")

	db, err := sqlite.Open(dbPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	hasher := auth.BcryptHasher{}
	hash, _ := hasher.Hash("admin123")
	if err := db.EnsureAdmin(ctx, "admin@food.local", hash); err != nil {
		log.Fatalf("seed admin: %v", err)
	}
	if err := db.SeedDemoData(ctx); err != nil {
		log.Fatalf("seed data: %v", err)
	}

	userRepo := sqlite.NewUserRepository(db)
	menuRepo := sqlite.NewMenuRepository(db)
	cartRepo := sqlite.NewCartRepository(db)
	orderRepo := sqlite.NewOrderRepository(db)
	courierRepo := sqlite.NewCourierRepository(db)

	tokenSvc := auth.NewJWTService(jwtSecret, 24*time.Hour)

	authCtrl := controllers.NewAuthController(userRepo, tokenSvc, hasher, courierRepo)
	menuCtrl := controllers.NewMenuController(menuRepo)
	cartCtrl := controllers.NewCartController(cartRepo, menuRepo)
	orderCtrl := controllers.NewOrderController(orderRepo, cartCtrl, courierRepo)
	courierCtrl := controllers.NewCourierController(courierRepo, userRepo)

	router := api.NewRouter(api.Handlers{
		Auth:    handlers.NewAuthHandler(authCtrl),
		Menu:    handlers.NewMenuHandler(menuCtrl),
		Cart:    handlers.NewCartHandler(cartCtrl),
		Order:   handlers.NewOrderHandler(orderCtrl),
		Courier: handlers.NewCourierHandler(courierCtrl),
		AuthMW:  middleware.NewAuthMiddleware(tokenSvc),
	})

	log.Printf("server listening on %s, swagger: http://localhost%s/swagger/index.html", addr, addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
