package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-auth/internal/order"
	"net/http"
	"os"
	"strings"
	"time"
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

type paystackInitResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type paystackVerifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status    string `json:"status"`
		Reference string `json:"reference"`
		Amount    int64  `json:"amount"`
		Metadata  struct {
			OrderID string `json:"orderId"`
			UserID  string `json:"userId"`
		} `json:"metadata"`
	} `json:"data"`
}

func (svc *Service) CreatePayment(ctx context.Context, userId string, orderId string, email string) (string, string, error) {
	secretKey := os.Getenv("PAYSTACK_SECRET_KEY")

	// Get order
	order, err := svc.orderRepo.GetById(ctx, orderId)
	if err != nil {
		return "", "", fmt.Errorf("unable to get order: %v", err)
	}

	if order.UserID != userId {
		return "", "", fmt.Errorf("unauthorized user")
	}

	// Amount in kobo (Paystack uses smallest currency unit)
	amountKobo := int64(order.Total * 100)

	payload := fmt.Sprintf(`{
    "email": "%s",
    "amount": %d,
    "currency": "NGN",
    "callback_url": "%s",
    "metadata": {
        "orderId": "%s",
        "userId": "%s"
    }
}`, email, amountKobo, os.Getenv("PAYSTACK_CALLBACK_URL"), orderId, userId)

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.paystack.co/transaction/initialize",
		strings.NewReader(payload),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to call paystack: %v", err)
	}
	defer resp.Body.Close()

	var psResp paystackInitResponse
	if err := json.NewDecoder(resp.Body).Decode(&psResp); err != nil {
		return "", "", fmt.Errorf("failed to decode paystack response: %v", err)
	}

	if !psResp.Status {
		return "", "", fmt.Errorf("paystack error: %s", psResp.Message)
	}

	// Save payment record
	_, err = svc.repo.Create(ctx, Payment{
		OrderID:           orderId,
		UserID:            userId,
		Amount:            order.Total,
		Currency:          "NGN",
		Status:            "pending",
		PaystackReference: psResp.Data.Reference,
		AuthorizationURL:  psResp.Data.AuthorizationURL,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to save payment: %v", err)
	}

	return psResp.Data.AuthorizationURL, psResp.Data.Reference, nil
}

func (svc *Service) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	secretKey := os.Getenv("PAYSTACK_SECRET_KEY")

	// Verify Paystack webhook signature
	mac := hmac.New(sha512.New, []byte(secretKey))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return fmt.Errorf("invalid webhook signature")
	}

	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to parse webhook: %v", err)
	}

	eventType, _ := event["event"].(string)

	switch eventType {
	case "charge.success":
		data, _ := event["data"].(map[string]interface{})
		reference, _ := data["reference"].(string)
		metadataRaw, _ := data["metadata"].(map[string]interface{})
		orderId, _ := metadataRaw["orderId"].(string)
		userId, _ := metadataRaw["userId"].(string)

		if err := svc.orderService.UpdateOrder(ctx, orderId, userId, "paid", reference); err != nil {
			return fmt.Errorf("failed to update order: %v", err)
		}

		if err := svc.repo.UpdateStatus(ctx, reference, "paid"); err != nil {
			return fmt.Errorf("failed to update payment: %v", err)
		}
	}

	return nil
}

// VerifyPayment — called after redirect to confirm payment
func (svc *Service) VerifyPayment(ctx context.Context, reference string) (*paystackVerifyResponse, error) {
	secretKey := os.Getenv("PAYSTACK_SECRET_KEY")

	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.paystack.co/transaction/verify/%s", reference),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+secretKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify payment: %v", err)
	}
	defer resp.Body.Close()

	var psResp paystackVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&psResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &psResp, nil
}
