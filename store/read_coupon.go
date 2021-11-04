package store

import (
	"github.com/stripe/stripe-go/v72"
)

func (s *Store) ReadCoupon(couponId string) (*stripe.Coupon, error) {
	log.Infof("ReadCoupon(%s)", couponId)
	return s.stripeAdapter.ReadCoupon(couponId)
}
