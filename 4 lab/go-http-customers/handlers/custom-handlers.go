package handlers

import (
	"fmt"
	"go-http-customers/database"
	"go-http-customers/models"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Секретный ключ для создания JWT
var secretKey = []byte("your_secret_key")

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Функция для проверки пароля
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Функция для создания JWT токена
func createJWT(userID uint) (string, error) {
	// Устанавливаем срок действия токена (например, 1 час)
	expirationTime := time.Now().Add(1 * time.Hour).Unix()

	// Создаем полезную нагрузку
	payload := jwt.MapClaims{
		"sub": fmt.Sprintf("%d", userID),
		"exp": expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// Обработчик для авторизации пользователя
func LoginCustomer(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		UserPass string `json:"user_pass"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	var customer models.Customer
	if result := database.DB.Where("email = ?", loginData.Email).First(&customer); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if !checkPasswordHash(loginData.UserPass, customer.UserPass) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := createJWT(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CreateCustomer(c *gin.Context) {
	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	hashedPassword, err := hashPassword(customer.UserPass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	customer.UserPass = hashedPassword

	if result := database.DB.Create(&customer); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create customer"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// Обработчик для получения данных покупателя
func GetCustomerData(c *gin.Context) {
	idParam := c.Param("id")
	requestedID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Извлекаем ID пользователя из токена (добавлено в middleware)
	userIDFromToken, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Authorization failed"})
		return
	}

	// Проверяем, что запрашиваемый ID соответствует ID из токена
	if strconv.Itoa(requestedID) != userIDFromToken {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}

	var customer models.Customer
	if result := database.DB.First(&customer, requestedID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    customer.ID,
		"email": customer.Email,
	})
}
