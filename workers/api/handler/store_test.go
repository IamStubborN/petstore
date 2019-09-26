package handler

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHandler_createOrder(t *testing.T) {
	raw := `
{
  "id": 1,
  "pet_id": 2,
  "user_id": 2,
  "quantity": 12,
  "ship_date": "2019-08-17T20:30:39.563Z",
  "status": "placed",
  "complete": false
}`

	request, err := http.NewRequest("POST", "/api/v2/store/order", strings.NewReader(raw))
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/v2/store/order", createOrder)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_deletePurchaseByID(t *testing.T) {
	request, err := http.NewRequest("DELETE", "/api/v2/store/order/1", nil)
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/api/v2/store/order/{orderID}", deletePurchaseByID)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_findOrderByID(t *testing.T) {
	request, err := http.NewRequest("GET", "/api/v2/store/order/1", nil)
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/api/v2/store/order/{orderID}", findOrderByID)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_inventoriesByStatus(t *testing.T) {
	request, err := http.NewRequest("GET", "/api/v2/store/order/inventory", nil)
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/api/v2/store/order/inventory", inventoriesByStatus)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}
