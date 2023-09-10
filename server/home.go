package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/3ssalunke/leetcode-clone/db"
	"github.com/3ssalunke/leetcode-clone/util"
	"go.mongodb.org/mongo-driver/bson"
)

func (server *Server) home(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	jwtToken, err := util.ExtractCookieFromHeader(r)
	if err != nil {
		log.Printf("failed to extract cookie from request header %v", err)
		http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)
		return
	} else {
		payload, tokenerr := server.tokenMaker.VerifyToken(jwtToken)
		if tokenerr != nil {
			log.Printf("failed to verify token %v", tokenerr)
			http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)
			return
		}

		var user db.User
		dberr := server.db.Collection("users").FindOne(ctx, bson.M{"username": payload.Username}).Decode(&user)
		if dberr != nil {
			log.Printf("failed to fetch user details %v", dberr)
			http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)
			return
		}

		data := struct {
			Title string
			User  db.User
		}{Title: "LeetCode - The World's Leading Programming Learning Platform", User: user}
		t.Execute(w, data)
	}
}
