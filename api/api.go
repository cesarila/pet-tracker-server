package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: Figure out if using a type like this is helpful or not
type Pet struct {
	Name   string `json:"name" binding:"required"`
	Inside bool   `json:"inside" binding:"required"`
}

//Some code written with reference to https://github.com/gin-gonic/examples/blob/master/basic/main.go

// CORS stuff helped by this post: https://jgunnink.substack.com/p/gin-framework-in-go-implementing-cors-effectively

func CorsMiddleware(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("%s:%d", config.Frontend.Host, config.Frontend.Port))
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Accept-Patch", "application/json")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Placeholder for actual database
var db = make(map[string]string)

func setupRouter(config *Config) *gin.Engine {
	r := gin.Default()
	r.Use(CorsMiddleware(config))

	/*Gets all pets and their status */
	r.GET("/pets", func(c *gin.Context) {

		// I think this is ok because db placeholder is a map of string to string, and gin.H is a shortcut for a map of string to any.
		c.JSON(http.StatusOK, db)
	})

	/*Post a new pet*/
	r.POST("/pets", func(c *gin.Context) {
		var json struct {
			Value string `json:"new_pet_name" binding:"required"`
		}

		err := c.Bind(&json)
		if err == nil {
			_, exists := db[json.Value]
			if exists {
				c.JSON(http.StatusConflict, gin.H{"status": "A pet with this name already exists."})
			} else {
				//new pet assumed to be inside
				db[json.Value] = "inside"
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON Binding Failed: " + err.Error()})
		}
	})

	/*Update status for an existing pet*/
	r.PATCH("/pets/:petId", func(c *gin.Context) {
		pet := c.Params.ByName("petId")
		_, ok := db[pet]
		if ok {
			var json struct {
				Value string `json:"updated_status" binding:"required"`
			}

			err := c.Bind(&json)
			if err == nil {
				db[pet] = json.Value
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "JSON Binding Failed: " + err.Error()})
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"status": "Pet ID Not Found"})
		}
	})

	r.DELETE("/pets/:petId", func(c *gin.Context) {
		pet := c.Params.ByName("petId")
		_, ok := db[pet]
		if ok {
			delete(db, pet)
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"status": "Pet ID Not Found"})
		}

	})

	return r

}

func api(config *Config) {
	r := setupRouter(config)
	r.Run(fmt.Sprintf(":%d", config.ApiPort))
}
