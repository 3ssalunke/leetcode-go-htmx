package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func GetTemplateDir() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("error occured while getting current directory path", err)
		return "", err
	}
	return filepath.Join(pwd, "views"), nil
}

func (server *Server) home(w http.ResponseWriter, r *http.Request) {
	layoutsDir, err := GetTemplateDir()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mainTemplate := fmt.Sprint(layoutsDir, "\\base.html")
	headerTemplate := fmt.Sprint(layoutsDir, "\\header.html")
	homeTemplate := fmt.Sprint(layoutsDir, "\\home.html")

	data := struct{ Title string }{Title: "LeetCode - The World's Leading Programming Learning Platform"}
	t, err := template.ParseFiles(mainTemplate, headerTemplate, homeTemplate)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}
