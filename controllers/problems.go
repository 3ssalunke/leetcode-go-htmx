package controllers

import (
	"context"

	"github.com/3ssalunke/leetcode-clone/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProblemWithDetails struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"title"`
	DetailsID primitive.ObjectID `bson:"details_id"`
	Content   string             `bson:"content"`
}

func GetProblems(ctx context.Context, database db.Database, userId primitive.ObjectID) ([]db.Problem, error) {
	var problems []db.Problem

	cursor, err := database.Collection("problem_set").Find(ctx, bson.M{})
	if err != nil {
		return problems, err
	}
	err = cursor.All(ctx, &problems)

	return problems, err
}

func GetProblemBySlug(ctx context.Context, database db.Database, problemSlug string, userId primitive.ObjectID) ([]ProblemWithDetails, error) {
	pipeline := mongo.Pipeline{
		{
			// $match stage
			{
				Key: "$match",
				Value: bson.D{
					{Key: "slug", Value: problemSlug},
				},
			},
		},
		{
			// $lookup stage
			{
				Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "problem_details"},  // The name of the second collection
					{Key: "localField", Value: "details_id"}, // Field from the first collection
					{Key: "foreignField", Value: "_id"},      // Field from the second collection
					{Key: "as", Value: "details"},            // Alias for the joined data
				},
			},
		},
		{
			// $unwind stage
			{Key: "$unwind", Value: "$details"},
		},
		{
			// $project stage
			{
				Key: "$project",
				Value: bson.D{
					{Key: "details_id", Value: "$details._id"},  // Rename _id to orderId
					{Key: "content", Value: "$details.content"}, // Include customer name
					{Key: "title", Value: 1},
				},
			},
		},
	}

	var results []ProblemWithDetails

	cursor, err := database.Collection("problem_set").Aggregate(ctx, pipeline)
	if err != nil {
		return results, err
	}

	err = cursor.All(ctx, &results)

	return results, err
}
