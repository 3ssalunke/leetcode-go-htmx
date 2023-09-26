package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/3ssalunke/leetcode-clone-app/controllers"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func (server *Server) oAuthHandler(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (server *Server) oAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	provider := vars["provider"]

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("failed to get user details from provider - %s, %v", provider, err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, "/accounts/signin", http.StatusTemporaryRedirect)
		return
	}

	err = controllers.OAuthSignUp(ctx, server.Db, user)
	if err != nil {
		log.Printf("failed to store oauth user in db - %v", err)
	}

	token, err := server.tokenMaker.CreateToken(user.Name, user.Email, time.Duration(24*time.Hour))
	if err != nil {
		log.Printf("failed to create token - %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, "/accounts/signin", http.StatusTemporaryRedirect)
		return
	}

	authCookie := &http.Cookie{
		Name:     "leetcode_auth",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, authCookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
