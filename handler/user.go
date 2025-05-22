package handler

import (
	"ebiznes/model"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte("123")

type userCredentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type jwtClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func (h *Handler) CreateUser(c echo.Context) error {
	credentials := userCredentials{}
	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if credentials.Name == "" || credentials.Password == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 14)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := model.User{Name: credentials.Name, Password: string(bytes)}
	if err := h.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, nil)
}

func (h *Handler) LoginUser(c echo.Context) error {
	credentials := userCredentials{}
	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if credentials.Name == "" || credentials.Password == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}

	user := model.User{}
	if err := h.DB.First(&user, model.User{Name: credentials.Name}).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, nil)
	}

	claims := jwtClaims{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	cookie := http.Cookie{Name: "token", Value: tokenString, Secure: true, HttpOnly: true}
	c.SetCookie(&cookie)
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetCurrentUser(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &jwtClaims{}, func(t *jwt.Token) (any, error) { return secret, nil })
	if err != nil {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	claims, ok := token.Claims.(*jwtClaims)
	if claims == nil || !ok {
		return c.JSON(http.StatusInternalServerError, "unknown jwt claims type")
	}

	user := model.User{}
	if err := h.DB.First(&user, claims.UserID).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, nil)
	}

	return c.JSON(http.StatusOK, user)
}
