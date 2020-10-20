package dynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gitlab.com/projectreferral/payment-api/lib/stripe-api/resources/models"
	"gitlab.com/projectreferral/util/pkg/dynamodb"
)

type SubRepo struct {
	//dynamo client
	DC		*dynamodb.Wrapper
}
//implement only the necessary methods for each repository
//available to be consumed by the API
type Repo interface{
	Create(*models.Subscription) (string, error)
	Del(string) (string, error)
}

//get all the adverts for a specific account
//token validated
func (s *SubRepo) Create(body *models.Subscription) (string, error) {

	av, errM := dynamodbattribute.MarshalMap(body)

	if errM != nil {
		return "", errM
	}

	fmt.Println(av)
	err := s.DC.CreateItem(av)

	if err != nil {
		// Need to handle changing premium status here, will need to call endpoint
		return "Failed", err
	}
	return "Success", err
}

func (s *SubRepo) Del(email string) (string, error) {
	err := s.DC.DeleteItem(email)

	if err != nil {
		return "Failed to delete item", err
	}
	return "Item deleted", nil
}

