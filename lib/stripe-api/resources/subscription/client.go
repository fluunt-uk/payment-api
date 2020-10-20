package card

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
	stripe_api "gitlab.com/projectreferral/payment-api/lib/stripe-api"
	"gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/models"
	"net/http"
)

//interface with the implemented methods will be injected in this variable
type Client interface {
	Put(c *stripe.Customer, pt string) (*models.Subscription, error)
	Get(http.ResponseWriter, *http.Request)
	Cancel(http.ResponseWriter, *http.Request)
	Patch(http.ResponseWriter, *http.Request)
	GetBatch(http.ResponseWriter, *http.Request)
}

type APIHelper struct{}

func (ah *APIHelper) Put(c *stripe.Customer, pt string) (*models.Subscription, error){
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(c.ID),

		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(pt),
			},
		},
	}
	s, e := sub.New(params)

	if e != nil {
		return nil, e
	}

	var sm = &models.Subscription{
		Email:          c.Email,
		AccountID:      s.Customer.ID,
		SubscriptionID: s.ID,
		PlanID:         s.Plan.ID,
		PlanType:       s.Plan.Nickname,
		Price:			s.Plan.Amount,
	}

	//TODO: might not be needed?
	//status, err := cw.DynamoSubRepo.Create(sm)
	//
	//if err != nil{
	//	fmt.Println(status, err)
	//}

	return sm, nil
}

func (ah *APIHelper) Get(w http.ResponseWriter, r *http.Request) {
	s, _ := sub.Get("sub_H6qCxUjOuCCmfj", nil)

	stripe_api.ReturnSuccessJSON(w, &s)
}

func (ah *APIHelper) Patch(w http.ResponseWriter, r *http.Request) {
	params := &stripe.SubscriptionParams{}
	params.AddMetadata("order_id", "0001")
	s, _ := sub.Update("sub_H6qCxUjOuCCmfj", params)

	stripe_api.ReturnSuccessJSON(w, &s)
}

func (ah *APIHelper) Cancel(w http.ResponseWriter, r *http.Request) {
	s, _ := sub.Cancel("sub_H6qCxUjOuCCmfj", nil)

	stripe_api.ReturnSuccessJSON(w, &s)
	//status, err = DeleteSubscription()
}

//it return 3 ReturnSuccessJSON as per the limit
//but SOMEHOW (to-be figured out) the method is auto called as many times as needed to get all Subs
func (ah *APIHelper) GetBatch(w http.ResponseWriter, r *http.Request) {
	params := &stripe.SubscriptionListParams{}
	//A limit on the number of objects to be returned. Limit can range between 1 and 100, and the default is 10.
	params.Filters.AddFilter("limit", "", "3")
	i := sub.List(params)
	for i.Next() {
		s := i.Subscription()
		stripe_api.ReturnSuccessJSON(w, &s)
	}
}
