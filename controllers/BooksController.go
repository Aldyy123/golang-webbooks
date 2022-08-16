package BooksController

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	Model "project/models"

	validator "github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type books struct {
	Name        string `validate:"required"`
	Author      string `validate:"required"`
	Description string `validate:"required,min=10"`
	ImageCover  string `validate:"required"`
}

func sqliteDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func tempFile() *os.File {
	tempFile, errTemp := ioutil.TempFile("assets/images", "cover-*.png")

	if errTemp != nil {
		panic("Error TempDir: " + errTemp.Error())
	}

	return tempFile
}

func saveFileImage(file multipart.File, tempFile *os.File) string {

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		panic("Error ReadAll: " + err.Error())
	}

	tempFile.Write(fileBytes)
	var fileName = strings.Replace(tempFile.Name(), "assets", "static", -1)
	fileName = strings.Replace(fileName, "\\", "/", -1)
	return fileName
}

func uploadFile(w http.ResponseWriter, r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	file, header, errForm := r.FormFile("cover_image")
	return file, header, errForm
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

	data := map[string]interface{}{
		"Error": nil,
	}

	if r.Method == "POST" {

		file, headerFile, errForm := uploadFile(w, r)

		book := books{
			Name:        r.FormValue("name"),
			Author:      r.FormValue("author"),
			Description: r.FormValue("description"),
			ImageCover:  headerFile.Filename,
		}

		validate := validator.New()
		err := validate.Struct(book)
		fmt.Println(err)
		var multiErrors = []interface{}{
			err,
		}

		if errForm != nil {
			multiErrors = append(multiErrors, errForm.Error())
		}

		data["Error"] = multiErrors

		if err != nil {
			renderTemplateHTML("create", w, data)
			return
		}
		defer file.Close()

		tempFile := tempFile()
		defer tempFile.Close()

		filename := saveFileImage(file, tempFile)
		book.ImageCover = filename
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

	data := map[string]interface{}{
		"Error": nil,
	}
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		file, _, _ := uploadFile(w, r)

		book.Name = r.FormValue("name")
		book.Description = r.FormValue("description")
		book.Author = r.FormValue("author")

		validate := validator.New()
		err := validate.Struct(book)
		fmt.Println(file)
		data["Error"] = err

		if err != nil {
			renderTemplateHTML("update", w, data)
			return
		}

		if file != nil {
			defer file.Close()

			tempFile := tempFile()
			defer tempFile.Close()

			filename := saveFileImage(file, tempFile)
			book.ImageCover = filename
		}

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
