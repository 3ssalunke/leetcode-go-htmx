package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/3ssalunke/leetcode-clone/controllers"
	"github.com/3ssalunke/leetcode-clone/db"
	"github.com/3ssalunke/leetcode-clone/middleware"
	"github.com/3ssalunke/leetcode-clone/util"
	"github.com/gorilla/mux"
)

func (server *Server) ProblemsAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data struct {
		Title        string
		UserID       string
		User         *db.User
		ProblemsList []db.Problem
	}

	data.Title = "Problems - LeetCode"

	user, ok := r.Context().Value(middleware.ContextUserKey).(db.User)
	if !ok {
		log.Printf("user is not authorized")
		http.Redirect(w, r, "/accounts/signin", http.StatusTemporaryRedirect)
		return
	} else {
		data.UserID = user.ID.String()
		data.User = &user
	}

	layoutsDir, err := util.GetTemplateDir()
	if err != nil {
		log.Printf("failed to get view template directory - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mainTemplate := fmt.Sprint(layoutsDir, "\\common\\base.html")
	headerTemplate := fmt.Sprint(layoutsDir, "\\common\\main_header.html")
	problemsListTemplate := fmt.Sprint(layoutsDir, "\\problems\\list.html")

	t, err := template.ParseFiles(mainTemplate, headerTemplate, problemsListTemplate)
	if err != nil {
		log.Printf("failed to parse view templates - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	problems, err := controllers.GetProblems(ctx, server.Db, user.ID)
	if err != nil {
		log.Printf("failed to fetch problems - %v", err)
	}
	log.Println(problems[0].Title)
	data.ProblemsList = problems

	t.Execute(w, data)
}

func (server *Server) Problem(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data struct {
		Title   string
		UserID  string
		User    *db.User
		Problem db.Problem
	}

	vars := mux.Vars(r)
	problemSlug := vars["problem"]

	data.Title = problemSlug + " - The World's Leading Programming Learning Platform"

	user, ok := r.Context().Value(middleware.ContextUserKey).(db.User)
	if !ok {
		log.Printf("user is not authorized")
		http.Redirect(w, r, "/accounts/signin", http.StatusTemporaryRedirect)
		return
	} else {
		data.UserID = user.ID.String()
		data.User = &user
	}

	layoutsDir, err := util.GetTemplateDir()
	if err != nil {
		log.Printf("failed to get view template directory - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mainTemplate := fmt.Sprint(layoutsDir, "\\common\\base.html")
	headerTemplate := fmt.Sprint(layoutsDir, "\\common\\main_header.html")
	problemsListTemplate := fmt.Sprint(layoutsDir, "\\problems\\problem.html")

	t, err := template.ParseFiles(mainTemplate, headerTemplate, problemsListTemplate)
	if err != nil {
		log.Printf("failed to parse view templates - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}
