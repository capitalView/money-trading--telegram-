package utils

import (
	"encoding/json"
)

func ParseJson(jsonStr string) (interface{}, error) {
	var response any
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		return nil, err
	}
	return response, nil
}
