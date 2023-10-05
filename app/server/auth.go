package server

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/3ssalunke/leetcode-clone-app/db"
	"github.com/3ssalunke/leetcode-clone-app/middleware"
	"github.com/3ssalunke/leetcode-clone-app/services"
	"github.com/3ssalunke/leetcode-clone-app/util"
)

func (server *Server) signIn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	layoutsDir, err := util.GetTemplateDir()
	if err != nil {
		log.Printf("failed to get view template directory - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	baseTemplate := fmt.Sprint(layoutsDir, "\\common\\base.html")
	headerTemplate := fmt.Sprint(layoutsDir, "\\common\\main_header.html")
	authBaseTemplate := fmt.Sprint(layoutsDir, "\\auth\\auth_base.html")
	authSignInTemplate := fmt.Sprint(layoutsDir, "\\auth\\signin.html")

	t, err := template.ParseFiles(baseTemplate, headerTemplate, authBaseTemplate, authSignInTemplate)
	if err != nil {
		log.Printf("failed to parse view templates - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		_, ok := r.Context().Value(middleware.ContextUserKey).(db.User)
		if ok {
			log.Println("User has valid logged in session")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		data := struct {
			Title   string
			Message string
			UserID  string
		}{Title: "Account Login - LeetCode", Message: "", UserID: ""}

		t.Execute(w, data)
		return
	} else {
		user, err := services.SignIn(ctx, server.Db, r)
		if err != nil {
			data := struct {
				Status  int
				Message string
			}{Status: http.StatusBadRequest, Message: err.Error()}

			var signInOutputBuffer bytes.Buffer

			err = t.ExecuteTemplate(&signInOutputBuffer, "auth_form", data)
			if err != nil {
				log.Printf("failed to parse auth form template - %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(signInOutputBuffer.Bytes())
			return
		}

		token, err := server.tokenMaker.CreateToken(user.Username, user.Email, time.Duration(24*time.Hour))
		if err != nil {
			log.Printf("failed to create token - %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authCookie := &http.Cookie{
			Name:     "leetcode_auth",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		}

		http.SetCookie(w, authCookie)
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(301)
		return
	}
}

func (server *Server) signUp(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	layoutsDir, err := util.GetTemplateDir()
	if err != nil {
		log.Printf("failed to get view template directory - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	signUpTemplate := fmt.Sprint(layoutsDir, "\\auth\\signup.html")

	t, err := template.ParseFiles(signUpTemplate)
	if err != nil {
		log.Printf("failed to parse auth form template - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		data := struct {
			Message string
		}{Message: ""}
		t.Execute(w, data)
		return
	} else {
		user, err := services.SignUp(ctx, server.Db, r)
		if err != nil {
			data := struct {
				Status  int
				Message string
			}{Status: http.StatusBadRequest, Message: err.Error()}
			t.Execute(w, data)
			return
		}

		token, err := server.tokenMaker.CreateToken(user.Username, user.Email, time.Duration(24*time.Hour))
		if err != nil {
			log.Printf("failed to create token - %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authCookie := &http.Cookie{
			Name:     "leetcode_auth",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		}

		http.SetCookie(w, authCookie)
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(301)
		return
	}
}

func (server *Server) logOut(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, ok := r.Context().Value(middleware.ContextUserKey).(db.User)
		if !ok {
			log.Println("User does not have valid logged in session")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		authCookie := &http.Cookie{
			Name:     "leetcode_auth",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
		}

		http.SetCookie(w, authCookie)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
}
