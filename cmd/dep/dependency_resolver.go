package dep

import (
	"gitlab.com/projectreferral/payment-api/configs"
	"gitlab.com/projectreferral/payment-api/internal"
	"gitlab.com/projectreferral/payment-api/internal/service"
	"gitlab.com/projectreferral/payment-api/lib/dynamodb/repo"
	"gitlab.com/projectreferral/payment-api/lib/rabbitmq"
	"gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/card"
	customer "gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/customer"
	sub "gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/subscription"
	token "gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/token"
	"gitlab.com/projectreferral/util/pkg/dynamodb"
	"gitlab.com/projectreferral/queueing-api/client"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"log"
)

//methods that are implemented on util
//and will be used
type ConfigBuilder interface{
	LoadEnvConfigs()
	LoadDynamoDBConfigs() *dynamodb.Wrapper
	LoadRabbitMQConfigs() *client.DefaultQueueClient
}

func Inject(builder ConfigBuilder){

	builder.LoadEnvConfigs()

	//setup dynamo library
	dynamoClient := builder.LoadDynamoDBConfigs()
	//connect to the instance
	log.Println("Connecting to Dynamo Client")
	dynamoClient.DefaultConnect()

	rabbitMQClient := builder.LoadRabbitMQConfigs()
	//dependency injection to our resource
	//we inject the rabbitmq client
	LoadRabbitMQClient(rabbitMQClient)

	stripe.Key = configs.StripeKey

	subscriptionServ := service.Subscription{
		CustomerClient: &customer.Wrapper{},
		SubClient:      &sub.Wrapper{DynamoSubRepo: &repo.Wrapper{DC: dynamoClient}},
		TokenClient:    &token.Wrapper{},
		CardClient:     &card.Wrapper{},
	}
	log.Println("Loading endpoints...")
	eb := internal.EndpointBuilder{}

	eb.SetupRouter(mux.NewRouter())
	eb.InjectSubscriptionServ(subscriptionServ)
	eb.SetupEndpoints()
	log.Println("All Dependencies injected")
}

func LoadRabbitMQClient(c client.QueueClient){
	log.Println("Injecting RabbitMQ Client")
	rabbitmq.Client = c
}