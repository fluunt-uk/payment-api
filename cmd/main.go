package main

import (
	"fmt"
	"gitlab.com/projectreferral/payment-api/cmd/dep"
	"gitlab.com/projectreferral/payment-api/configs"
	"gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/models"
	"gitlab.com/projectreferral/util/util"
	"log"
	"os"
)

func main() {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	f, err := os.OpenFile(path + "/logs/paymentAPI_log.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	//gets all the necessary configs into our object
	//completes connections
	//assigns connections to repos
	dep.Inject(&util.ServiceConfigs{
		Environment: os.Getenv("ENV"),
		Region:       configs.EU_WEST_2,
		Table:        configs.TABLE,
		SearchParam:  configs.UNIQUE_IDENTIFIER,
		GenericModel: models.Subscription{},
		BrokerUrl:    configs.QAPI_URL,
		Port:		  configs.PORT,
	})
}
