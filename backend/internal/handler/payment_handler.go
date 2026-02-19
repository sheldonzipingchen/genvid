package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/genvid/backend/internal/config"
	"github.com/genvid/backend/internal/model"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/checkout/session"
)

type PaymentHandler struct {
	cfg *config.Config
}

func NewPaymentHandler(cfg *config.Config) *PaymentHandler {
	stripe.Key = cfg.External.Stripe.SecretKey
	return &PaymentHandler{cfg: cfg}
}

type CheckoutRequest struct {
	PlanID     string `json:"plan_id"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}

type CheckoutResponse struct {
	SessionURL string `json:"session_url"`
	SessionID  string `json:"session_id"`
}

var planConfigs = map[string]struct {
	Name   string
	Amount int64
	Credits int
}{
	"starter_monthly":  {"Starter Monthly", 1900, 15},
	"starter_yearly":   {"Starter Yearly", 19000, 180},
	"pro_monthly":      {"Pro Monthly", 4900, 50},
	"pro_yearly":       {"Pro Yearly", 49000, 600},
	"business_monthly": {"Business Monthly", 9900, 150},
	"business_yearly":  {"Business Yearly", 99000, 1800},
}

func (h *PaymentHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	planConfig, ok := planConfigs[req.PlanID]
	if !ok {
		respondError(w, http.StatusBadRequest, "INVALID_PLAN", "Invalid plan ID", nil)
		return
	}

	successURL := req.SuccessURL
	if successURL == "" {
		successURL = fmt.Sprintf("%s/dashboard?payment=success", h.cfg.Server.AppURL)
	}

	cancelURL := req.CancelURL
	if cancelURL == "" {
		cancelURL = fmt.Sprintf("%s/pricing?payment=canceled", h.cfg.Server.AppURL)
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("usd"),
					UnitAmount: stripe.Int64(planConfig.Amount),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(planConfig.Name),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Metadata: map[string]string{
			"user_id":  userID,
			"plan_id":  req.PlanID,
			"credits":  fmt.Sprintf("%d", planConfig.Credits),
		},
	}

	s, err := session.New(params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "STRIPE_ERROR", "Failed to create checkout session", nil)
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse(CheckoutResponse{
		SessionURL: s.URL,
		SessionID:  s.ID,
	}))
}

func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	payload := make([]byte, r.ContentLength)
	r.Body.Read(payload)
	
	eventType := r.Header.Get("Stripe-Event")
	
	if eventType == "checkout.session.completed" {
		var event struct {
			Data struct {
				Object struct {
					Metadata map[string]string `json:"metadata"`
				} `json:"object"`
			} `json:"data"`
		}
		
		if err := json.Unmarshal(payload, &event); err == nil {
			userID := event.Data.Object.Metadata["user_id"]
			_ = userID
		}
	}

	w.WriteHeader(http.StatusOK)
}
