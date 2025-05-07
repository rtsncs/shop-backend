package model

import (
    "gorm.io/gorm"
)

type CartItem struct {
    gorm.Model
    CartID    uint    `json:"cart_id"`
    ProductID uint    `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Product   Product `gorm:"foreignKey:ProductID"`
}
