package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	paypalClientID     = os.Getenv("PAYPAL_CLIENT_ID")
	paypalClientSecret = os.Getenv("PAYPAL_CLIENT_SECRET")
	port               = os.Getenv("PORT")
	baseURL            = "https://api-m.sandbox.paypal.com"
)

type orderPayload struct {
	Intent        string `json:"intent"`
	PurchaseUnits []struct {
		Amount struct {
			CurrencyCode string `json:"currency_code"`
			Value        string `json:"value"`
		} `json:"amount"`
	} `json:"purchase_units"`
}

func GenerateAccessToken() (string, error) {
	if paypalClientID == "" || paypalClientSecret == "" {
		return "", fmt.Errorf("missing API credentials")
	}

	auth := base64.StdEncoding.EncodeToString([]byte(paypalClientID + ":" + paypalClientSecret))
	url := fmt.Sprintf("%s/v1/oauth2/token", baseURL)
	payload := []byte("grant_type=client_credentials")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	accessToken, ok := data["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get access token")
	}

	return accessToken, nil
}

func createOrder(c *gin.Context) {
	var cart map[string]interface{}
	if err := c.BindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	accessToken, err := GenerateAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	url := fmt.Sprintf("%s/v2/checkout/orders", baseURL)
	payload := orderPayload{
		Intent: "CAPTURE",
		PurchaseUnits: []struct {
			Amount struct {
				CurrencyCode string `json:"currency_code"`
				Value        string `json:"value"`
			} `json:"amount"`
		}{
			{
				Amount: struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				}{
					CurrencyCode: "INR",
					Value:        "100.00",
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal payload"})
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	defer resp.Body.Close()

	result, err := handleResponse(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(resp.StatusCode, result)
}

func captureOrder(c *gin.Context) {
	orderID := c.Param("orderID")

	accessToken, err := GenerateAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	url := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", baseURL, orderID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to capture order"})
		return
	}
	defer resp.Body.Close()

	result, err := handleResponse(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(resp.StatusCode, result)
}

func handleResponse(resp *http.Response) (gin.H, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return gin.H{
		"jsonResponse":     data,
		"httpStatusCode":   resp.StatusCode,
		"originalResponse": body,
	}, nil
}

func Mmain() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	r.POST("/api/orders", createOrder)
	r.POST("/api/orders/:orderID/capture", captureOrder)

	r.Run(":" + port)
}
