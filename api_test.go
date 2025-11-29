package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Made with consulting the following sources:
// https://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go
// https://github.com/gin-gonic/examples/blob/master/basic/main_test.go
// TODO: Change package name from main, see: https://appliedgo.net/testmain/
// TODO: Read on this: https://mortenvistisen.com/posts/integration-tests-with-docker-and-go
// Also: https://github.com/MBvisti/integration-test-in-go/tree/vanilla-approach/running-integration-tests-using-std-library
func TestPostRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"new_pet_name": "cat1"}`)
	req, _ := http.NewRequest("POST", "/pets", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"status":"ok"}`, w.Body.String())

}

func TestPostRouteAlreadyExists(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"new_pet_name": "cat1"}`)
	req, _ := http.NewRequest("POST", "/pets", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, `{"status":"A pet with this name already exists."}`, w.Body.String())
}

func TestPatchRouteSuccess(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"updated_status": "outside"}`)
	req, _ := http.NewRequest("PATCH", "/pets/cat1", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"status":"ok"}`, w.Body.String())
}

func TestPatchRouteNotFound(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"updated_status": "outside"}`)
	req, _ := http.NewRequest("PATCH", "/pets/cat2", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, `{"status":"Pet ID Not Found"}`, w.Body.String())
}

func TestDeleteRouteSuccess(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"updated_status": "outside"}`)
	req, _ := http.NewRequest("DELETE", "/pets/cat1", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"status":"ok"}`, w.Body.String())
}

func TestDeleteRouteNotFound(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"updated_status": "outside"}`)
	req, _ := http.NewRequest("DELETE", "/pets/cat1", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, `{"status":"Pet ID Not Found"}`, w.Body.String())
}
