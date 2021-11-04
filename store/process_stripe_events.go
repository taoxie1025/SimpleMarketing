package store

import (
	"email_action/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/stripe/stripe-go/v72"
	"strings"
)

func (s *Store) ProcessStripeEvents(reqBody []byte, header string) error {
	log.Infof("ProcessStripeEvents():")
	event, err := s.GetStripeAdapter().ConstructStripeEvent(reqBody, header)
	if err != nil {
		return err
	}
	switch event.Type {
	case "checkout.session.completed":
		log.Info("ProcessStripeEvents(): checkout.session.completed")
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Errorf("ProcessStripeEvents(): error = %v", err)
			return err
		}
		subscriptionId := session.Subscription.ID
		clientReferenceId := session.ClientReferenceID
		subscriptionPriceId := strings.Split(clientReferenceId, ",")[1]
		log.Infof("ProcessStripeEvents(): %v", clientReferenceId)
		targetPlan := models.SubscriptionPlanFree
		if subscriptionPriceId == beego.AppConfig.String("pro_plan_price_id") {
			targetPlan = models.SubscriptionPlanPro
		} else if subscriptionPriceId == beego.AppConfig.String("ultra_plan_price_id") {
			targetPlan = models.SubscriptionPlanUltra
		}
		return s.UpdateSubscriptionPlan(session.CustomerEmail, targetPlan, 0, subscriptionId,
			subscriptionPriceId, false)
	case "customer.subscription.updated":
		log.Info("ProcessStripeEvents(): customer.subscription.updated")
		var customer stripe.Customer
		err := json.Unmarshal(event.Data.Raw, &customer)
		if err != nil {
			log.Errorf("ProcessStripeEvents(): error = %v", err)
			return err
		}
		invoiceStatus := customer.Subscriptions.Data[0].LatestInvoice.Status
		if invoiceStatus == "past_due" {
			return s.UpdateAccountStatus(customer.Email, true, models.PaymentStatusFailed)
		}
	default:
		log.Infof("ProcessStripeEvents(): Unhandled event type: %s", event.Type)
	}

	return nil
}

func (s *Store) UpdateSubscriptionPlan(email string, targetPlan int, emailUsageInCycle int64, subscriptionId,
	subscriptionPriceId string, isBlock bool) error {
	log.Infof("UpdateSubscriptionPlan(): email = %s, targetPlan = %d, emailUsageInCycle = %d, "+
		"subscriptionId = %s, subscriptionPriceId = %s, isBlock = %s", email, targetPlan, emailUsageInCycle,
		subscriptionId, subscriptionPriceId, isBlock)
	return s.GetDynamodbAdapter().UpdateUserSubscriptionPlan(email, targetPlan, models.PaymentStatusCompleted,
		emailUsageInCycle, subscriptionId, subscriptionPriceId, isBlock)
}

func (s *Store) UpdateAccountStatus(email string, isBlock bool, paymentStatus int) error {
	log.Infof("UpdateAccountStatus(): email = %s, isBlock = %s, paymentStatus = %d", email, isBlock, paymentStatus)
	return s.GetDynamodbAdapter().UpdateAccountStatus(email, isBlock, paymentStatus)
}
