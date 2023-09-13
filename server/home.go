package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/3ssalunke/leetcode-clone/db"
	"github.com/3ssalunke/leetcode-clone/middleware"
	"github.com/3ssalunke/leetcode-clone/util"
)

func (server *Server) home(w http.ResponseWriter, r *http.Request) {
	layoutsDir, err := util.GetTemplateDir()
	if err != nil {
		log.Printf("failed to get view template directory %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mainTemplate := fmt.Sprint(layoutsDir, "\\common\\base.html")
	headerTemplate := fmt.Sprint(layoutsDir, "\\common\\header.html")
	homeTemplate := fmt.Sprint(layoutsDir, "\\common\\home.html")

	t, err := template.ParseFiles(mainTemplate, headerTemplate, homeTemplate)
	if err != nil {
		log.Printf("failed to parse view templates %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
	} else {
		data.UserID = user.ID.String()
		data.User = &user
	}

	t.Execute(w, data)
}
