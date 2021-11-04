package store

import (
	"github.com/stripe/stripe-go/v72"
)

func (s *Store) CreateCheckoutSession(stripeCustomerId, priceId, plan, email, coupon, successURL, cancelURL string) (*stripe.CheckoutSession, error) {
	log.Infof("CreateCheckoutSession(): stripeCustomerId = %s, priceId = %s, email = %s", stripeCustomerId, priceId, email)
	sess, err := s.GetStripeAdapter().CreateSession(stripeCustomerId, priceId, plan, coupon, email, successURL, cancelURL)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
