package dep

import (
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"gitlab.com/projectreferral/payment-api/configs"
	"gitlab.com/projectreferral/payment-api/external/dynamodb"
	"gitlab.com/projectreferral/payment-api/internal"
	"gitlab.com/projectreferral/payment-api/internal/rabbitmq"
	"gitlab.com/projectreferral/payment-api/internal/service"
	"gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/card"
	customer "gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/customer"
	sub "gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/subscription"
	token "gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/token"
	rabbit "gitlab.com/projectreferral/util/client/rabbitmq"
	util_dynamo "gitlab.com/projectreferral/util/pkg/dynamodb"
	"log"
)

//methods that are implemented on util
//and will be used
type ConfigBuilder interface{
	SetEnvConfigs()
	SetDynamoDBConfigsAndBuild() *util_dynamo.Wrapper
	SetRabbitMQConfigsAndBuild() *rabbit.DefaultQueueClient
}

func Inject(builder ConfigBuilder){

	builder.SetEnvConfigs()

	//setup dynamo library
	dynamoClient := builder.SetDynamoDBConfigsAndBuild()
	//connect to the instance
	log.Println("Connecting to Dynamo Client")
	dynamoClient.DefaultConnect()

	rabbitMQClient := builder.SetRabbitMQConfigsAndBuild()
	//dependency injection to our resource
	//we inject the rabbitmq client
	LoadRabbitMQClient(rabbitMQClient)

	stripe.Key = configs.StripeKey

	subscriptionServ := service.Subscription{
		CustomerClient: &customer.APIHelper{},
		SubClient:      &sub.APIHelper{},
		TokenClient:    &token.APIHelper{},
		CardClient:     &card.APIHelper{},
		SubscriptionRepo: 	&dynamodb.SubRepo{},

	}
	log.Println("Loading endpoints...")
	eb := internal.EndpointBuilder{}

	eb.SetupRouter(mux.NewRouter())
	eb.InjectSubscriptionServ(subscriptionServ)
	eb.SetupEndpoints()
	log.Println("All Dependencies injected")
}

func LoadRabbitMQClient(c rabbit.QueueClient){
	log.Println("Injecting RabbitMQ Client")
	rabbitmq.Client = c
}