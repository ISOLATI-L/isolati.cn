package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"isolati.cn/constant_define"
)

type Video struct {
	Vid      int
	Vtitle   string
	Vcontent template.HTML
	Vcover   string
	Vtime    string
}

type VideoList struct {
	Videos    []Video
	Page      int64
	TotalPage int64
}

const MAX_PER_PAGE = 10

var videosTemplate = template.New("videos")
var videoTemplate = template.New("video")

var videosPattern *regexp.Regexp
var numberPattern *regexp.Regexp

func handleVideos(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/videos" ||
		r.URL.Path == "/videos/" {
		var page int64
		var err error
		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			page = 1
		} else {
			pageMatches := numberPattern.FindStringSubmatch(pageStr)
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
		row := constant_define.DB.QueryRow(`SELECT COUNT(Vid) FROM videos;`)
		var totalPage int64
		err = row.Scan(
			&totalPage,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		totalPage = (totalPage + MAX_PER_PAGE - 1) / MAX_PER_PAGE
		// log.Println(totalPage)
		if page < 1 || page > totalPage {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		query := fmt.Sprintf(
			`SELECT Vid, Vtitle, Vcover, Vtime FROM videos
			ORDER BY Vid DESC LIMIT %d, %d;`,
			(page-1)*MAX_PER_PAGE,
			MAX_PER_PAGE,
		)
		var rows *sql.Rows
		rows, err = constant_define.DB.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		defer rows.Close()

		videos := VideoList{
			Videos:    []Video{},
			Page:      page,
			TotalPage: totalPage,
		}
		var video Video
		for rows.Next() {
			video = Video{}
			err = rows.Scan(
				&video.Vid,
				&video.Vtitle,
				&video.Vcover,
				&video.Vtime,
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			videos.Videos = append(videos.Videos, video)
			// log.Println(video)
		}
		// log.Println("Done!")

		videosTemplate.ExecuteTemplate(w, "layout", layoutMsg{
			PageName: "videos",
			ContainerData: sliderContainerData{
				LeftSliderData:  constant_define.LEFT_SLIDER,
				RightSliderData: constant_define.RIGHT_SLIDER,
				ContentData:     videos,
			},
		})
	} else {
		matches := videosPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			query := fmt.Sprintf(
				`SELECT Vid, Vcontent FROM videos
				WHERE Vid=%v`,
				matches[1],
			)
			row := constant_define.DB.QueryRow(query)
			video := Video{}
			var err error
			row.Scan(
				&video.Vid,
				&video.Vcontent,
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			// log.Println(video)
			if video.Vid != 0 {
				videoTemplate.ExecuteTemplate(w, "layout", layoutMsg{
					PageName: "videos",
					ContainerData: sliderContainerData{
						LeftSliderData:  constant_define.LEFT_SLIDER,
						RightSliderData: constant_define.RIGHT_SLIDER,
						ContentData:     video.Vcontent,
					},
				})
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func registerVideosRoutes() {
	var err error
	videosPattern, err = regexp.Compile(`/videos/(\d+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	numberPattern, err = regexp.Compile(`^(\d+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	template.Must(
		videosTemplate.ParseFiles(
			constant_define.ROOT_PATH+"/wwwroot/layout.html",
			constant_define.ROOT_PATH+"/wwwroot/sliderContainer.html",
			constant_define.ROOT_PATH+"/wwwroot/videos.html",
		),
	)
	template.Must(
		videoTemplate.ParseFiles(
			constant_define.ROOT_PATH+"/wwwroot/layout.html",
			constant_define.ROOT_PATH+"/wwwroot/sliderContainer.html",
			constant_define.ROOT_PATH+"/wwwroot/video.html",
		),
	)
	http.HandleFunc("/videos", handleVideos)
	http.HandleFunc("/videos/", handleVideos)
}
