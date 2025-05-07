package handler

import (
    "strconv"
	"net/http"
	"github.com/labstack/echo/v4"
    "ebiznes/model"
)

func (h *Handler) CreateCart(c echo.Context) error {
    cart := model.Cart{}
    if err := h.DB.Create(&cart).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusCreated, cart)
}

func (h *Handler) GetAllCarts(c echo.Context) error {
    carts := []model.Cart{}
    if err := h.DB.Preload("CartItems").Preload("CartItems.Product").Find(&carts).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusOK, carts)
}

func (h *Handler) GetCartByID(c echo.Context) error {
    id := c.Param("id")
    cart := model.Cart{}
    if err := h.DB.Preload("CartItems").Preload("CartItems.Product").First(&cart, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, err.Error())
    }
    return c.JSON(http.StatusOK, cart)
}

func (h *Handler) AddItemToCart(c echo.Context) error {
    cartID, err :=  strconv.ParseUint(c.Param("cart_id"), 10, 32)
    if  err != nil {
        return c.JSON(http.StatusBadRequest, "Invalid cart ID")
    }
    cartItem := model.CartItem{}
    if err := c.Bind(&cartItem); err != nil {
        return c.JSON(http.StatusBadRequest, err.Error())
    }
    cartItem.CartID = uint(cartID)
    if err := h.DB.Create(&cartItem).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusCreated, cartItem)
}

func (h *Handler) RemoveItemFromCart(c echo.Context) error {
    cartItemID := c.Param("cart_item_id")
    if err := h.DB.Delete(&model.CartItem{}, cartItemID).Error; err != nil {
        return c.JSON(http.StatusNotFound, err.Error())
    }
    return c.NoContent(http.StatusNoContent)
}

func (h *Handler) UpdateCartItem(c echo.Context) error {
    cartItemID := c.Param("cart_item_id")
    cartItem := model.CartItem{}
    if err := h.DB.First(&cartItem, cartItemID).Error; err != nil {
        return c.JSON(http.StatusNotFound, err.Error())
    }
    if err := c.Bind(&cartItem); err != nil {
        return c.JSON(http.StatusBadRequest, err.Error())
    }
    if err := h.DB.Save(&cartItem).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusOK, cartItem)
}
