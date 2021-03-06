package controller

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"isolati.cn/database"
	"isolati.cn/db"
	"isolati.cn/global"
)

type VideoList struct {
	Videos    []database.Video
	Page      int64
	TotalPage int64
}

const MAX_VIDEOS_PER_PAGE = 10

var videosTemplate = template.New("videos")
var videoTemplate = template.New("video")

var videosPattern *regexp.Regexp
var vnumberPattern *regexp.Regexp

func showVideoPage(w http.ResponseWriter, r *http.Request) {
	matches := videosPattern.FindStringSubmatch(r.URL.Path)
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
			`SELECT Vid, Vcontent FROM videos
			WHERE Vid=?;`,
			matches[1],
		)
		video := database.Video{}
		row.Scan(
			&video.Vid,
			&video.Vcontent,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		transaction.Commit()
		// log.Println(video)
		if video.Vid != 0 {
			videoTemplate.ExecuteTemplate(w, "layout", layoutMsg{
				CssFiles: []string{"/css/sliderContainer.css", "/css/video.css"},
				JsFiles:  []string{},
				PageName: "videos",
				ContainerData: sliderContainerData{
					ContentData: video.Vcontent,
				},
			})
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func showVideosPage(w http.ResponseWriter, r *http.Request) {
	var page int64
	var err error
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		page = 1
	} else {
		pageMatches := vnumberPattern.FindStringSubmatch(pageStr)
		if len(pageMatches) > 0 {
			page, err = strconv.ParseInt(pageMatches[1], 10, 64)
			if err != nil {
				log.Println(err.Error())
				page = 1
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	// log.Println(page)
	transaction, err := db.DB.Begin()
	if err != nil {
		if transaction != nil {
			transaction.Rollback()
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	row := transaction.QueryRow(`SELECT COUNT(Vid) FROM videos;`)
	var totalPage int64
	err = row.Scan(
		&totalPage,
	)
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	totalPage = (totalPage + MAX_VIDEOS_PER_PAGE - 1) / MAX_VIDEOS_PER_PAGE
	if totalPage == 0 {
		totalPage = 1
	}
	// log.Println(totalPage)
	if page < 1 || page > totalPage {
		transaction.Rollback()
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var rows *sql.Rows
	rows, err = transaction.Query(
		`SELECT Vid, Vtitle, Vcover, Vtime FROM videos
			ORDER BY Vid DESC LIMIT ?, ?;`,
		(page-1)*MAX_VIDEOS_PER_PAGE,
		MAX_VIDEOS_PER_PAGE,
	)
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	defer rows.Close()

	videos := VideoList{
		Videos:    []database.Video{},
		Page:      page,
		TotalPage: totalPage,
	}
	var video database.Video
	for rows.Next() {
		video = database.Video{}
		var timeStr string
		err = rows.Scan(
			&video.Vid,
			&video.Vtitle,
			&video.Vcover,
			&timeStr,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		video.Vtime, err = time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		videos.Videos = append(videos.Videos, video)
		// log.Println(video)
	}
	transaction.Commit()
	// log.Println("Done!")

	videosTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		CssFiles: []string{"/css/sliderContainer.css", "/css/videos.css"},
		JsFiles:  []string{},
		PageName: "videos",
		ContainerData: sliderContainerData{
			ContentData: videos,
		},
	})
}

func handleVideos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/videos", "/videos/":
			showVideosPage(w, r)
		default:
			showVideoPage(w, r)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func registerVideosRoutes() {
	var err error
	videosPattern, err = regexp.Compile(`/videos/(\d+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	vnumberPattern, err = regexp.Compile(`^(\d+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	template.Must(
		videosTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/videos.html",
		),
	)
	template.Must(
		videoTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/video.html",
		),
	)
	http.HandleFunc("/videos", handleVideos)
	http.HandleFunc("/videos/", handleVideos)
}
