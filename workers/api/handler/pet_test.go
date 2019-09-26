package handler

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHandler_addPetToStore(t *testing.T) {
	raw := `
{
  "id": 3,
  "category": {
    "id": 2,
    "name": "Cat"
  },
  "name": "TestPet",
  "photo_urls": [
    "string-try",
    "string-try"
  ],
  "tags": [
    {
      "id": 1,
      "name": "small"
    }
  ],
  "status": "available"
}`

	request, err := http.NewRequest("POST", "/api/v2/pet", strings.NewReader(raw))
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/v2/pet", addPetToStore)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_deletePetByID(t *testing.T) {
	request, err := http.NewRequest("DELETE", "/api/v2/pet/1", nil)
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/api/v2/pet/{petID}", deletePetByID)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_findPetsByStatus(t *testing.T) {
	reqArgs := url.Values{}
	reqArgs.Add("status", "available")
	reqURL, _ := url.Parse(host)
	reqURL.Path = "/api/v2/pet/findByStatus"
	reqURL.RawQuery = reqArgs.Encode()

	request, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/api/v2/pet/findByStatus", findPetsByStatus)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_getPetByID(t *testing.T) {
	request, err := http.NewRequest("GET", "/api/v2/pet/1", nil)
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/api/v2/pet/{petID}", getPetByID)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestHandler_updatePetByID(t *testing.T) {
	reqArgs := url.Values{}
	reqArgs.Add("name", "TestPet")
	reqArgs.Add("status", "pending")

	request, err := http.NewRequest("POST", "/api/v2/pet/1", nil)
	if err != nil {
		log.Println(err)
	}

	request.PostForm = reqArgs

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/v2/pet/{petID}", updatePetByID)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_updatePetInStore(t *testing.T) {
	raw := `
{
  "id": 1,
  "category": {
    "id": 2,
    "name": "Cat"
  },
  "name": "TestPet",
  "photo_urls": [
    "string-try",
    "string-try"
  ],
  "tags": [
    {
      "id": 1,
      "name": "small"
    }
  ],
  "status": "available"
}`

	request, err := http.NewRequest("POST", "/api/v2/pet", strings.NewReader(raw))
	if err != nil {
		log.Println(err)
	}

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/v2/pet", updatePetInStore)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestHandler_updatePetPhotosByID(t *testing.T) {
	reqArgs := url.Values{}
	reqArgs.Add("images",
		"https://i.ytimg.com/vi/YCaGYUIfdy4/maxresdefault.jpg,"+
			"https://d17fnq9dkz9hgj.cloudfront.net/uploads/2012/11/152964589-welcome-home-new-cat-632x475.jpg")

	request, err := http.NewRequest("POST", "/api/v2/pet/1/uploadImage", nil)
	if err != nil {
		log.Println(err)
	}

	request.PostForm = reqArgs

	request = addToCtxWriteTimeout(request)
	response := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/v2/pet/{petID}/uploadImage", updatePetPhotosByID)

	r.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}
