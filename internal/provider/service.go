package provider

import (
	_ "embed"
	"encoding/json"
)

//go:embed aiservices.json
var aiservicesRaw []byte

var AiServiceEndpoints map[string]string

func init() {
	AiServiceEndpoints = make(map[string]string)
	// fill up the init function to load the aiservices.json file
	if err := json.Unmarshal(aiservicesRaw, &AiServiceEndpoints); err != nil {
		panic("failed to unmarshal aiservices.json: " + err.Error())
	}
}
