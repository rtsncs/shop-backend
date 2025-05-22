package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ebiznes/handler"
	"ebiznes/model"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		e.Logger.Fatal("failed to connect database", err)
	}

	db.AutoMigrate(&model.Product{})
	db.AutoMigrate(&model.Cart{})
	db.AutoMigrate(&model.CartItem{})
	db.AutoMigrate(&model.User{})

	h := &handler.Handler{DB: db}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	p := e.Group("/products")
	p.POST("", h.CreateProduct)
	p.GET("", h.GetAllProducts)
	p.GET("/:id", h.GetProductByID)
	p.PUT("/:id", h.UpdateProduct)
	p.DELETE("/:id", h.DeleteProduct)

	c := e.Group("/carts")
	c.POST("", h.CreateCart)
	c.GET("", h.GetAllCarts)
	c.GET("/:id", h.GetCartByID)
	c.POST("/:cart_id/items", h.AddItemToCart)
	c.DELETE("/:cart_id/items/:cart_item_id", h.RemoveItemFromCart)
	c.PUT("/:cart_id/items/:cart_item_id", h.UpdateCartItem)

	u := e.Group("/users")
	u.POST("", h.CreateUser)
	u.POST("/login", h.LoginUser)
	u.GET("/current", h.GetCurrentUser)

	e.Logger.Fatal(e.Start(":1323"))

}
