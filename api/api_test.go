package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// Made with consulting the following sources:
// https://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go
// https://github.com/gin-gonic/examples/blob/master/basic/main_test.go
// TODO: Change package name from main, see: https://appliedgo.net/testmain/
// TODO: Read on this: https://mortenvistisen.com/posts/integration-tests-with-docker-and-go
// Also: https://github.com/MBvisti/integration-test-in-go/tree/vanilla-approach/running-integration-tests-using-std-library

type ApiTestSuite struct {
	suite.Suite
	Config *Config
	Router *gin.Engine
}

func (suite *ApiTestSuite) SetupTest() {
	suite.Config = New()
	suite.Router = setupRouter(suite.Config)
}

func postCatOne(suite *ApiTestSuite) {
	w := httptest.NewRecorder()
	var jsonString = []byte(`{"new_pet_name": "cat1"}`)
	req, _ := http.NewRequest("POST", "/pets", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	suite.Router.ServeHTTP(w, req)
}

func (suite *ApiTestSuite) TestPostRouteSuccess() {

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"new_pet_name": "newPostedCat"}`)
	req, _ := http.NewRequest("POST", "/pets", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	suite.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(`{"status":"ok"}`, w.Body.String())

}

func (suite *ApiTestSuite) TestGetRoute() {
	postCatOne(suite) //set up for test

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pets", nil)
	suite.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "cat1")
	suite.Contains(w.Body.String(), "inside")
}

func (suite *ApiTestSuite) TestPostRouteAlreadyExists() {
	postCatOne(suite)

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"new_pet_name": "cat1"}`)
	req, _ := http.NewRequest("POST", "/pets", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	suite.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusConflict, w.Code)
	suite.Equal(`{"status":"A pet with this name already exists."}`, w.Body.String())
}

func (suite *ApiTestSuite) TestPatchRouteSuccess() {
	postCatOne(suite)

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"updated_status": "outside"}`)
	req, _ := http.NewRequest("PATCH", "/pets/cat1", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	suite.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(`{"status":"ok"}`, w.Body.String())
}

func (suite *ApiTestSuite) TestPatchRouteNotFound() {

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"updated_status": "outside"}`)
	req, _ := http.NewRequest("PATCH", "/pets/cat2", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	suite.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Equal(`{"status":"Pet ID Not Found"}`, w.Body.String())
}

func (suite *ApiTestSuite) TestDeleteRouteSuccess() {
	postCatOne(suite)

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"updated_status": "outside"}`)
	req, _ := http.NewRequest("DELETE", "/pets/cat1", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	suite.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(`{"status":"ok"}`, w.Body.String())
}

func (suite *ApiTestSuite) TestDeleteRouteNotFound() {

	w := httptest.NewRecorder()
	var jsonString = []byte(`{"updated_status": "outside"}`)
	req, _ := http.NewRequest("DELETE", "/pets/cat1", bytes.NewBuffer(jsonString))
	req.Header.Add("content-type", "application/json")
	suite.Router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Equal(`{"status":"Pet ID Not Found"}`, w.Body.String())
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
