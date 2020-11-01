package utils

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type ResponseJson struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
	Status     int         `json:"status"`
	Message    string      `json:"message"`
}

func MarshalJson(v interface{}) string {
	var json = jsoniter.ConfigFastest
	buf, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Json Marshall Error. : " + err.Error())
		panic(err)
	}
	// fmt.Printf("buf is", string(buf))
	return string(buf)
}
