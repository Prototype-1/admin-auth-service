package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var jwtSecretKey []byte

func init() {
    err := godotenv.Load("config/.env")
    if err != nil {
        Log.Fatal("Error loading .env file", zap.Error(err))
    }
    key := os.Getenv("JWT_SECRET_KEY")
    if key == "" {
        log.Fatal("JWT_SECRET_KEY is not set")
    }
    jwtSecretKey = []byte(key)
}

func GenerateJWT(adminID int, role string, secretKey string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "admin_id": adminID,
        "role":     role,  
        "exp":      time.Now().Add(time.Hour * 1).Unix(),
    })
    return token.SignedString([]byte(secretKey))
}

func ParseJWT(tokenStr string) (uint, string, error) {
    log.Printf("Using JWT Secret Key: %s\n", string(jwtSecretKey))
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return jwtSecretKey, nil
    })
    if err != nil {
        log.Println("Error parsing JWT:", err)
        return 0, "", err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return 0, "", errors.New("invalid token")
    }

    adminIDFloat, ok := claims["admin_id"].(float64)
    if !ok {
        return 0, "", errors.New("invalid admin_id in token")
    }

    role, ok := claims["role"].(string)
    if !ok {
        return 0, "", errors.New("invalid role in token")
    }

    return uint(adminIDFloat), role, nil
}
