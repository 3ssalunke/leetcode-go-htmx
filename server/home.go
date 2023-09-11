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

	var data struct {
		Title  string
		UserID string
		User   *db.User
	}
	data.Title = "LeetCode - The World's Leading Programming Learning Platform"

	jwtToken, err := util.ExtractCookieFromHeader(r)
	if err != nil {
		log.Printf("failed to extract cookie from request header %v", err)
		data.UserID = ""
		data.User = nil
	} else {
		payload, tokenerr := server.tokenMaker.VerifyToken(jwtToken)
		if tokenerr != nil {
			log.Printf("failed to verify token %v", tokenerr)
			data.UserID = ""
			data.User = nil
		} else {
			var user db.User
			dberr := server.db.Collection("users").FindOne(ctx, bson.M{"username": payload.Username}).Decode(&user)
			if dberr != nil {
				log.Printf("failed to fetch user details %v", dberr)
				data.UserID = ""
				data.User = nil
			} else {
				data.UserID = user.ID.String()
				data.User = &user
			}
		}
	}

	log.Println(data)
	t.Execute(w, data)
}
