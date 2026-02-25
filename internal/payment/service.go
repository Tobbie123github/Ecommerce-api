package payment

import (
	"context"
	"encoding/json"

	//"errors"
	"fmt"
	"go-auth/internal/config"
	"go-auth/internal/order"
	"os"
	"time"

	//"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/v81"

	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/webhook"
)

type Service struct {
	repo         *Repo
	orderRepo    *order.Repo
	orderService *order.Service
}

func NewService(repo *Repo, orderRepo *order.Repo, orderService *order.Service) *Service {
	return &Service{
		repo:         repo,
		orderRepo:    orderRepo,
		orderService: orderService,
	}
}

func (svc Service) CreatePayment(ctx context.Context, userId string, orderId string) (string, error) {

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Get Order

	order, err := svc.orderRepo.GetById(ctx, orderId)

	if err != nil {
		return "", fmt.Errorf("Unable to get order: %v", err)
	}

	// check if order is for user

	if order.UserID != userId {
		return "", fmt.Errorf("Unauthorized User: %v", err)
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(order.Total * 100)),
		Currency: stripe.String("usd"),
		Metadata: map[string]string{
			"orderId": orderId,
			"userId":  userId,
		},
	}

	params.SetIdempotencyKey(orderId)

	pi, err := paymentintent.New(params)

	if err != nil {
		return "", fmt.Errorf("Error creating PaymentIntent: %v", err)
	}

	_, err = svc.repo.Create(ctx, Payment{
		OrderID:         orderId,
		UserID:          userId,
		Amount:          order.Total,
		Currency:        "usd",
		Status:          "pending",
		StripePaymentID: pi.ID,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	})

	if err != nil {
		return "", fmt.Errorf("failed to save payment: %v", err)
	}

	return pi.ClientSecret, nil

}

func (svc *Service) HandleWebhook(ctx context.Context, payload []byte, signature string) error {

	// when webhook is successfull/ payment is successful, we want to
	// update Order status from pending to paid
	// and also update payment from pending to paid

	cfg, err := config.Load()

	event, err := webhook.ConstructEventWithOptions(payload, signature, cfg.STRIPE_WEBHOOK_SECRET, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		return fmt.Errorf("Invalid signature %v", err)
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent

		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return fmt.Errorf("failed to parse payment: %v", err)
		}

		orderId := pi.Metadata["orderId"]
		userId := pi.Metadata["userId"]

		if err := svc.orderService.UpdateOrder(ctx, orderId, userId, "paid", pi.ID); err != nil {
			return fmt.Errorf("failed to update order: %v", err)
		}

		if err := svc.repo.UpdateStatus(ctx, pi.ID, "paid"); err != nil {
			return fmt.Errorf("failed to update payment: %v", err)
		}

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent

		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return fmt.Errorf("failed to parse payment: %v", err)
		}

		orderId := pi.Metadata["orderId"]
		userId := pi.Metadata["userId"]

		paymentRef := pi.ID

		if err := svc.orderService.UpdateOrder(ctx, orderId, userId, "failed", paymentRef); err != nil {
			return fmt.Errorf("failed to update order: %v", err)
		}

		if err := svc.repo.UpdateStatus(ctx, pi.ID, "failed"); err != nil {
			return fmt.Errorf("failed to update payment: %v", err)
		}

	}

	return nil

}
