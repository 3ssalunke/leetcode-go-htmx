package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/3ssalunke/leetcode-clone-app/db"
	"github.com/3ssalunke/leetcode-clone-app/middleware"
	"github.com/3ssalunke/leetcode-clone-app/util"
)

func (server *Server) home(w http.ResponseWriter, r *http.Request) {
	layoutsDir, err := util.GetTemplateDir()
	if err != nil {
		log.Printf("failed to get view template directory - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var headerTemplate string

	var data struct {
		Title  string
		UserID string
		User   *db.User
	}
	data.Title = "LeetCode - The World's Leading Programming Learning Platform"

	user, ok := r.Context().Value(middleware.ContextUserKey).(db.User)
	if !ok {
		log.Printf("user is not authorized")
		data.UserID = ""
		data.User = nil

		headerTemplate = fmt.Sprint(layoutsDir, "\\common\\home_header.html")
	} else {
		data.UserID = user.ID.String()
		data.User = &user

		headerTemplate = fmt.Sprint(layoutsDir, "\\common\\main_header.html")
	}

	mainTemplate := fmt.Sprint(layoutsDir, "\\common\\base.html")
	homeTemplate := fmt.Sprint(layoutsDir, "\\common\\home.html")

	t, err := template.ParseFiles(mainTemplate, headerTemplate, homeTemplate)
	if err != nil {
		log.Printf("failed to parse view templates - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}
