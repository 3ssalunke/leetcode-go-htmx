package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/3ssalunke/leetcode-clone/controllers"
	"github.com/3ssalunke/leetcode-clone/db"
	"github.com/3ssalunke/leetcode-clone/docker"
	"github.com/3ssalunke/leetcode-clone/middleware"
	"github.com/3ssalunke/leetcode-clone/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getFileExtension(lang string) string {
	switch lang {
	case "javascript":
		return "js"
	case "python":
		return "py"
	default:
		return ""
	}
}

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

	data.ProblemsList = problems

	t.Execute(w, data)
}

func (server *Server) Problem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data struct {
		Title   string
		UserID  string
		User    *db.User
		Problem struct {
			ID      primitive.ObjectID
			Title   string
			Slug    string
			Content template.HTML
		}
	}

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

	problems, err := controllers.GetProblemBySlug(ctx, server.Db, problemSlug, user.ID)
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

	data.Problem = struct {
		ID      primitive.ObjectID
		Title   string
		Slug    string
		Content template.HTML
	}{ID: problems[0].ID, Title: problems[0].Title, Slug: problems[0].Slug, Content: template.HTML(problems[0].Content)}

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

type CodeExecuteRequest struct {
	ProblemID string `json:"problem_id"`
	TypedCode string `json:"typed_code"`
	Lang      string `json:"lang"`
}

func (server *Server) ExecuteCode(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
	defer cancel()

	var requestData CodeExecuteRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = util.WriteFile(requestData.Lang, getFileExtension(requestData.Lang), requestData.TypedCode)
	if err != nil {
		log.Printf("failed to write to the file - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		log.Printf("failed to create docker client - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dockerClient.Client.Close()

	dockerImageTag := fmt.Sprintf("%s-docker-img", requestData.Lang)

	image, err := dockerClient.CreateDockerImage(ctx, requestData.Lang, dockerImageTag)
	if err != nil {
		log.Printf("failed to create docker image - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer image.Body.Close()
	log.Println("docker image created successfully.")

	container, err := dockerClient.RunDockerContainer(ctx, dockerImageTag)
	if err != nil {
		log.Printf("failed to run docker container - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("docker container started successfully.")

	reader, err := dockerClient.GetContainerLogs(ctx, container.ID)
	if err != nil {
		log.Printf("failed to get docker container logs - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	if err = dockerClient.RemoveDockerContainer(ctx, container.ID); err != nil {
		log.Printf("failed to remove docker container - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("container removed successfully.")

	if err = dockerClient.RemoveDockerImage(ctx, dockerImageTag); err != nil {
		log.Printf("failed to remove docker image - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("image removed successfully.")

	w.WriteHeader(http.StatusOK)
}
