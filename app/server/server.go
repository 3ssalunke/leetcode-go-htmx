package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/3ssalunke/leetcode-clone-app/db"
	"github.com/3ssalunke/leetcode-clone-app/middleware"
	"github.com/3ssalunke/leetcode-clone-app/mq"
	"github.com/3ssalunke/leetcode-clone-app/token"
	"github.com/3ssalunke/leetcode-clone-app/util"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

type Server struct {
	http.Server
	tokenMaker *token.TokenMaker
	config     util.Config
	Db         db.Database
	Mq         *mq.RabbitMQ
}

func NewServer() *Server {
	server := &Server{}

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load enviroment variables - %v", err)
	}

	tokenMaker, err := token.NewTokenMaker(config.TokenSecret)
	if err != nil {
		log.Fatalf("failed to create token maker - %v", err)
	}

	rabbitmq, err := mq.NewRabbitMQ(config)
	if err != nil {
		log.Fatalf("failed to make connection to rabbitmq server - %v", err)
	}

	err = rabbitmq.CreateChannel()
	if err != nil {
		log.Fatalf("failed to create channel to rabbitmq server - %v", err)
	}

	server.WriteTimeout = 15 * time.Second
	server.ReadTimeout = 15 * time.Second
	server.config = config
	server.Db = db.NewMongoDatabase(config)
	server.tokenMaker = tokenMaker
	server.Handler = server.setupRoutes()
	server.Mq = rabbitmq

	return server
}

func (server *Server) Start(host string) error {
	server.Addr = fmt.Sprintf("%s:%s", host, server.config.AppPort)
	log.Println("server listening on", server.Addr)
	return server.ListenAndServe()
}

func (server *Server) setupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/static/images").Handler(http.StripPrefix("/static/images", http.FileServer(http.Dir("public/images"))))
	r.PathPrefix("/static/css").Handler(http.StripPrefix("/static/css", http.FileServer(http.Dir("public/css"))))

	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.AuthMiddleware(server.tokenMaker, server.Db))

	callback_uri := "http://127.0.0.1:8080/accounts/auth/google/callback"
	goth.UseProviders(
		google.New(server.config.GoogleOAuthClientId, server.config.GoogleOAuthClientSecret, callback_uri, "email", "profile"),
	)

	r.HandleFunc("/", server.home).Methods("GET")

	r.HandleFunc("/accounts/signin", server.signIn).Methods("GET", "POST")
	r.HandleFunc("/accounts/signup", server.signUp).Methods("GET", "POST")
	r.HandleFunc("/accounts/logout", server.logOut).Methods("GET")
	r.HandleFunc("/accounts/auth/{provider}", server.oAuthHandler).Methods("GET")
	r.HandleFunc("/accounts/auth/{provider}/callback", server.oAuthCallbackHandler).Methods("GET")

	r.HandleFunc("/problemset/all", server.ProblemsAll).Methods("GET")
	r.HandleFunc("/problems/{problem}", server.Problem).Methods("GET")
	r.HandleFunc("/problems/run", server.RunCode).Methods("POST")

	return r
}
