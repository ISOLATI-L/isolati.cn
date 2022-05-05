package controller

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"isolati.cn/database"
	"isolati.cn/db"
	"isolati.cn/global"
)

var imagesTemplate = template.New("images")

var imagesPattern *regexp.Regexp

func showImagePage(w http.ResponseWriter, r *http.Request) {
	matches := imagesPattern.FindStringSubmatch(r.URL.Path)
	if len(matches) > 0 {
		transaction, err := db.DB.Begin()
		if err != nil {
			if transaction != nil {
				transaction.Rollback()
			}
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		row := transaction.QueryRow(
			`SELECT Iid, Iname FROM images
			WHERE Iid=?;`,
			matches[1],
		)
		image := database.Image{}
		row.Scan(
			&image.Iid,
			&image.Iname,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		transaction.Commit()
		// log.Println(image)
		if image.Iid != 0 {
			http.ServeFile(w, r, global.ROOT_PATH+"/wwwroot/images/"+image.Iname)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func showImagesPage(w http.ResponseWriter, r *http.Request) {
	imagesTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		CssFiles: []string{"/css/sliderContainer.css", "/css/images.css"},
		JsFiles:  []string{"/js/images.js", "/js/get.js"},
		PageName: "images",
		ContainerData: sliderContainerData{
			LeftSliderData:  global.LEFT_SLIDER,
			RightSliderData: global.RIGHT_SLIDER,
			ContentData:     nil,
		},
	})
}

func apiGetList(w http.ResponseWriter, r *http.Request) {
	var err error
	transaction, err := db.DB.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	var rows *sql.Rows
	rows, err = transaction.Query(
		`SELECT Iid FROM images
			ORDER BY Iid DESC;`,
	)
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	defer rows.Close()

	images := make([]string, 0)
	for rows.Next() {
		var imageName string
		err = rows.Scan(
			&imageName,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		images = append(images, imageName)
	}
	transaction.Commit()

	res, err := json.Marshal(images)
	w.Header().Set("Content-type", "application/json")
	w.Write(res)
}

func handleImages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/images", "/images/":
			showImagesPage(w, r)
		case "/images/api/list":
			apiGetList(w, r)
		default:
			showImagePage(w, r)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func registerImagesRoutes() {
	var err error
	imagesPattern, err = regexp.Compile(`/images/(\d+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	template.Must(
		imagesTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/images.html",
		),
	)
	registerWritingRoutes()
	http.HandleFunc("/images", handleImages)
	http.HandleFunc("/images/", handleImages)
}
