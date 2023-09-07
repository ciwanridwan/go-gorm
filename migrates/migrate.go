package main

import (
	"fmt"
	"log"

	"go-gorm/models"

	"go-gorm/initializers"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal(" Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Article{})
	fmt.Println("Migration complete")
}
