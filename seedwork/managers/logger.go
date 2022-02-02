package managers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type LoggerManager struct {
	BaseAddress            string
	ErrorEndPointUri       string
	InformationEndPointUri string
	Source                 string
	WarningEndPointUri     string
}

type loggerDTO struct {
	Source    string `json:"source"`
	Message   string `json:"message"`
	Trace     string `json:"trace"`
	UserAgent string `json:"userAgent"`
	Username  string `json:"username"`
}

func NewLoggerManager(baseAddress string, source string) LoggerManager {

	return LoggerManager{
		BaseAddress:            baseAddress,
		ErrorEndPointUri:       "api/v1/events/errors",
		InformationEndPointUri: "api/v1/events/informations",
		Source:                 source,
		WarningEndPointUri:     "api/v1/events/warnings",
	}
}

func (log LoggerManager) Error(message string, trace string, username string, userAgent string) string {
	return log.logger(log.ErrorEndPointUri, message, trace, username, userAgent)
}

func (log LoggerManager) logger(endPoint string, message string, trace string, username string, userAgent string) string {
	dto := loggerDTO{log.Source, message, trace, userAgent, username}

	jsonReq, err := json.Marshal(dto)

	client := &http.Client{}

	responseBody := bytes.NewBuffer(jsonReq)
	req, err := http.NewRequest(http.MethodPost, log.BaseAddress+endPoint, responseBody)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	req.Header.Add("Authorization", "Basic dXNybG9nZ2luZzpBMTIzNDU2YQ==")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println(message)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println(message)
		return ""
	}

	return string(body)
}
