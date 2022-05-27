package util

import (
	"encoding/json"
	"github.com/PaackEng/paackit/errorx"
	"os"
)

// Functions that would be included in paackit

type response struct {
	Data json.RawMessage `json:"data"`
	Err  *responseError  `json:"error"`
}

type responseError struct {
	Msg  string `json:"msg"`
	Code string `json:"code"`
}

func ParseResponse(payload []byte, data interface{}) error {
	var res response

	err := json.Unmarshal(payload, &res)
	if err != nil {
		return err
	}

	if res.Err != nil {
		return errorx.New(res.Err.Msg, res.Err.Code)
	}

	return json.Unmarshal(res.Data, data)
}

func GetEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("missing mandatory environment variable: " + key)
	}
	return value
}
