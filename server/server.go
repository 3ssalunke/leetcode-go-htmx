package server

import (
	"log"
	"net/http"
	"time"

	"github.com/3ssalunke/leetcode-clone/db"
	"github.com/3ssalunke/leetcode-clone/util"
	"github.com/gorilla/mux"
)

type Server struct {
	http.Server
	config util.Config
	db     db.Database
}

func NewServer() *Server {
	server := &Server{}
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Error occured while loading enviroment variables", err)
	}

	server.config = config
	server.db = db.NewMongoDatabase(config)
	server.Handler = server.setupRoutes()
	server.WriteTimeout = 15 * time.Second
	server.ReadTimeout = 15 * time.Second

	return server
}

func (server *Server) Start(addr string) error {
	server.Addr = addr
	log.Println("server listening on", addr)
	return server.ListenAndServe()
}

func (server *Server) setupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/static/images").Handler(http.StripPrefix("/static/images", http.FileServer(http.Dir("public/images"))))
	r.PathPrefix("/static/css").Handler(http.StripPrefix("/static/css", http.FileServer(http.Dir("public/css"))))

	r.Use(util.LoggerMiddleware)

	r.HandleFunc("/", server.home).Methods("GET")
	r.HandleFunc("/signin", server.signIn).Methods("GET", "POST")
	r.HandleFunc("/signup", server.signUp).Methods("GET", "POST")

	return r
}
