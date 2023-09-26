package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/3ssalunke/leetcode-clone-app/db"
	"github.com/markbates/goth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OAuthSignUp(ctx context.Context, database db.Database, oauthUser goth.User) error {
	filter := bson.M{
		"$or": []bson.M{
			{"email": oauthUser.Email},
			{"username": oauthUser.Name},
		},
	}

	var existingUser db.User

	err := database.Collection("users").FindOne(ctx, filter).Decode(&existingUser)
	if err == nil {
		return fmt.Errorf("user already exist for given username or email")
	}

	trimmedUsername := strings.Trim(oauthUser.Name, " ")
	hypenatedUsername := strings.ReplaceAll(trimmedUsername, " ", "-")

	user := &db.User{
		ID:         primitive.NewObjectID(),
		Username:   hypenatedUsername,
		Email:      oauthUser.Email,
		Password:   "",
		ImageUrl:   "",
		Created_at: primitive.NewDateTimeFromTime(time.Now()),
		Updated_at: primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err = database.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert new user - %w", err)
	}

	return nil
}
