package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"lab7-oop/api/handlers"
	"lab7-oop/api/middleware"
	"lab7-oop/models"
	_ "lab7-oop/docs"
)

type Handlers struct {
	Auth    *handlers.AuthHandler
	Menu    *handlers.MenuHandler
	Cart    *handlers.CartHandler
	Order   *handlers.OrderHandler
	Courier *handlers.CourierHandler
	AuthMW  *middleware.AuthMiddleware
}

func NewRouter(h Handlers) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.Health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/auth/register", h.Auth.Register)
	r.POST("/auth/login", h.Auth.Login)

	api := r.Group("/")
	api.Use(h.AuthMW.RequireAuth())
	{
		api.GET("/auth/me", h.Auth.Me)

		api.GET("/restaurants", h.Menu.ListRestaurants)
		api.GET("/restaurants/:id", h.Menu.GetRestaurant)
		api.GET("/restaurants/:id/dishes", h.Menu.ListDishes)

		customer := api.Group("/")
		customer.Use(h.AuthMW.RequireRoles(models.RoleCustomer, models.RoleAdmin))
		{
			customer.GET("/cart", h.Cart.GetCart)
			customer.POST("/cart/items", h.Cart.AddToCart)
			customer.DELETE("/cart/items/:dishId", h.Cart.RemoveFromCart)

			customer.POST("/orders", h.Order.CreateOrder)
			customer.GET("/orders", h.Order.ListOrders)
			customer.GET("/orders/:id", h.Order.GetOrder)
		}

		staff := api.Group("/")
		staff.Use(h.AuthMW.RequireRoles(models.RoleAdmin))
		{
			staff.PUT("/orders/:id/status", h.Order.SetOrderStatus)
			staff.POST("/orders/:id/courier", h.Order.AssignCourier)
			staff.GET("/couriers", h.Courier.ListCouriers)
		}

		ops := api.Group("/")
		ops.Use(h.AuthMW.RequireRoles(models.RoleAdmin, models.RoleCourier))
		{
			ops.POST("/orders/:id/advance", h.Order.AdvanceOrderStatus)
		}
	}

	return r
}
