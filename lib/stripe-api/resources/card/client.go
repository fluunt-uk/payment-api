package card

import (
	stripe_api "gitlab.com/projectreferral/payment-api/lib/stripe-api"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/card"
	"net/http"
)

type Client interface {
	Put(c *stripe.Customer, t *stripe.Token) (*stripe.Card, error)
	Get(http.ResponseWriter, *http.Request)
	Del(http.ResponseWriter, *http.Request)
	Patch(http.ResponseWriter, *http.Request)
	GetBatch(http.ResponseWriter, *http.Request)
}

type APIHelper struct{}

func (ah *APIHelper) Put(c *stripe.Customer, t *stripe.Token) (*stripe.Card, error) {
	params := &stripe.CardParams{
		Customer: stripe.String(c.ID),
		Token: stripe.String(t.ID),
	}
	card, err := card.New(params)

	if err != nil {
		return nil, err
	}

	return card, nil
}

func (ah *APIHelper) Get(w http.ResponseWriter, r *http.Request)  {
	params := &stripe.CardParams{
		Customer: stripe.String("cus_H7EDAZGzu81IFr"),
	}
	c, _ := card.Get(
		"card_1GYzmlGhy1brUyYItScV9Lwo",
		params,
	)

	stripe_api.ReturnSuccessJSON(w, &c)
}

func (ah *APIHelper) Patch(w http.ResponseWriter, r *http.Request)  {
	params := &stripe.CardParams{
		Customer: stripe.String("cus_H7Dt44weDWU4s5"),
		Name: stripe.String("Jenny Rosen"),
	}
	c, _ := card.Update(
		"card_1GYzmlGhy1brUyYItScV9Lwo",
		params,
	)

	stripe_api.ReturnSuccessJSON(w, &c)
}

func (ah *APIHelper) Del(w http.ResponseWriter, r *http.Request)  {
	params := &stripe.CardParams{
		Customer: stripe.String("cus_H7Dt44weDWU4s5"),
	}
	c, _ := card.Del(
		"card_1GYzmlGhy1brUyYItScV9Lwo",
		params,
	)

	stripe_api.ReturnSuccessJSON(w, &c)
}

func (ah *APIHelper) GetBatch(w http.ResponseWriter, r *http.Request)  {
	params := &stripe.CardListParams{
		Customer: stripe.String("cus_H7Dt44weDWU4s5"),
	}
	params.Filters.AddFilter("limit", "", "3")
	i := card.List(params)
	for i.Next() {
		c := i.Card()

		stripe_api.ReturnSuccessJSON(w, &c)
	}
}