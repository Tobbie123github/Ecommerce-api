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
		return Payment{}, fmt.Errorf("failed to insert payment: %v", err)
	}

	return p, nil
}

func (r *Repo) GetAll(ctx context.Context, userId string) ([]Payment, error) {
	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{"userId": userId}

	var payments []Payment

	cursor, err := r.col.Find(childCtx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %v", err)
	}
	defer cursor.Close(childCtx)

	if err := cursor.All(childCtx, &payments); err != nil {
		return nil, fmt.Errorf("failed to decode payments: %v", err)
	}

	return payments, nil
}

func (r *Repo) GetById(ctx context.Context, userId string, paymentId string) (Payment, error) {
	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	id, err := bson.ObjectIDFromHex(paymentId)
	if err != nil {
		return Payment{}, fmt.Errorf("invalid payment id: %v", err)
	}

	var payment Payment

	filter := bson.M{
		"_id":    id,
		"userId": userId,
	}

	err = r.col.FindOne(childCtx, filter).Decode(&payment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Payment{}, mongo.ErrNoDocuments
		}
		return Payment{}, fmt.Errorf("failed to get payment: %v", err)
	}

	return payment, nil
}

func (r *Repo) UpdateStatus(ctx context.Context, reference string, status string) error {
	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{
		"paystackReference": reference,
	}

	update := bson.M{"$set": bson.M{
		"status":    status,
		"updatedAt": time.Now().UTC(),
	}}

	_, err := r.col.UpdateOne(childCtx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	return nil
}


func (r *Repo) Orders(ctx context.Context) ([]Payment, error) {

	childCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	filter := bson.M{}

	cursor, err := r.col.Find(childCtx, filter)

	if err != nil {
		return nil, err
	}

	var payments []Payment

	if err := cursor.All(childCtx, &payments); err != nil {
		return nil, err
	}

	return payments, nil

	
}