package stripe

import (
	"email_action/logging"
	"github.com/astaxie/beego"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/coupon"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/webhook"
	"strconv"
	"time"
)

var (
	log = logging.NewZapLogger()
)

type StripeAdapter struct {
	apiKey        string
	webhookSecret string
}

func NewStripeAdapter() *StripeAdapter {
	apiKey := beego.AppConfig.String("stripe_secret_key")
	stripe.Key = apiKey
	stripeAdapter := &StripeAdapter{
		apiKey: apiKey,
	}
	stripeAdapter.CreateWebhook()

	return stripeAdapter
}

func (d *StripeAdapter) CreateCustomer(email, firstName, lastName string) (*stripe.Customer, error) {
	log.Infof("CreateCustomer(): email = %s", email)
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(firstName + " " + lastName),
	}
	customer, err := customer.New(params)
	if err != nil {
		log.Errorf("CreateCustomer(): error = %v", err)
	}
	return customer, err
}

func (d *StripeAdapter) CreateSession(customerId, priceId, plan, email, coupon, successURL, cancelURL string) (*stripe.CheckoutSession, error) {
	log.Infof("CreateSession(): customerId = %s, priceId %s = , email = ", customerId, priceId, email)
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Items: []*stripe.CheckoutSessionSubscriptionDataItemsParams{
				{
					Plan:     stripe.String(priceId),
					Quantity: stripe.Int64(1),
				},
			},
			Metadata: map[string]string{
				"plan":    plan,
				"priceId": priceId,
				"email":   email,
			},
		},
		ClientReferenceID: stripe.String(email + "," + priceId + "," + plan + "," + strconv.FormatInt(time.Now().Unix(), 10)),
		SuccessURL:        stripe.String(successURL),
		CancelURL:         stripe.String(cancelURL),
	}
	params.CustomerEmail = stripe.String(email)

	if coupon != "" {
		params.SubscriptionData.Coupon = stripe.String(coupon)
	}

	checkoutSession, err := session.New(params)
	if err != nil {
		log.Errorf("CreateSession(): error = %v", err)
	}
	return checkoutSession, err
}

func (d *StripeAdapter) CreateWebhook() {
	webhookSecret := beego.AppConfig.String("stripe_webhook_secret")
	d.webhookSecret = webhookSecret

	/* create webhook dynamically
	log.Infof("CreateWebhook(): %s%s", httpaddr, "/api/v1/webhook")
	params := &stripe.WebhookEndpointParams{
		URL:           stripe.String(httpaddr + "/api/v1/webhook"),
		EnabledEvents: stripe.StringSlice([]string{"*"}),
		Description:   stripe.String("For user accounts' subscription plan only."),
	}

	endpoint, err := webhookendpoint.New(params)
	if err != nil {
		log.Errorf("CreateWebhook(): error = %v", err)
	}
	d.webhookSecret = endpoint.Secret
	*/
}

func (d *StripeAdapter) ConstructStripeEvent(reqBody []byte, header string) (*stripe.Event, error) {
	log.Infof("ConstructStripeEvent():")
	event, err := webhook.ConstructEvent(reqBody, header, d.webhookSecret)
	if err != nil {
		log.Errorf("ConstructStripeEvent(): error = %v", err)
		return nil, err
	}
	return &event, nil
}

func (d *StripeAdapter) CancelSubscription(subscriptionId string) error {
	log.Infof("CancelSubscription(): subscriptionId = %s", subscriptionId)
	_, err := sub.Cancel(subscriptionId, nil)
	if err != nil {
		log.Errorf("CancelSubscription(): error = %v", err)
	}
	return err
}

func (d *StripeAdapter) UpdateSubscription(subscriptionId, priceId string) error {
	log.Infof("UpdateSubscription(): subscriptionId = %s, priceId = %s", subscriptionId, priceId)
	existingSubscription, err := sub.Get(subscriptionId, nil)
	param := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(existingSubscription.Items.Data[0].ID),
				Price: stripe.String(priceId),
			},
		},
		ProrationBehavior: stripe.String(string(stripe.SubscriptionProrationBehaviorCreateProrations)),
		CancelAtPeriodEnd: stripe.Bool(false),
	}
	_, err = sub.Update(subscriptionId, param)
	if err != nil {
		log.Errorf("UpdateSubscription(): error = %v", err)
	}
	return err
}

func (d *StripeAdapter) ReadCoupon(couponId string) (*stripe.Coupon, error) {
	log.Infof("ReadCoupon(%s)", couponId)
	coupon, err := coupon.Get(couponId, nil)
	if err != nil {
		log.Warnf("ReadCoupon(): error = %v", err)
	}
	return coupon, err
}
