package controller

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"isolati.cn/global"
)

var filesTemplate = template.New("about")

var fileHandler = http.StripPrefix(
	"/files/",
	http.FileServer(http.Dir(global.SHARE_FILES_PATH)),
)
var filePattern *regexp.Regexp
var slashPattern *regexp.Regexp

type fileOrDir struct {
	Name  string
	IsDir bool
}

type fileView struct {
	Path  string
	Files []fileOrDir
}

func handleFiles(w http.ResponseWriter, r *http.Request) {
	matches := filePattern.FindStringSubmatch(r.URL.Path)
	var filePath = global.SHARE_FILES_PATH
	if len(matches) > 0 {
		filePath += "/" + matches[1]
	}
	s, err := os.Stat(filePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if s.IsDir() {
		rd, err := ioutil.ReadDir(filePath)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		filesList := []fileView{}
		if r.URL.Path != "/files/" && r.URL.Path != "/files" {
			filesList = append(filesList, fileView{
				Path: slashPattern.ReplaceAllLiteralString(
					strings.Replace(filepath.Dir(r.URL.Path),
						`\`, "/", -1),
					"/",
				),
				Files: []fileOrDir{},
			})
			filesList[0].Files = append(filesList[0].Files, fileOrDir{
				Name:  "..",
				IsDir: true,
			})
		}

		filesList = append(filesList, fileView{
			Path:  slashPattern.ReplaceAllLiteralString(r.URL.Path+"/", "/"),
			Files: []fileOrDir{},
		})
		for _, fi := range rd {
			filesList[len(filesList)-1].Files =
				append(filesList[len(filesList)-1].Files, fileOrDir{
					Name:  fi.Name(),
					IsDir: fi.IsDir(),
				})
		}
		filesTemplate.ExecuteTemplate(w, "layout", layoutMsg{
			PageName: "files",
			ContainerData: sliderContainerData{
				LeftSliderData:  global.LEFT_SLIDER,
				RightSliderData: global.RIGHT_SLIDER,
				ContentData:     filesList,
			},
		})
	} else {
		fileHandler.ServeHTTP(w, r)
	}
}

func registerFilesRoutes() {
	template.Must(
		filesTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/layout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/files.html",
		),
	)
	var err error
	filePattern, err = regexp.Compile(`/files/(.+)$`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	slashPattern, err = regexp.Compile(`/+`)
	if err != nil {
		log.Fatalln(err.Error())
	}
	http.HandleFunc("/files", handleFiles)
	http.HandleFunc("/files/", handleFiles)
}
