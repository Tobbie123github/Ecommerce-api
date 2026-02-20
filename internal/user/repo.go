package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repo struct {
	col *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		col: db.Collection("users"),
	}
}


func (r *Repo) FindByEmail(ctx context.Context, email string) (User, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	filter := bson.M{
		"email": email,
	}

	var u User

	err := r.col.FindOne(ctx, filter).Decode(&u)
	if err != nil{
		if errors.Is(err, mongo.ErrNoDocuments){
			return User{}, mongo.ErrNoDocuments
		}

		return User{}, fmt.Errorf("Find by email failed: %s", err)
	}

	return u, nil
}



func (r *Repo) Create(ctx context.Context, u User) (User, error){
	childCtx, cancel := context.WithTimeout(ctx, 15*time.Second)

	defer cancel()

	u.ID = bson.NewObjectID()

	_, err := r.col.InsertOne(childCtx, u)

	if err != nil {
		return User{}, fmt.Errorf("Error inserting New user: %v", err)
	}

	return u, nil

}