package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/3ssalunke/leetcode-clone/db"
	"github.com/3ssalunke/leetcode-clone/token"
	"github.com/3ssalunke/leetcode-clone/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type contextUserKey string

const ContextUserKey contextUserKey = "LoggedInUser"

func AuthMiddleware(tokenMaker *token.TokenMaker, database db.Database) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			jwtToken, err := util.ExtractCookieFromHeader(r)
			if err != nil {
				log.Printf("failed to extract cookie from request header - %v", err)
			} else {
				payload, tokenerr := tokenMaker.VerifyToken(jwtToken)
				if tokenerr != nil {
					log.Printf("failed to verify token - %v", tokenerr)
				} else {
					var user db.User
					filter := bson.M{
						"$or": []bson.M{
							{"email": payload.Email},
							{"username": payload.Username},
						},
					}

					dberr := database.Collection("users").FindOne(ctx, filter).Decode(&user)
					if dberr != nil {
						log.Printf("failed to fetch user details - %v", dberr)
					}
					if user.ID != primitive.NilObjectID {
						r = r.WithContext(context.WithValue(r.Context(), ContextUserKey, user))
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
