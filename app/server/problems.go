package server

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/3ssalunke/leetcode-clone-app/db"
	"github.com/3ssalunke/leetcode-clone-app/middleware"
	"github.com/3ssalunke/leetcode-clone-app/services"
	"github.com/3ssalunke/leetcode-clone-app/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	problems, err := services.GetProblems(ctx, server.Db, user.ID)
	if err != nil {
		log.Printf("failed to fetch problems - %v", err)
	}

	data.ProblemsList = problems

	t.Execute(w, data)
}

type ProblemData struct {
	ID           primitive.ObjectID
	Title        string
	Slug         string
	Content      template.HTML
	TestCaseList []string
	CodeSnippets []db.CodeSnippet
}
type ProblemViewData struct {
	Title   string
	UserID  string
	User    *db.User
	Problem ProblemData
}

func (server *Server) Problem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data ProblemViewData

	vars := mux.Vars(r)
	problemSlug := vars["problem"]

	data.Title = problemSlug + " - LeetCode"

	user, ok := r.Context().Value(middleware.ContextUserKey).(db.User)
	if !ok {
		log.Printf("user is not authorized")
		http.Redirect(w, r, "/accounts/signin", http.StatusTemporaryRedirect)
		return
	} else {
		data.UserID = user.ID.String()
		data.User = &user
	}

	problems, err := services.GetProblemBySlug(ctx, server.Db, problemSlug, user.ID)
	if err != nil {
		log.Printf("failed to fetch details for problem with slug %s - %v", problemSlug, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(problems) == 0 {
		log.Printf("no details for problem with slug %s", problemSlug)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data.Problem = ProblemData{
		ID:           problems[0].ID,
		Title:        problems[0].Title,
		Slug:         problems[0].Slug,
		Content:      template.HTML(problems[0].Content),
		TestCaseList: problems[0].TestCaseList,
		CodeSnippets: problems[0].CodeSnippets,
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

func (server *Server) RunCode(w http.ResponseWriter, r *http.Request) {
	requestPayloadBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read the request payload - %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	requestPayloadString := string(requestPayloadBytes)
	if err := server.Mq.PublishMessage(server.config.RabbitMQQueueName, requestPayloadString); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
