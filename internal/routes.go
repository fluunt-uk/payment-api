package internal

import (
	"fmt"
	"gitlab.com/projectreferral/payment-api/configs"
	"gitlab.com/projectreferral/payment-api/internal/service"
	"gitlab.com/projectreferral/util/pkg/security"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)


type EndpointBuilder struct {
	router       	*mux.Router
	ss 				service.Subscription
}

func (eb *EndpointBuilder) SetupRouter(route *mux.Router) {
	eb.router = route
}

func (eb *EndpointBuilder) InjectSubscriptionServ(ss service.Subscription) {
	eb.ss = ss
}

func (eb *EndpointBuilder) SetupEndpoints() {

	eb.router.HandleFunc("/premium/subscribe", security.WrapHandlerWithSpecialAuth(eb.ss.SubscribeToPremiumPlan, configs.AUTH_AUTHENTICATED)).Methods("POST")
	eb.router.HandleFunc("/log", displayLog).Methods("GET")

	log.Fatal(http.ListenAndServe(configs.PORT, eb.router))
}

func displayLog(w http.ResponseWriter, r *http.Request){
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	b, _ := ioutil.ReadFile(path + "/logs/paymentAPI_log.txt")

	w.Write(b)
}
