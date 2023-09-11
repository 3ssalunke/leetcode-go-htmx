package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/3ssalunke/leetcode-clone/db"
	"github.com/3ssalunke/leetcode-clone/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(ctx context.Context, database db.Database, r *http.Request) (*db.User, error) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	var existingUser db.User

	err := database.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return nil, fmt.Errorf("user already exists with given username")
	}

	err = database.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		return nil, fmt.Errorf("user already exists with given email")
	}

	if password != confirmPassword {
		return nil, fmt.Errorf("entered passwords does not match")
	}

	password, err = util.HashPassword(confirmPassword)
	if err != nil {
		return nil, err
	}

	user := &db.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Email:    email,
		Password: password,
		ImageUrl: "",
	}

	_, err = database.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to insert new user %w", err)
	}

	return user, nil
}

func SignIn(ctx context.Context, database db.Database, r *http.Request) (*db.User, error) {
	usernameOrEmail := r.FormValue("usernameOrEmail")
	password := r.FormValue("password")

	var user db.User
	filter := bson.M{
		"$or": []bson.M{
			{"email": usernameOrEmail},
			{"username": usernameOrEmail},
		},
	}
	err := database.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Printf("user does not exist for given username or email %v", err)
		return nil, fmt.Errorf("user does not exist for given username or email")
	}

	err = util.CheckPassword(user.Password, password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
