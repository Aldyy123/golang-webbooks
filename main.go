package main

import (
	"fmt"
	"net/http"
	BooksController "project/controllers"

	Model "project/models"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Model.Books{})

	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("assets"))

	router.GET("/", BooksController.Index)
	router.GET("/create", BooksController.Create)
	router.POST("/create", BooksController.Create)
	router.GET("/update/:id", BooksController.Update)
	router.POST("/update/:id", BooksController.Update)
	router.GET("/delete/:id", BooksController.DeleteBook)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
