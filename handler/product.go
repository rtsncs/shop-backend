package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
    "ebiznes/model"
)

func (h *Handler) CreateProduct(c echo.Context) error {
    product := model.Product{}
    if err := c.Bind(&product); err != nil {
        return c.JSON(http.StatusBadRequest, err.Error())
    }
    if err := h.DB.Create(&product).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusCreated, product)
}

func (h *Handler) GetProductByID(c echo.Context) error {
    id := c.Param("id")
    product := model.Product{}
    if err := h.DB.First(&product, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, err.Error())
    }
    return c.JSON(http.StatusOK, product)
}

func (h *Handler) GetAllProducts(c echo.Context) error {
    products := []model.Product{}
    if err := h.DB.Find(&products).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusOK, products)
}

func (h *Handler) UpdateProduct(c echo.Context) error {
    id := c.Param("id")
    product := model.Product{}
    if err := h.DB.First(&product, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, err.Error())
    }
    if err := c.Bind(&product); err != nil {
        return c.JSON(http.StatusBadRequest, err.Error())
    }
    if err := h.DB.Save(&product).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusOK, product)
}

func (h *Handler) DeleteProduct(c echo.Context) error {
    id := c.Param("id")
    if err := h.DB.Delete(&model.Product{}, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, err.Error())
    }
    return c.NoContent(http.StatusNoContent)
}
