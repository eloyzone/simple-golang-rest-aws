# Simple Golang Rest API on AWS

A simple implementation of Rest-API on AWS that uses following tech stack:

* [Serverless Framework] 
* [Go Programming language]
* AWS API Gateway
* AWS Lambda
* AWS DynamoDB



## Project structure

This project is based on AWS's [serverless architecture]. Serverless doesn't mean that we don't have any server, it means that we don't care about server's manipulation and we just pay for processing time, not for server's idle time.

Provided codes have two main responsibility, `adding` a new device to database and `fetching` a device based on its id. To achieve this goal we've created two AWS's Lambda functions namely `addNewDevice`, `getDeviceById`.

AWS provides various programming options for creating lambda functions like Java, C# and etc. In this project we've used Golang which is recently added to AWS's supported programming languages list.  


##### Why Serverless Framework?
AWS has created different web-consoles for developing and deploying your projects but as you project get bigger and bigger it's seems that using these web-consoles just waste our time.

Serverless Framework is a toolkit for deploying and operating serverless architectures, by using it you can enhance your programming speed and most of tough and repetitive tasks will be automatically.

The main file of Serverlesss Framework is `serverless.yml` that contains most of your configurations like `service-name`, `database-name`, `table-name`, `permissions`, `aws-region`, `function-name` and etc.






## Rest-API's description

The API accepts the following JSON requests and produces the corresponding HTTP responses:

##### Request 1:

Requests to insert a new device to database(DynamoDB).

```
HTTP Method: POST
URL: https://<api-gateway-url>/api/devices
content-type: application/json
Body:
{
  "id": "/devices/id1",
  "deviceModel": "/devicemodels/id1",
  "name": "Sensor",
  "note": "Testing a sensor.",
  "serial": "A020000102"
}
```

##### Response 1 - Success:
Provided data inserted to database(DynamoDB) successfully.

```
HTTP-Statuscode: HTTP 201
content-type: application/json
Body:
{
	"status": "requested item inserted",
	"data": {
		"id": "/devices/id1",
		"deviceModel": "/devicemodels/id1",
		"name": "Sensor",
		"note": "Testing a sensor.",
		"serial": "A020000102"
	}
}

```

##### Response 1 - Failure 1:
If any of the payload fields are missing. Response will have a descriptive error message for client user.

```
HTTP-Statuscode: HTTP 400
content-type: application/json
body:
{
	"error": {
		"code": 400,
		"message": "Following fields are not provided: id, serial "
	}
}

```


##### Response 1 - Failure 2:
If any exceptional situation occurs on the server side.

```
HTTP-Statuscode: HTTP 500
content-type: application/json
body:
{
	"error": {
		"code": 500,
		"message": "Internal Server's Error occured, serial "
	}
}
```

##### Request 2:
Get a device based on provided id.

```
HTTP Method: GET
URL: https://`API-GATEWAY-URL`/api/devices/{id}

Example: https://api123.amazonaws.com/api/devices/id1
```

##### Response 2 - Success:

```
HTTP-Statuscode: HTTP 200
content-type: application/json
body:
{
	"data": {
		"id": "/devices/id1",
		"deviceModel": "/devicemodels/id1",
		"name": "Sensor",
		"note": "Testing a sensor.",
		"serial": "A020000102"
	}
}

```

##### Response 2 - Failure 1:
Requested device with provided id not founded.

```
{
	"error": {
		"code": 404,
		"message": "Desired device with provided id was not founded"
	}
}
```

##### Response 2 - Failure 2:
If any exceptional situation occurs on the server side.

```
{
	"error": {
		"code": 500,
		"message": "Internal Server's Error occured"
	}
}
```

These JSON structured is suggested by [Google JSON Guideline]


## Getting Started

In order to use these code you have to install some applications and having one AWS's account is necessary.

### Prerequisites

Although you can use Windows or Mac as your local machine but it's highly recommended to use GNU-Linux. Following installation steps are based on [Fedora] distribution.



### Installing

##### Installing NodeJS and npm

We don't use nodejs as a programming language we just need to have npm for installing our serverless framework.

