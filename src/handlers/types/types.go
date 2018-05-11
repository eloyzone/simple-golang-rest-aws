package types


// struct that contains device information, as json
type Device struct {
    ID          string  `json:"id"`
    DeviceModel string  `json:"deviceModel"`
    Name        string  `json:"name"`
    Note  		string  `json:"note"`
    Serial   	string  `json:"serial"`
}


// struct that contains errors for showing to clinet, as json
type ErrorResponse struct {
   ErrorMessage   ErrorMessage    `json:"error"`
}

type ErrorMessage struct {
   Code   int     `json:"code"`
   Message string  `json:"message"`
}

type DatabseStruct struct {
    SessionError error
    TableName *string
}