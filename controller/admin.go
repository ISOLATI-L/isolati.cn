package controller

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"net/http"

	"isolati.cn/global"
	"isolati.cn/session"
)

var adminTemplate = template.New("admin")

func isRequestAdmin(r *http.Request) (bool, error) {
	transaction, err := session.UserSession.BeginTransaction()
	if err != nil {
		return false, err
	}
	vIdentity, err := session.UserSession.GetByRequest(transaction, r, "identity")
	if err != nil {
		if err == session.ErrNoCookies {
			transaction.Commit()
			return false, nil
		} else {
			transaction.Rollback()
			return false, err
		}
	} else {
		var identity string
		err = json.Unmarshal(vIdentity, &identity)
		if err != nil {
			transaction.Rollback()
			return false, err
		} else {
			transaction.Commit()
			return identity == "admin", nil
		}
	}
}

func showAdminPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/paragraphs", http.StatusSeeOther)
	// adminTemplate.ExecuteTemplate(w, "layout", layoutMsg{
	// 	CssFiles: []string{"/css/sliderContainer.css", "/css/admin.css"},
	// 	JsFiles:  []string{},
	// 	PageName: "admin",
	// 	ContainerData: sliderContainerData{
	// 		ContentData: nil,
	// 	},
	// })
}

func handleAdmin(w http.ResponseWriter, r *http.Request) {
	isAdmin, err := isRequestAdmin(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !isAdmin {
		http.Redirect(w,
			r,
			"/login?ref=\""+base64.URLEncoding.EncodeToString([]byte(r.URL.String()))+"\"",
			http.StatusSeeOther,
		)
		return
	}

	switch r.URL.Path {
	case "/admin", "/admin/":
		switch r.Method {
		case http.MethodGet:
			showAdminPage(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "/admin/writing":
		switch r.Method {
		case http.MethodGet:
			showWritingPage(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "/admin/paragraphs":
		switch r.Method {
		case http.MethodGet:
			showEditParagraphsPage(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "/admin/images":
		switch r.Method {
		case http.MethodGet:
			showEditImagesPage(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "/admin/api/update":
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "/admin/api/paragraph":
		switch r.Method {
		case http.MethodPost:
			apiUploadParagraph(w, r)
		case http.MethodDelete:
			apiDeleteParagraph(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "/admin/api/image":
		switch r.Method {
		case http.MethodPost:
			apiUploadImage(w, r)
		case http.MethodDelete:
			apiDeleteImage(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func registerAdminRoutes() {
	template.Must(
		adminTemplate.ParseFiles(
			global.ROOT_PATH+"/wwwroot/adminLayout.html",
			global.ROOT_PATH+"/wwwroot/sliderContainer.html",
			global.ROOT_PATH+"/wwwroot/admin.html",
		),
	)
	registerEditParagraphsRoutes()
	registerWritingRoutes()
	registerEditImagesRoutes()
	http.HandleFunc("/admin", handleAdmin)
	http.HandleFunc("/admin/", handleAdmin)
}
