package BooksController

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	Model "project/models"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func sqliteDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func uploadFile(w http.ResponseWriter, r *http.Request) string {
	file, _, errForm := r.FormFile("cover_image")

	if errForm != nil {
		panic("Error FormFile: " + errForm.Error())
	}

	defer file.Close()

	tempFile, errTemp := ioutil.TempFile("assets/images", "cover-*.png")

	if errTemp != nil {
		panic("Error TempDir: " + errTemp.Error())
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		panic("Error ReadAll: " + err.Error())
	}

	tempFile.Write(fileBytes)
	var fileName = strings.Replace(tempFile.Name(), "assets", "static", -1)
	fileName = strings.Replace(fileName, "\\", "/", -1)
	return fileName
}

func renderTemplateHTML(htmlTmp string, w http.ResponseWriter, data interface{}) {
	files := []string{
		"views/" + htmlTmp + ".html",
		"views/base.html",
	}
	tmpt, err := template.ParseFiles(files...)

	if err != nil {
		panic("Error Template: " + err.Error())
	}

	errExec := tmpt.ExecuteTemplate(w, "base", data)

	if errExec != nil {
		panic("Error Execute: " + errExec.Error())
	}

}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := sqliteDB()
	var books []Model.Books
	db.Find(&books)
	datas := map[string]interface{}{
		"Books": books,
	}

	renderTemplateHTML("index", w, datas)
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := sqliteDB()

	if r.Method == "POST" {
		filename := uploadFile(w, r)
		fmt.Println(filename)

		book := Model.Books{
			Name:        r.FormValue("name"),
			Author:      r.FormValue("author"),
			Description: r.FormValue("description"),
			ImageCover:  filename,
		}

		db.Create(&book)

		http.Redirect(w, r, "/", http.StatusFound)

	} else {
		renderTemplateHTML("create", w, nil)
	}
}

func Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	db := sqliteDB()
	book := Model.Books{}
	err := db.First(&book, params.ByName("id")).Error

	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		book.Name = r.FormValue("name")
		book.Description = r.FormValue("description")
		book.Author = r.FormValue("author")
		db.Save(&book)

		http.Redirect(w, r, "/", http.StatusFound)
	} else {

		datas := map[string]interface{}{
			"Book": book,
		}
		renderTemplateHTML("update", w, datas)
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	db := sqliteDB()
	book := Model.Books{}
	err := db.First(&book, params.ByName("id")).Error

	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	db.Delete(&book, params.ByName("id"))

	http.Redirect(w, r, "/", http.StatusFound)
}