```
curl --silent --location https://rpm.nodesource.com/setup_8.x | sudo bash -

sudo yum -y install nodejs
```
For other operating systems check [NodeJs] website.





##### Installing Serverless Framework

```
npm install -g serverless
```

To run serverless commands that interface with your AWS account, you will need to setup your AWS account credentials on your machine.

```
sls config credentials --provider aws --key aws_access_key_id --secret aws_secret_access_key
```
Replace aws_access_key_id and aws_secret_access_key with your AWS's API Key & Secret.

For more information check serverless installation and credentials links.



##### Installing golang

Following commands work fine on Fedora for more information check [Go Programming language] guidelines.

```
sudo dnf install golang
mkdir -p $HOME/go
echo 'export GOPATH=$HOME/go' >> $HOME/.bashrc
source $HOME/.bashrc

```

##### Installing dep

dep is a prototype dependency management tool for Go. It requires Go 1.9 or newer to compile. dep is safe for production use. 

Following command works fine on Fedora for more information [dep]'s github.

```
sudo dnf install dep
```


##### Cloning path [Important]

As golang and dep are new stack techs in some environment it might be a little bit bugy so it's imprtant to clone this project to follwing path

```
cd  ${GOPATH}/src
```

If you have any more problem in this step, check [GOPATH].


##### Installing wkhtmltopdf - [ optional ]

In order to have a `pdf` version of our `testing` report we use wkhtmltopdf application. It's an optional application but it's recommended as our `test.sh` uses this file.

```
sudo dnf install wkhtmltopdf
```


## Deployment
The first easiest way for deploying is compiling your `.go` files then make a `.zip` archive and upload it to AWS, but doing these steps are very tedious so we want to use serverless framework to automates these essentials steps for us.

Don't forget by using Serverless Framework you can automate most of steps like creating DynamoDB's table and assigning permission to it.

In order to build (compile) your project and deploy it to AWS, use following commands

```
make
sls deploy
```

`make` command calls Makefile that contains some scripts for resolving golang dependencies and compiling `.go` files. `sls deploy` deploy our build project to AWS.


### Deployment - scricpts

If we add new lambda function to our project we would need to update Makefile and add those file to it. In order to make these progress more easier you can use `.build` file.

```
./scripts/build.sh
```

For deploying you can use following command which builds your code and then deploy it to AWS 
```
./scripts/deploy.sh
```

## Testing
After deploying, AWS gives you two links, one for adding new device and one for getting a device by its id. (follwing links are just sample)

```
POST https://API-GATEWAY-URL/dev/devices

GET https://API-GATEWAY-URL/dev/devices/{id}
```

###### Test - Create Device:

```
curl -i -H "Content-Type: application/json" -X POST https://API-GATEWAY-URL/devices -d '{"id":"/devices/id1","deviceModel":"/devicemodels/id1","name":"Sensor","note":"Testing a sensor.","serial":"A020000102"}' 
```

###### Test - Get Device:

```
curl -i https://API-GATEWAY-URL/devices/13
```



## Unit Testing
Go has a built-in testing command called `go test` and a package testing which combine to give a minimal but complete testing experience. 

All files that ends to `_test.go` are considered as testing codes that can be founded in each package.

For running `unit-testing` and `coverage-test` you can use `.test.sh` file.

```
./scripts/test.sh
```

This scripts will run three different command for each packages. these commands test each `.go` file, create reports in `cover.out`, `cover.html` and `cover.pdf`

```
go test -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html
wkhtmltopdf cover.html cover.pdf
```


[Serverless Framework]: https://serverless.com/
[Go Programming language ]: https://golang.org/
[serverless architecture]: https://martinfowler.com/articles/serverless.html
[Google JSON Guideline]: https://google.github.io/styleguide/jsoncstyleguide.xml
[Fedora]: https://getfedora.org/
[NodeJs]: https://nodejs.org/en/download/
[installation]: https://serverless.com/framework/docs/providers/aws/guide/installation/
[credentials]: https://serverless.com/framework/docs/providers/aws/guide/credentials/
[dep]: https://github.com/golang/dep
[GOPATH]: https://github.com/golang/go/wiki/GOPATH
