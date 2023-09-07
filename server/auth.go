package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/3ssalunke/leetcode-clone/controllers"
)

func (server *Server) signIn(w http.ResponseWriter, r *http.Request) {
	layoutsDir, err := GetTemplateDir()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mainTemplate := fmt.Sprint(layoutsDir, "\\base.html")
	headerTemplate := fmt.Sprint(layoutsDir, "\\header.html")
	signInTemplate := fmt.Sprint(layoutsDir, "\\signin.html")

	if r.Method == "GET" {
		data := struct{ Title string }{Title: "Account Login - LeetCode"}
		t, err := template.ParseFiles(mainTemplate, headerTemplate, signInTemplate)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Execute(w, data)
	} else {
		username := r.FormValue("username")
		password := r.FormValue("password")

		fmt.Println(username, password)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (server *Server) signUp(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	layoutsDir, err := GetTemplateDir()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	signUpTemplate := fmt.Sprint(layoutsDir, "\\signup.html")

	if r.Method == "GET" {
		data := struct{ Title string }{Title: "Account Login - LeetCode"}
		t, err := template.ParseFiles(signUpTemplate)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.Execute(w, data)
	} else {
		controllers.SignUp(ctx, server.db, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
