package test

import (
	"bytes"
	"jwtreact/controllers"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestLogin(t *testing.T) {
	app := fiber.New()
	app.Post("/login", controllers.Login)

	// Test case 1: Invalid email
	req1, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer([]byte(`{"email": "invalid_email@example.com", "password": "password"}`)))
	req1.Header.Set("Content-Type", "application/json")
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", fiber.StatusNotFound, resp1.StatusCode)
	}

	// Test case 2: Invalid password
	req2, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer([]byte(`{"email": "testing@example.com", "password": "bayu"}`)))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", fiber.StatusBadRequest, resp2.StatusCode)
	}

	// Test case 3: Valid credentials
	req3, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer([]byte(`{"email": "Ubay@example.com", "password": "testing"}`)))
	req3.Header.Set("Content-Type", "application/json")
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status code %d, but got %d", fiber.StatusOK, resp3.StatusCode)
	}
}
