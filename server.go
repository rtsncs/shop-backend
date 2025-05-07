package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

    "gorm.io/gorm"
    "gorm.io/driver/sqlite"

    "ebiznes/model"
    "ebiznes/handler"
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

    h := &handler.Handler{DB: db}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
    e.POST("/products", h.CreateProduct)
    e.GET("/products", h.GetAllProducts)
    e.GET("/products/:id", h.GetProductByID)
    e.PUT("/products/:id", h.UpdateProduct)
    e.DELETE("/products/:id", h.DeleteProduct)
    e.POST("/carts", h.CreateCart)
    e.GET("/carts", h.GetAllCarts)
    e.GET("/carts/:id", h.GetCartByID)
    e.POST("/carts/:cart_id/items", h.AddItemToCart)
    e.DELETE("/carts/:cart_id/items/:cart_item_id", h.RemoveItemFromCart)
    e.PUT("/carts/:cart_id/items/:cart_item_id", h.UpdateCartItem)
	e.Logger.Fatal(e.Start(":1323"))
}
