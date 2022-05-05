package controller

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"isolati.cn/database"
	"isolati.cn/db"
	"isolati.cn/global"
)

type ParagraphsList struct {
	Paragraphs []database.Paragraph
	Page       int64
}

const MAX_PARAGRAPHS_PER_PAGE = 10

var paragraphsTemplate = template.New("paragraphs")
var paragraphTemplate = template.New("paragraph")

var paragraphsPattern *regexp.Regexp
var pnumberPattern *regexp.Regexp

func showParagraphPage(w http.ResponseWriter, r *http.Request) {
	matches := paragraphsPattern.FindStringSubmatch(r.URL.Path)
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
			`SELECT Pid, Pcontent FROM paragraphs
			WHERE Pid=?;`,
			matches[1],
		)
		paragraph := database.Paragraph{}
		row.Scan(
			&paragraph.Pid,
			&paragraph.Pcontent,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		transaction.Commit()
		// log.Println(paragraph)
		if paragraph.Pid != 0 {
			f, err := os.OpenFile(
				global.ROOT_PATH+"wwwroot"+paragraph.Pcontent,
				os.O_RDONLY,
				0666,
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			buf, err := io.ReadAll(f)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			err = f.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			ContentData := template.HTML(buf)
			paragraphTemplate.ExecuteTemplate(w, "layout", layoutMsg{
				CssFiles: []string{
					"/css/sliderContainer.css",
					"/css/paragraph.css",
					"/editormd/css/editormd.min.css",
				},
				JsFiles: []string{
					"/jquery/jquery-3.6.0.min.js",
					"/editormd/lib/marked.min.js",
					"/editormd/lib/prettify.min.js",

					"/editormd/lib/raphael.min.js",
					"/editormd/lib/underscore.min.js",
					"/editormd/lib/sequence-diagram.min.js",
					"/editormd/lib/flowchart.min.js",
					"/editormd/lib/jquery.flowchart.min.js",

					"/editormd/editormd.min.js",
					"/js/paragraph.js",
				},
				PageName: "paragraphs",
				ContainerData: sliderContainerData{
					LeftSliderData:  global.LEFT_SLIDER,
					RightSliderData: global.RIGHT_SLIDER,
					ContentData:     ContentData,
				},
			})
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func showParagraphsPage(w http.ResponseWriter, r *http.Request) {
	var page int64
	var err error
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		page = 1
	} else {
		pageMatches := pnumberPattern.FindStringSubmatch(pageStr)
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
	row := transaction.QueryRow(`SELECT COUNT(Pid) FROM paragraphs;`)
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
	totalPage = (totalPage + MAX_PARAGRAPHS_PER_PAGE - 1) / MAX_PARAGRAPHS_PER_PAGE
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
		`SELECT Pid, Ptitle, Ptime FROM paragraphs
			ORDER BY Pid DESC LIMIT ?, ?;`,
		(page-1)*MAX_PARAGRAPHS_PER_PAGE,
		MAX_PARAGRAPHS_PER_PAGE,
	)
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	defer rows.Close()

	paragraphs := ParagraphsList{
		Paragraphs: []database.Paragraph{},
		Page:       page,
	}
	var paragraph database.Paragraph
	for rows.Next() {
		paragraph = database.Paragraph{}
		var timeStr string
		err = rows.Scan(
			&paragraph.Pid,
			&paragraph.Ptitle,
			&timeStr,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		paragraph.Ptime, err = time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		paragraphs.Paragraphs = append(paragraphs.Paragraphs, paragraph)
		// log.Println(paragraphs)
	}
	transaction.Commit()
	// log.Println("Done!")

	paragraphsTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		CssFiles: []string{"/css/sliderContainer.css", "/css/paragraphs.css"},
		JsFiles:  []string{},
		PageName: "paragraphs",
		ContainerData: sliderContainerData{
			LeftSliderData:  global.LEFT_SLIDER,
			RightSliderData: global.RIGHT_SLIDER,
			ContentData:     paragraphs,
		},
	})
}

func handleParagraphs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/paragraphs", "/paragraphs/":
			showParagraphsPage(w, r)
		default:
			showParagraphPage(w, r)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func registerParagraphsRoutes() {
	var err error
	paragraphsPattern, err = regexp.Compile(`/paragraphs/(\d+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	pnumberPattern, err = regexp.Compile(`^(\d+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	template.Must(
		paragraphsTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/paragraphs.html",
		),
	)
	template.Must(
		paragraphTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/paragraph.html",
		),
	)
	http.HandleFunc("/paragraphs", handleParagraphs)
	http.HandleFunc("/paragraphs/", handleParagraphs)
}
