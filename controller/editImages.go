package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"isolati.cn/db"
	"isolati.cn/global"
)

type ImageData struct {
	Suffix string `json:"suffix"`
	Data   string `json:"data"`
}

var editImagesTemplate = template.New("editParagraphs")

func showEditImagesPage(w http.ResponseWriter, r *http.Request) {
	editImagesTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		CssFiles: []string{"/css/sliderContainer.css", "/css/editImages.css", "/css/dialog.css"},
		JsFiles: []string{
			"/js/editImages.js",
			"/js/get.js",
			"/js/post.js",
			"/js/delete.js",
		},
		PageName: "images",
		ContainerData: sliderContainerData{
			ContentData: nil,
		},
	})
}

func apiDeleteImage(w http.ResponseWriter, r *http.Request) {
	var id int64
	var err error
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	idStr := string(buf)
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transaction, err := db.DB.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	_, err = transaction.Exec(
		`DELETE FROM images
		WHERE Iid = ?;`,
		id,
	)
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	transaction.Commit()
}

func apiUploadImage(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var imageData ImageData
	err = json.Unmarshal(buf, &imageData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transaction, err := db.DB.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	res, err := transaction.Exec(
		`INSERT INTO images (Isuffix)
			VALUES (?);`,
		imageData.Suffix,
	)
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	f, err := os.OpenFile(
		global.ROOT_PATH+"wwwroot/images/"+fmt.Sprint(id)+imageData.Suffix,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	data, err := base64.StdEncoding.DecodeString(imageData.Data)
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	_, err = f.Write(data)
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	err = f.Close()
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	transaction.Commit()
}

func registerEditImagesRoutes() {
	template.Must(
		editImagesTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/adminLayout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/editImages.html",
		),
	)
}
