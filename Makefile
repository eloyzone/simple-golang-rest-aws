build:
	dep ensure
	env GOOS=linux go build -o bin/handlers/addDevice src/handlers/addDevice/addDevice.go
	env GOOS=linux go build -o bin/handlers/getDeviceById src/handlers/getDeviceById/getDeviceById.go
	env GOOS=linux go build -o bin/handlers/types src/handlers/types/types.go