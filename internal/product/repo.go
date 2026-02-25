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


func (r *Repo) GetById(ctx context.Context, id bson.ObjectID) (Product, error) {

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


func (r *Repo) UpdateById(ctx context.Context, id bson.ObjectID, req UpdateProduct) (Product, error) {
	 
		childCtx, cancel := context.WithTimeout(ctx, 15*time.Second) 

		defer cancel()

		filter := bson.M{
			"_id":id,
		}


		update := bson.M{
			"$set":bson.M{
				"name": req.Name,
				"description": req.Description,
				"price":req.Price,
				"discount":req.Discount,
				"category" : req.Category,
				"instock": req.InStock,
				"imageUrl": req.ImageURL,
				"publicIds":req.PublicIDs,
				"updatedAt": time.Now().UTC(),
				
			},
		}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		var updated Product

		err := r.col.FindOneAndUpdate(childCtx, filter, update, opts).Decode(&updated)

		if err !=  nil {
			return Product{}, fmt.Errorf("error inserting update: %v", err)
		}

	return updated, nil

}


func (r *Repo) DeleteItem(ctx context.Context, id bson.ObjectID) error {
	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	filter := bson.M{
		"_id": id,
	} 

	_, err := r.col.DeleteOne(childCtx, filter)

	if err != nil {
		return fmt.Errorf("error with deleting: %v", err)
	}

	return nil
}
