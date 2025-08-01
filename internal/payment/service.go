package payment

import (
	"context"
	"fmt"

	"github.com/purushothdl/ecommerce-api/configs"
	"github.com/purushothdl/ecommerce-api/internal/domain"
	"github.com/purushothdl/ecommerce-api/internal/shared/dto"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/refund"
)

// stripeService doesn't need to hold a client instance.
// The stripe-go library manages this globally.
type stripeService struct{}

// NewStripeService creates a new service for interacting with Stripe.
func NewStripeService(cfg configs.StripeConfig) domain.PaymentService {
	stripe.Key = cfg.SecretKey
	return &stripeService{}
}

// CreatePaymentIntent creates a payment intent on Stripe.
func (s *stripeService) CreatePaymentIntent(ctx context.Context, amount float64) (*dto.PaymentIntent, error) {
	// Stripe expects the amount in the smallest currency unit (e.g., cents).
	amountInPaise := int64(amount * 100)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInPaise),
		Currency: stripe.String(string(stripe.CurrencyINR)), 
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create stripe payment intent: %w", err)
	}

	// Map the Stripe response to our internal model.
	return &dto.PaymentIntent{
		ID:           pi.ID,
		ClientSecret: pi.ClientSecret,
		Amount:       amount,
		Currency:     string(pi.Currency),
		Status:       string(pi.Status),
	}, nil
}

// RefundPaymentIntent issues a full refund for a given Payment Intent.
func (s *stripeService) RefundPaymentIntent(ctx context.Context, paymentIntentID string) error {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
	}

	_, err := refund.New(params)
	if err != nil {
		// Check for specific Stripe errors if needed
		return fmt.Errorf("failed to create stripe refund: %w", err)
	}

	return nil
}