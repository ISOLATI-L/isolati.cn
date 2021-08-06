package controller

import (
	"WEB_ISOLATI/constant_define"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type Video struct {
	Vid      int
	Vtitle   string
	Vcontent template.HTML
	Vcover   string
	Vtime    string
}

var videosTemplate = template.New("videos")
var videoTemplate = template.New("video")

var videosPattern *regexp.Regexp

func handleVideos(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/videos" ||
		r.URL.Path == "/videos/" {
		rows, err := constant_define.DB.Query(
			"SELECT Vid, Vtitle, Vcontent, Vcover, Vtime FROM videos")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		defer rows.Close()

		var videos []Video
		var video Video
		for rows.Next() {
			video = Video{}
			err = rows.Scan(
				&video.Vid,
				&video.Vtitle,
				&video.Vcontent,
				&video.Vcover,
				&video.Vtime,
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			videos = append([]Video{video}, videos...)
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
	// http.HandleFunc("/videos/", handleVideos)
}
