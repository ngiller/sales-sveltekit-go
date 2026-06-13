package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBStockName string
	JWTSecret   string
}

var AppConfig *Config

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("FATAL: Environment variable %s is not set", key)
	}
	return val
}

func InitConfig() {
	AppConfig = &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  requireEnv("DB_PASSWORD"),
		DBName:      getEnv("DB_NAME", "magnum_sales_svelte_go"),
		DBStockName: getEnv("DB_STOCK_NAME", "magnum_stock_db"),
		JWTSecret:   requireEnv("JWT_SECRET"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var DB *gorm.DB
var StockDB *gorm.DB

func InitDB() *gorm.DB {
	InitConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?collation=utf8mb4_general_ci&parseTime=True&loc=Local",
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to sales database: %v", err)
	}

	setupDBPool(db)
	DB = db
	return db
}

func InitStockDB() *gorm.DB {
	if AppConfig == nil {
		InitConfig()
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?collation=utf8mb4_general_ci&parseTime=True&loc=Local",
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBStockName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to stock database: %v", err)
	}

	setupDBPool(db)
	StockDB = db
	return db
}

func setupDBPool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal server error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		msg = e.Message
	}
	log.Printf("HTTP Error %d: %v", code, err)
	return c.Status(code).JSON(fiber.Map{
		"message": msg,
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID int64, email string, inisial string, userGroupID *uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"inisial": inisial,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}
	if userGroupID != nil {
		claims["user_group_id"] = *userGroupID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(AppConfig.JWTSecret))
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
