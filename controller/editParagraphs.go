package controller

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"isolati.cn/database"
	"isolati.cn/db"
	"isolati.cn/global"
)

var editParagraphsTemplate = template.New("editParagraphs")

func showEditParagraphsPage(w http.ResponseWriter, r *http.Request) {
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
		TotalPage:  totalPage,
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

	editParagraphsTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		CssFiles: []string{"/css/sliderContainer.css", "/css/editParagraphs.css", "/css/dialog.css"},
		JsFiles:  []string{"/js/editParagraphs.js", "/js/delete.js"},
		PageName: "paragraphs",
		ContainerData: sliderContainerData{
			ContentData: paragraphs,
		},
	})
}

func apiDeleteParagraph(w http.ResponseWriter, r *http.Request) {
	var p int64
	var err error
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	pStr := string(buf)
	if pStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err = strconv.ParseInt(pStr, 10, 64)
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
		`DELETE FROM paragraphs
		WHERE Pid = ?;`,
		p,
	)
	if err != nil {
		transaction.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	transaction.Commit()
}

func registerEditParagraphsRoutes() {
	template.Must(
		editParagraphsTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/adminLayout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/editParagraphs.html",
		),
	)
}
