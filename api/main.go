package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := New()
	initDatabase(fmt.Sprintf("%s/%s", conf.Database.SqliteFolderPath, conf.Database.SqliteFileName))
	api(conf)
}
