package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repo struct {
	col *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		col: db.Collection("order"),
	}
}

// Create a new order

func (r *Repo) Create(ctx context.Context, o Order) (Order, error) {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	o.ID = bson.NewObjectID()

	_, err := r.col.InsertOne(childCtx, o)

	if err != nil {
		return Order{}, fmt.Errorf("Unable to insert product: %v", err)
	}

	return o, nil

}

func (r *Repo) GetById(ctx context.Context, orderId string) (Order, error) {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	id, errr := bson.ObjectIDFromHex(orderId)

	if errr != nil {
		return Order{}, fmt.Errorf("failed to get id: %v", errr)
	}

	var order Order
	err := r.col.FindOne(childCtx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Order{}, mongo.ErrNoDocuments
		}
		return Order{}, fmt.Errorf("failed to get order: %v", err)
	}

	return order, nil

}

func (r *Repo) GetByUserId(ctx context.Context, userId string) ([]Order, error) {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	filter := bson.M{"userId": userId}

	cursor, err := r.col.Find(childCtx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	defer cursor.Close(childCtx)

	var orders []Order
	if err := cursor.All(childCtx, &orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %v", err)
	}

	return orders, nil
}

func (r *Repo) Delete(ctx context.Context, orderId bson.ObjectID) error {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := r.col.DeleteOne(childCtx, bson.M{"_id": orderId})
	if err != nil {
		return fmt.Errorf("failed to delete order: %v", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("order not found")
	}

	return nil
}

func (r *Repo) UpdateStatus(ctx context.Context, orderId string, status string, paymentRef string) error {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	id, err := bson.ObjectIDFromHex(orderId)

	if err != nil {
		return fmt.Errorf("Unable to convert id, %v", err)
	}

	filter := bson.M{
		"_id": id,
	}

	option := bson.M{"$set": bson.M{
		"status":     status,
		"paymentRef": paymentRef,
		"updatedAt":  time.Now().UTC(),
	},
	}

	_, errr := r.col.UpdateOne(childCtx, filter, option)
	// _, err := r.col.UpdateOne(childCtx, filter, option)

	if errr != nil {
		return fmt.Errorf("failed to update payment status: %v", errr)
	}

	return nil

}
