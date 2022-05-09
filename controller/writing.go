package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"isolati.cn/database"
	"isolati.cn/db"
	"isolati.cn/global"
)

type ParagraphData struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var writingTemplate = template.New("writing")

func showWritingPage(w http.ResponseWriter, r *http.Request) {
	paragraphsID := r.URL.Query().Get("p")
	var paragraphData ParagraphData
	if len(paragraphsID) > 0 {
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
			`SELECT Pid, Ptitle FROM paragraphs
			WHERE Pid=?;`,
			paragraphsID,
		)
		paragraph := database.Paragraph{}
		row.Scan(
			&paragraph.Pid,
			&paragraph.Ptitle,
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
				global.ROOT_PATH+"wwwroot/md/"+fmt.Sprint(paragraph.Pid)+".md",
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
			paragraphData.Title = paragraph.Ptitle
			paragraphData.Content = string(buf)
		}
	}
	writingTemplate.ExecuteTemplate(w, "layout", layoutMsg{
		CssFiles: []string{
			"/css/writing.css",
			"/css/dialog.css",
			"/editormd/css/editormd.min.css",
		},
		JsFiles: []string{
			"/jquery/jquery-3.6.0.min.js",
			"/editormd/editormd.min.js",
			"/js/writing.js",
			"/js/post.js",
		},
		PageName:      "paragraphs",
		ContainerData: paragraphData,
	})
}

func apiUploadParagraph(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var paragraphData ParagraphData
	err = json.Unmarshal(buf, &paragraphData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var p int64
	pStr := r.URL.Query().Get("p")
	if pStr == "" {
		p = 0
	} else {
		pMatches := pnumberPattern.FindStringSubmatch(pStr)
		if len(pMatches) > 0 {
			p, err = strconv.ParseInt(pMatches[1], 10, 64)
			if err != nil {
				log.Println(err.Error())
				p = 0
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
	if p != 0 {
		_, err = transaction.Exec(
			`UPDATE paragraphs SET Ptitle = ?
			WHERE Pid = ?;`,
			paragraphData.Title,
			p,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		f, err := os.OpenFile(
			global.ROOT_PATH+"wwwroot/md/"+fmt.Sprint(p)+".md",
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0666,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		_, err = f.Write([]byte(paragraphData.Content))
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
	} else {
		res, err := transaction.Exec(
			`INSERT INTO paragraphs (Ptitle)
			VALUES (?);`,
			paragraphData.Title,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		p, err = res.LastInsertId()
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		f, err := os.OpenFile(
			global.ROOT_PATH+"wwwroot/md/"+fmt.Sprint(p)+".md",
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0666,
		)
		if err != nil {
			transaction.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		_, err = f.Write([]byte(paragraphData.Content))
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
	}
	transaction.Commit()
}

func registerWritingRoutes() {
	template.Must(
		writingTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/adminLayout.html",
			// global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/writing.html",
		),
	)
}
