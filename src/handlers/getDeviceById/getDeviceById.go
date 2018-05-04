package main

import (
    "eloy-aws-api-service/src/handlers/types"

    "os"
    "encoding/json"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var ERROR_MISSING_ID_FIELD = 1
var ERROR_INTERNAL_SERVERS_DATABAE = 2
var ERROR_NO_ITEM_FOUNDED = 3


type SuccessResponse struct {
   Device   types.Device    `json:"data"`
}


type dynamoDBAPI struct {
    DynamoDB dynamodbiface.DynamoDBAPI
}


// get a device from DynamoDB database with provided id
func (ig *dynamoDBAPI) getFromDatabase(id string) ( *dynamodb.GetItemOutput, error) {
    // Get table name from OS's environment
    tableName := aws.String(os.Getenv("DEVICES_TABLE_NAME"))

    var input = &dynamodb.GetItemInput{
        TableName: tableName,
        
        Key: map[string]*dynamodb.AttributeValue{
            "id": {
                S: aws.String(id),
            },
        },
    }

     result, err := ig.DynamoDB.GetItem(input)

     return result, err
}


// main AWS lambda function starting point.
// It gets an id from client, parse it and tries to get corresponding device fromdynamodb.
func GetDeviceById(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        
    region := os.Getenv("AWS_REGION")
    var getter = new(dynamoDBAPI)
    sess, err := session.NewSession(&aws.Config{Region: &region},)

    svc := dynamodb.New(sess)
    getter.DynamoDB = dynamodbiface.DynamoDBAPI(svc)



    // get requested id from APIGatewayProxyRequest 
    id := request.PathParameters["id"]

    // If no id provided, return HTTP error 404
    if id == "" {
        return events.APIGatewayProxyResponse{
            Body:       createErrorResponseJson(ERROR_MISSING_ID_FIELD),
            StatusCode: 404,
        }, nil
    }

    result, err := getter.getFromDatabase(id)

    validationResult := validateDatabaseResult(result, err)

    return validationResult , nil
}


func validateDatabaseResult(result *dynamodb.GetItemOutput, err error)( events.APIGatewayProxyResponse) {

    // If an internal error occured in the database, return HTTP error 500
    // todo: log here
    if err != nil {
        return events.APIGatewayProxyResponse{ 
            Body:       createErrorResponseJson(ERROR_INTERNAL_SERVERS_DATABAE),
            StatusCode: 500,
        }
    }

    
    // If no item founded, return error 404
    if len(result.Item) == 0 {
        return events.APIGatewayProxyResponse{
            Body:       createErrorResponseJson(ERROR_NO_ITEM_FOUNDED),
            StatusCode: 404,
        }
    }


    // returned founded item as json file with 200 HTTP status code.
    // foundedDeviceAsJson, _ := json.Marshal(item)
    return events.APIGatewayProxyResponse{ 
        Body: createSuccessResponseJson(result),
        StatusCode: 200,
    }
}


func createErrorResponseJson(errorState int) (jsonString string) {


    if errorState == ERROR_MISSING_ID_FIELD {
            errorResponse := types.ErrorResponse { ErrorMessage: types.ErrorMessage { Code: 404,Message: "No ID Field Provided",},}
            errorResponseJson, _ := json.MarshalIndent(&errorResponse, "", "\t")
            return string(errorResponseJson)

    } else if errorState == ERROR_INTERNAL_SERVERS_DATABAE {
            errorResponse := types.ErrorResponse { ErrorMessage: types.ErrorMessage { Code: 500,Message: "Internal Server's Error occured",},}
            errorResponseJson, _ := json.MarshalIndent(&errorResponse, "", "\t")
            return string(errorResponseJson)

    } else {
            errorResponse := types.ErrorResponse { ErrorMessage: types.ErrorMessage { Code: 404,Message: "Desired device with provided id was not founded",},}
            errorResponseJson, _ := json.MarshalIndent(&errorResponse, "", "\t")
            return string(errorResponseJson)
    }
}


func createSuccessResponseJson(result *dynamodb.GetItemOutput) (jsonString string) {

    // create json file of database's returned Item 
    item := types.Device{}
    dynamodbattribute.UnmarshalMap(result.Item, &item)

    successResponse := SuccessResponse { 
        item,
    }
    successResponseJson, _ := json.MarshalIndent(&successResponse, "", "\t")

    return string(successResponseJson)
}

func main() {

    lambda.Start(GetDeviceById)
}