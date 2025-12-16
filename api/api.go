package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func setupRouter(config *Config) *gin.Engine {
	r := gin.Default()
	r.Use(CorsMiddleware(config))

	/*Gets all pets and their status */
	r.GET("/pets", func(c *gin.Context) {
		petRows, err := getPets()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "something broke: " + err.Error()})
		}

		var pets = make(map[string]any)
		for _, row := range petRows {
			if row.Pet.Inside == true {
				pets[row.Pet.Name] = "inside"
			} else {
				pets[row.Pet.Name] = "outside"
			}
		}
		c.JSON(http.StatusOK, pets)

	})

	/*Post a new pet*/
	r.POST("/pets", func(c *gin.Context) {
		var json struct {
			Value string `json:"new_pet_name" binding:"required"`
		}

		err := c.Bind(&json)
		if err == nil {
			newPet := Pet{Name: json.Value, Inside: true}
			rowsAffected, err := addPet(&newPet)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "something broke: " + err.Error()})
			}
			if rowsAffected == 1 {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			} else {
				c.JSON(http.StatusConflict, gin.H{"status": "A pet with this name already exists."})
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON Binding Failed: " + err.Error()})
		}
	})

	/*Update status for an existing pet*/
	r.PATCH("/pets/:petId", func(c *gin.Context) {
		petName := c.Params.ByName("petId")

		var json struct {
			Value string `json:"updated_status" binding:"required"`
		}

		err := c.Bind(&json)
		if err == nil {
			isNowInside := json.Value == "inside"
			updatedPet := Pet{Name: petName, Inside: isNowInside}
			rowsAffected, err := updatePetStatus(&updatedPet)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "something broke: " + err.Error()})
			} else if rowsAffected == 1 {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"status": "Pet ID Not Found"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON Binding Failed: " + err.Error()})
		}

	})

	r.DELETE("/pets/:petId", func(c *gin.Context) {
		pet := c.Params.ByName("petId")
		toDelete := Pet{Name: pet}
		rowsAffected, err := deletePet(&toDelete)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "something broke " + err.Error()})
		}
		if rowsAffected >= 1 {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": "Pet ID Not Found"})
		}

	})

	return r

}

func api(config *Config) {
	r := setupRouter(config)
	r.Run(fmt.Sprintf(":%d", config.ApiPort))
}
