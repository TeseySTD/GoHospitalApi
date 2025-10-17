package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	SecretKey = "my-super-secret-key-change-in-production"
	
	RoleAdmin  = "admin"
	RoleReader = "reader"
)

// Claims структура для JWT токена
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// User структура користувача
type User struct {
	Username string
	Password string
	Role     string
}

// Хардкодні користувачі (у реальному проекті - база даних)
var Users = map[string]User{
	"admin": {
		Username: "admin",
		Password: "admin123",
		Role:     RoleAdmin,
	},
	"reader": {
		Username: "reader",
		Password: "reader123",
		Role:     RoleReader,
	},
	"doctor": {
		Username: "doctor",
		Password: "doctor123",
		Role:     RoleAdmin,
	},
	"viewer": {
		Username: "viewer",
		Password: "viewer123",
		Role:     RoleReader,
	},
}

// GenerateToken генерує JWT токен для користувача
func GenerateToken(username, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	
	claims := &Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ValidateToken перевіряє та парсить JWT токен
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	return claims, nil
}

// AuthenticateUser перевіряє логін і пароль
func AuthenticateUser(username, password string) (*User, error) {
	user, exists := Users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	if user.Password != password {
		return nil, errors.New("invalid password")
	}
	
	return &user, nil
}

// IsAdmin перевіряє чи користувач має роль адміна
func IsAdmin(role string) bool {
	return role == RoleAdmin
}

// IsReader перевіряє чи користувач має роль читача
func IsReader(role string) bool {
	return role == RoleReader
}