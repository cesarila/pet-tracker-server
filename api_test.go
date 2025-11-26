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
