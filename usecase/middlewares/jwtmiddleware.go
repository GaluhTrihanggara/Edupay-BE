package middlewares

import (
	"Edupay/config"
	"Edupay/model"
	"net/http"

	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Middlewares interface {
	CreateToken(user model.User) (string, error)
	ExtractTokenUserId(userType string, c echo.Context) string
	AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

func CreateToken(user model.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = user.ID
	claims["user_type"] = user.UserType
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.SecretJWT))

}

func ExtractTokenUserId(userType string, c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.SecretJWT), nil
	})
	if err != nil || !token.Valid {
		return ""
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["userId"].(string)
	if userType == model.ALL_TYPE {
		return string(userId)
	} else if claims["user_type"].(string) == userType {
		return string(userId)
	} else {
		return ""
	}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid Authorization header")
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.SecretJWT), nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
		}

		claims := token.Claims.(jwt.MapClaims)
		userId, ok := claims["userId"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID")
		}

		c.Set("userId", userId)

		return next(c)
	}
}
