package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Задание 1
	router.GET("/greet", func(c *gin.Context) {
		name := c.Query("name")
		age := c.Query("age")
		response := "Меня зовут " + name + ", мне " + age + " лет"
		c.String(http.StatusOK, response)
	})

	// Задание 2
	router.GET("/add", func(c *gin.Context) {
		a := c.Query("a")
		b := c.Query("b")
		result, err := calculate(a, b, "add")
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusOK, "Результат: %s", result)
	})

	router.GET("/sub", func(c *gin.Context) {
		a := c.Query("a")
		b := c.Query("b")
		result, err := calculate(a, b, "sub")
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusOK, "Результат: %s", result)
	})

	router.GET("/mul", func(c *gin.Context) {
		a := c.Query("a")
		b := c.Query("b")
		result, err := calculate(a, b, "mul")
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusOK, "Результат: %s", result)
	})

	router.GET("/div", func(c *gin.Context) {
		a := c.Query("a")
		b := c.Query("b")
		result, err := calculate(a, b, "div")
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusOK, "Результат: %s", result)
	})

	router.POST("/count", func(c *gin.Context) {
		var jsonData struct {
			Text string `json:"text"`
		}
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		counts := make(map[rune]int)
		for _, char := range jsonData.Text {
			counts[char]++
		}
		c.JSON(http.StatusOK, counts)
	})

	router.Run(":8080")
}

func calculate(aStr, bStr, operation string) (string, error) {
	a, err := strconv.Atoi(aStr)
	if err != nil {
		return "", err
	}
	b, err := strconv.Atoi(bStr)
	if err != nil {
		return "", err
	}

	var result int
	switch operation {
	case "add":
		result = a + b
	case "sub":
		result = a - b
	case "mul":
		result = a * b
	case "div":
		if b == 0 {
			return "", errors.New("деление на ноль")
		}
		result = a / b
	default:
		return "", errors.New("неизвестная операция")
	}
	return strconv.Itoa(result), nil
}
