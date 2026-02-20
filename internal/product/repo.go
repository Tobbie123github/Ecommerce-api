package product

import (
	"context"
	"fmt"
	"time"

	
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repo struct {
	col *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		col: db.Collection("products"),

	}
}

func (r *Repo) Add(ctx context.Context, p Product) (Product, error) {
	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	p.ID = bson.NewObjectID()

	_, err := r.col.InsertOne(childCtx, p)

	if err != nil {
		return Product{}, fmt.Errorf("Unable to insert product: %v", err)
	}

	return p, nil 

}

func (r *Repo) GetAll(ctx context.Context) ([]Product, error) {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()


	filter := bson.M{}
	cursor , err := r.col.Find(childCtx, filter)

	if err != nil {
		return nil, fmt.Errorf("Failed to get all product: %v", err)
	}

	defer cursor.Close(childCtx)

	var result []Product

	if err := cursor.All(childCtx, &result); err != nil {
		return nil, fmt.Errorf("Failed to get Products")
	}

	return result, nil

}


func (r *Repo) GetId(ctx context.Context, id bson.ObjectID) (Product, error) {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	var product Product


	err := r.col.FindOne(childCtx, filter, options.FindOne()).Decode(&product)

	if err != nil {
		return Product{}, fmt.Errorf("Error getting product: %v", err)
	}

	return product, nil

}
