package main

import (
    "eloy-aws-api-service/src/handlers/types"
    "testing"
    "errors"
    
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// TestCase struct that contains all reuested and expected values for unit testing
type TestCase struct {
    Name                        string
    InputId                     events.APIGatewayProxyRequest
    InputIdString               string
    DatabaseOutput              dynamodb.GetItemOutput
    Error                       error
    ExpectedBody                string
    ExpectedStatusCode          int
    ExpectedDatabaseOutput      dynamodb.GetItemOutput
}



// A fakeDynamoDB instance for mocking test that emulates real DynamoDB
type FakeDynamoDBAPI struct {
    dynamodbiface.DynamoDBAPI
}


// a mocked version of DynamoDB's GetItem function.
// in testing state, instead of calling real DynamoDB's GetItem, we try to emulate it.
// Get function of getDeviceById.go calls this function in Testing state.  
func (fd *FakeDynamoDBAPI) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
    output := new(dynamodb.GetItemOutput)
    id := input.Key["id"].S
    
    if *id == "id_test" {
        output.SetItem(
            map[string]*dynamodb.AttributeValue{
                "id": &dynamodb.AttributeValue{S: aws.String("id_test")},
                "deviceModel": &dynamodb.AttributeValue{S: aws.String("deviceModel_test")},
                "name": &dynamodb.AttributeValue{S: aws.String("name_test")},
                "note": &dynamodb.AttributeValue{S: aws.String("note_test")},
                "serial": &dynamodb.AttributeValue{S: aws.String("serial_test")},
            },
        )
    }    

    return output, nil
}

func TestGetFromDatabase(t *testing.T) {

    // a valid Dynamodb's GetItemOutput
    // GetItemOutput is a returned type of Dynamodb's GetItem function
    output := dynamodb.GetItemOutput{}
    output.SetItem(
        map[string]*dynamodb.AttributeValue{
            "id": &dynamodb.AttributeValue{S: aws.String("id_test")},
            "deviceModel": &dynamodb.AttributeValue{S: aws.String("deviceModel_test")},
            "name": &dynamodb.AttributeValue{S: aws.String("name_test")},
            "note": &dynamodb.AttributeValue{S: aws.String("note_test")},
            "serial": &dynamodb.AttributeValue{S: aws.String("serial_test")},
        },
    )

    testCases := []TestCase{
        {
            Name:                   "** Requested id exists **",
            InputIdString:          "id_test",
            ExpectedDatabaseOutput: output,
        },
        {
            Name:                   "** Requested id does not exists **",
            InputIdString:          "id_test_no",
            ExpectedDatabaseOutput: dynamodb.GetItemOutput{},
        },
    }

    // create mocked database.
    getter := new(dynamoDBAPI)
    getter.DynamoDB = &FakeDynamoDBAPI{}


    for _, test := range testCases {

        // calls getDevicebyId.go's Get function.
        response, _ := getter.getFromDatabase(test.InputIdString)

        if len(response.GoString()) != len(test.ExpectedDatabaseOutput.GoString()) {
            t.Errorf("%s \n \t<expected output: \n%s> \n<resulted output: \n%s>", test.Name, test.ExpectedDatabaseOutput.GoString(), response.GoString())
        }
    }
} // end of TestGet function



func TestGetDeviceById(t *testing.T) {

    testCases := []TestCase{
        {
            Name:                "** Testing empty input id **",
            InputId:             events.APIGatewayProxyRequest{PathParameters: map[string]string{
                                        "id": "",},},
            ExpectedBody:         "{\n\t\"error\": {\n\t\t\"code\": 404,\n\t\t\"message\": \"No ID Field Provided\"\n\t}\n}",
            ExpectedStatusCode:   404,
        },
        {
            Name:                "** Testing database internal problem **",
            InputId:             events.APIGatewayProxyRequest{PathParameters: map[string]string{
                                        "id": "id_test",},},
            ExpectedBody:        "{\n\t\"error\": {\n\t\t\"code\": 500,\n\t\t\"message\": \"Internal Server's Error occured\"\n\t}\n}",
            ExpectedStatusCode:   500,
        },
    }

    
    for _, test := range testCases {

        databseStruct = new(types.DatabseStruct)
        databseStruct.TableName = aws.String("test_table_name");
        // calls getDeviceById.go's AddDevice function.
        response,_ := GetDeviceById(test.InputId)

        if response.StatusCode != test.ExpectedStatusCode ||  response.Body != test.ExpectedBody{
            t.Errorf("%s \n \t<expected error-code: %d> <resulted error-code: %d> \n \t<expected body: %s> <resulted body: %s>", test.Name, test.ExpectedStatusCode, response.StatusCode, test.ExpectedBody, response.Body)
        }
    }

} // end of TestAddDevice function





func TestValidateDatabaseResult(t *testing.T) {

    // a valid Dynamodb's GetItemOutput
    // GetItemOutput is a returned type of Dynamodb's GetItem function
    output := dynamodb.GetItemOutput{}
    output.SetItem(
        map[string]*dynamodb.AttributeValue{
            "id": &dynamodb.AttributeValue{S: aws.String("id_test")},
            "deviceModel": &dynamodb.AttributeValue{S: aws.String("deviceModel_test")},
            "name": &dynamodb.AttributeValue{S: aws.String("name_test")},
            "note": &dynamodb.AttributeValue{S: aws.String("note_test")},
            "serial": &dynamodb.AttributeValue{S: aws.String("serial_test")},
            },
    )

    testCases := []TestCase{
        {
            Name:                 "** Database Unexpected Error **",
            InputId:              events.APIGatewayProxyRequest{},
            DatabaseOutput:       dynamodb.GetItemOutput{},
            Error:                errors.New("Unexpected Error has occured"),
            ExpectedBody:         "{\n\t\"error\": {\n\t\t\"code\": 500,\n\t\t\"message\": \"Internal Server's Error occured\"\n\t}\n}",
            ExpectedStatusCode:   500,
        },
        {
            Name:                 "** Database Returns Empty Result **",
            InputId:              events.APIGatewayProxyRequest{},
            DatabaseOutput:       dynamodb.GetItemOutput{},
            ExpectedBody:         "{\n\t\"error\": {\n\t\t\"code\": 404,\n\t\t\"message\": \"Desired device with provided id was not founded\"\n\t}\n}",
            ExpectedStatusCode:   404,
        },
        {
            Name:                 "** Database Returns founded device **",
            DatabaseOutput:       output,
            ExpectedBody:         "{\n\t\"data\": {\n\t\t\"id\": \"id_test\",\n\t\t\"deviceModel\": \"deviceModel_test\",\n\t\t\"name\": \"name_test\",\n\t\t\"note\": \"note_test\",\n\t\t\"serial\": \"serial_test\"\n\t}\n}",
            ExpectedStatusCode:   200,
        },
    }


    for _, test := range testCases {

        response := validateDatabaseResult(&test.DatabaseOutput, test.Error)

        if response.StatusCode != test.ExpectedStatusCode ||  response.Body != test.ExpectedBody{
            t.Errorf("%s \n \t<expected error-code: %d> <resulted error-code: %d> \n \t<expected body: %s> <resulted body: %s>", test.Name, test.ExpectedStatusCode, response.StatusCode, test.ExpectedBody, response.Body)
        }
    }
} // end of TestValidateDatabaseResult function
