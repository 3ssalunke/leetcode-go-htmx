package controllers

import (
	"context"

	"github.com/3ssalunke/leetcode-clone/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProblems(ctx context.Context, database db.Database, userId primitive.ObjectID) ([]db.Problem, error) {
	var problems []db.Problem
	cursor, err := database.Collection("problem_set").Find(ctx, bson.M{})
	if err != nil {
		return problems, err
	}
	err = cursor.All(ctx, &problems)
	if problems == nil {
		return problems, err
	}

	return problems, err
}
