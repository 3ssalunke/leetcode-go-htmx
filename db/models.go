package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id"`
	Username   string             `bson:"username"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	ImageUrl   string             `bson:"image_url"`
	Created_at primitive.DateTime `bson:"created_at"`
	Updated_at primitive.DateTime `bson:"updated_at"`
}

type Problem struct {
	ID         primitive.ObjectID `bson:"_id"`
	Title      string             `bson:"title"`
	Slug       string             `bson:"slug"`
	Acceptance float64            `bson:"acceptance"`
	Difficulty string             `bson:"difficulty"`
	Created_at primitive.DateTime `bson:"created_at"`
	Updated_at primitive.DateTime `bson:"updated_at"`
}
