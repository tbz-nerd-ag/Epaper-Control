package rest

// ErrorResponse
type ErrorResponse struct {
	Error string `json:"error" example:"Ungültiger Token"`
}

// DeviceResponse
type DeviceResponse struct {
	Device string `json:"device"`
}

type badrequest struct {
	Error string `json:"error" example:"Ungültige Anfrage"`
}

type WartungResponse struct {
	Wartung bool `json:"wartung" example:"false"`
}

type FileNameResponse struct {
	FileName string `json:"log_filename" example:"/root/logs/log.log"`
}

type WartungSleepResponse struct {
	WartungSleepTime int `json:"wartung_sleep_time" example:"20"`
}
