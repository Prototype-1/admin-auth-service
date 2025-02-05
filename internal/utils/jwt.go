package utils

import (
    "errors"
    "github.com/golang-jwt/jwt/v4" 
    "time"
    "os"
    "log"
    "fmt"
    "github.com/joho/godotenv"
)

var jwtSecretKey []byte

func init() {
    err := godotenv.Load("config/.env")
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
    key := os.Getenv("JWT_SECRET_KEY")
    if key == "" {
        log.Fatal("JWT_SECRET_KEY is not set")
    }
    jwtSecretKey = []byte(key)
    fmt.Printf("Loaded JWT Secret Key: %s\n", key)
}

func GenerateJWT(adminID int, secretKey string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "admin_id": adminID,
        "exp":      time.Now().Add(time.Hour * 1).Unix(),
    })
    return token.SignedString([]byte(secretKey))
}

func ParseJWT(tokenStr string) (uint, error) {
    log.Printf("Using JWT Secret Key: %s\n", string(jwtSecretKey))
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return jwtSecretKey, nil
    })
    if err != nil {
        log.Println("Error parsing JWT:", err)
        return 0, err
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return 0, errors.New("invalid token")
    }
    
    adminIDFloat, ok := claims["admin_id"].(float64)
    if !ok {
        return 0, errors.New("invalid admin_id in token")
    }
    
    return uint(adminIDFloat), nil
}