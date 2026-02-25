package payment

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
		col: db.Collection("payment"),
	}
}

func (r *Repo) Create(ctx context.Context, p Payment) (Payment, error) {
	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	p.ID = bson.NewObjectID()

	_, err := r.col.InsertOne(childCtx, p)

	if err != nil {
		return Payment{}, fmt.Errorf("Unable to insert product: %v", err)
	}

	return p, nil
}

func (r *Repo) GetAll(ctx context.Context, userId string) ([]Payment, error) {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	filter := bson.M{
		"userId": userId,
	}

	var payments []Payment

	cursor, err := r.col.Find(childCtx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	defer cursor.Close(childCtx)

	if err := cursor.All(childCtx, &payments); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %v", err)
	}

	return payments, nil

}



func (r *Repo) GetById(ctx context.Context, userId string, paymentId string) (Payment, error) {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	id, errr := bson.ObjectIDFromHex(paymentId)

	if errr != nil {
		return Payment{}, fmt.Errorf("failed to get id: %v", errr)
	}

	var payment Payment

	filter := bson.M{
		"_id":    id,
		"userId": userId,
	}
	err := r.col.FindOne(childCtx, filter).Decode(&payment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Payment{}, mongo.ErrNoDocuments
		}
		return Payment{}, fmt.Errorf("failed to get order: %v", err)
	}

	return payment, nil

}

func (r *Repo) UpdateStatus(ctx context.Context, stripePaymentId string, status string) error {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	filter := bson.M{
		"stripePaymentId": stripePaymentId,
	}

    option := bson.M{"$set": bson.M{
            "status":    status,
            "updatedAt": time.Now().UTC(),
     },
    }

	_, errr := r.col.UpdateOne(childCtx, filter, option)
	

	if errr != nil {
		return fmt.Errorf("failed to update payment status: %v", errr)
	}

	return nil

}
