package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// First login to get token
	loginURL := "http://localhost:8080/api/v1/auth/login"
	loginPayload := map[string]string{
		"username": "Administrator",
		"password": "Password.1",
	}

	jsonData, _ := json.Marshal(loginPayload)

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating login request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making login request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var loginResponse map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &loginResponse); err != nil {
		fmt.Printf("Error parsing login response: %v\n", err)
		return
	}

	token, ok := loginResponse["token"].(string)
	if !ok {
		fmt.Printf("No token in response: %s\n", string(body))
		return
	}

	fmt.Printf("Login successful, token: %s...\n", token[:20])

	// Now test the global settings endpoint
	url := "http://localhost:8080/api/v1/admin/global-settings/discount_rate"

	payload := map[string]string{
		"value": "10",
	}

	jsonData, _ = json.Marshal(payload)

	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))
}