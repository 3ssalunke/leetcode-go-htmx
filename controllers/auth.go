package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/3ssalunke/leetcode-clone/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(ctx context.Context, database db.Database, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	fmt.Println(username, email, password, confirmPassword)

	user := db.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Email:    email,
		Password: password,
		ImageUrl: "",
	}
	database.Collection("users").InsertOne(ctx, user)
}
