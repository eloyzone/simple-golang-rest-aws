package main

import(
	"types"
	"testing"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)


// A fakeDynamoDB instance for mocking test that emulates real DynamoDB
type FakeDynamoDBAPI struct {
	dynamodbiface.DynamoDBAPI
}


// TestCase struct that contains all reuested and expected values for unit testing
type TestCase struct {
	Name 				string
	Request 			events.APIGatewayProxyRequest
	Device 				types.Device
	ExpectedBody 		string
	ExpectedStatusCode 	int
}


// a mocked version of DynamoDB's PutItem function.
// in testing state, instead of calling real DynamoDB's PutItem, we try to emulate it.
// insertItemToDatabase function of addDevice.go calls this function in Testing state.  
func (d *FakeDynamoDBAPI) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return new(dynamodb.PutItemOutput), nil
}


func TestAddDevice(t *testing.T) {

	testCases := []TestCase{
		{
			Name: 				"** Testing empty body input **",
			Request: 			events.APIGatewayProxyRequest{Body: ""},
			ExpectedBody: 		"{\n\t\"error\": {\n\t\t\"code\": 400,\n\t\t\"message\": \"No inputs provided, please provide inputs in json format.\"\n\t}\n}",
			ExpectedStatusCode:	400,
		},
		{
			Name: 				"** Testing wrong json format **",
			Request: 			events.APIGatewayProxyRequest{Body: "{{{}"},
			ExpectedBody: 		"{\n\t\"error\": {\n\t\t\"code\": 400,\n\t\t\"message\": \"Wrong format: Inputs must be a valid json.\"\n\t}\n}",
			ExpectedStatusCode:	400,
		},
		{
			Name:				"** Testing json with missing field {id} **",
			Request: 			events.APIGatewayProxyRequest{Body: "{\"id\":\"\" , \"deviceModel\":\"testDeviceModel\" , \"name\":\"testName\" , \"note\":\"testNote\" , \"serial\":\"testSerial\" }"},
			ExpectedBody: 		"{\n\t\"error\": {\n\t\t\"code\": 400,\n\t\t\"message\": \"Following fields are not provided: id, \"\n\t}\n}",
			ExpectedStatusCode:	400,
		},
		{
			Name:				"** Testing json with missing field {deviceModel, note} **",
			Request:			events.APIGatewayProxyRequest{Body: "{\"id\":\"1\" , \"deviceModel\":\"\" , \"name\":\"testName\" , \"note\":\"\" , \"serial\":\"testSerial\" }"},
			ExpectedBody:		"{\n\t\"error\": {\n\t\t\"code\": 400,\n\t\t\"message\": \"Following fields are not provided: deviceModel, note, \"\n\t}\n}",
			ExpectedStatusCode:	400,
		},
	
		{
			Name: 				"** Testing json with missing field {serial, name, deviceModel} **",
			Request: 			events.APIGatewayProxyRequest{Body: "{\"id\":\"1\" , \"deviceModel\":\"\" , \"name\":\"\" , \"note\":\"testNote\" , \"serial\":\"\" }"},
			ExpectedBody: 		"{\n\t\"error\": {\n\t\t\"code\": 400,\n\t\t\"message\": \"Following fields are not provided: deviceModel, name, serial, \"\n\t}\n}",
			ExpectedStatusCode:	400,
		},

		{
			// as we don't have any access to real database or os.environment, we will get error
			Name:				"** Testing valid json with all fields **",
			Request:			events.APIGatewayProxyRequest{Body: "{\"id\":\"1\" , \"deviceModel\":\"testDeviceModel\" , \"name\":\"testName\" , \"note\":\"testNote\" , \"serial\":\"testSerial\"}"},
			ExpectedBody:		"{\n\t\"error\": {\n\t\t\"code\": 500,\n\t\t\"message\": \"Internal Server's Error occured\"\n\t}\n}",
			ExpectedStatusCode:	500,
		},

	}

    
	for _, test := range testCases {

		// calls addDevice.go's AddDevice function.
		databseStruct = new(types.DatabseStruct)
		databseStruct.TableName = aws.String("test_table_name");
		response, _ := AddDevice(test.Request)

		if response.StatusCode != test.ExpectedStatusCode ||  response.Body != test.ExpectedBody{
			t.Errorf("%s \n \t<expected error-code: %d> <resulted error-code: %d> \n \t<expected body: %s> <resulted body: %s>", test.Name, test.ExpectedStatusCode, response.StatusCode, test.ExpectedBody, response.Body)
		}

	}

} // end of TestAddDevice function

func TestCreateSuccessResponseJson(t *testing.T){

	device := types.Device{
		ID:				"id_test",
		DeviceModel:	"deviceModel_test",
		Name:			"name_test",
		Note:			"note_test",
		Serial:			"serial_test",
	}

	testCases := []TestCase{
		{
			Name:				"** Testing Response **",
			Device:				device,
			ExpectedBody:		"{\n\t\"status\": \"requested item inserted\",\n\t\"data\": {\n\t\t\"id\": \"id_test\",\n\t\t\"deviceModel\": \"deviceModel_test\",\n\t\t\"name\": \"name_test\",\n\t\t\"note\": \"note_test\",\n\t\t\"serial\": \"serial_test\"\n\t}\n}",
			ExpectedStatusCode:	201,
		},
	}


    
	for _, test := range testCases {

		response,_ := createSuccessResponseJson(device)

		if response.StatusCode != test.ExpectedStatusCode ||  response.Body != test.ExpectedBody{
			t.Errorf("%s \n \t<expected error-code: %d> <resulted error-code: %d> \n \t<expected body: %s> <resulted body: %s>", test.Name, test.ExpectedStatusCode, response.StatusCode, test.ExpectedBody, response.Body)
		}
	}

} // end of TestCreateSuccessResponseJson fucntion
