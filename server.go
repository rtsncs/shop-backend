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
	e.Logger.Fatal(e.Start(":1323"))
}
