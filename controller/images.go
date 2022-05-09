package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"isolati.cn/db"
	"isolati.cn/global"
)

var imagesTemplate = template.New("images")

var imagesPattern *regexp.Regexp
var inumberPattern *regexp.Regexp

func showImagePage(w http.ResponseWriter, r *http.Request) {
	matches := imagesPattern.FindStringSubmatch(r.URL.Path)
	if len(matches) > 0 {
		http.ServeFile(w, r, global.ROOT_PATH+"/wwwroot/images/"+matches[1])
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
			ContentData: nil,
		},
	})
}

func apiGetImagesList(w http.ResponseWriter, r *http.Request) {
	var s, n int64
	var err error
	sStr := r.URL.Query().Get("s")
	if sStr == "" {
		s = 0
	} else {
		sMatches := pnumberPattern.FindStringSubmatch(sStr)
		if len(sMatches) > 0 {
			s, err = strconv.ParseInt(sMatches[1], 10, 64)
			if err != nil {
				log.Println(err.Error())
				s = 0
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	nStr := r.URL.Query().Get("n")
	if nStr == "" {
		n = -1
	} else {
		nMatches := pnumberPattern.FindStringSubmatch(nStr)
		if len(nMatches) > 0 {
			n, err = strconv.ParseInt(nMatches[1], 10, 64)
			if err != nil {
				log.Println(err.Error())
				n = -1
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
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
	var rows *sql.Rows
	rows, err = transaction.Query(
		`SELECT Iid, Isuffix FROM images
			ORDER BY Iid DESC LIMIT ?, ?;`,
		s,
		n,
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
		var imageID uint64
		var imageSuffix string
		err = rows.Scan(
			&imageID,
			&imageSuffix,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		images = append(images, fmt.Sprintf("%d%s", imageID, imageSuffix))
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
			apiGetImagesList(w, r)
		default:
			showImagePage(w, r)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func registerImagesRoutes() {
	var err error
	imagesPattern, err = regexp.Compile(`/images/(.+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	inumberPattern, err = regexp.Compile(`^(\d+)$`)
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
