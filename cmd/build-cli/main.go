package main

import (
	"fmt"
	"io/ioutil"

	schema "github.com/GuyARoss/windtunnel/pkg/windtunnel-schema"
)

func main() {
	config := readConfiguration()

	marshalSchema, schemaParseErr := schema.ParseFile(fmt.Sprintf("%s/%s", config.Schema.BaseDir, config.Schema.Stage))
	if schemaParseErr != nil {
		panic(schemaParseErr)
	}
	schemaParseErr = marshalSchema.ParseFile(fmt.Sprintf("%s/%s", config.Schema.BaseDir, config.Schema.Definition))
	if schemaParseErr != nil {
		panic(schemaParseErr)
	}

	schemaValidationErr := marshalSchema.ValidateDefinitionStageMatch()
	if schemaValidationErr != nil {
		panic(schemaValidationErr)
	}
	fmt.Println(marshalSchema.Stages)

	// generate client wrappers
	//  -- schema defintion to struct
	//  -- setup calling of the client
	//  -- Dockerfilex
}

func readConfiguration() *CompositionConfiguration {
	data, err := ioutil.ReadFile("./composition.yml")
	if err != nil {
		panic(err)
	}

	config := &CompositionConfiguration{}
	err = config.marshal(data)
	if err != nil {
		panic(err)
	}

	validationErr := config.validate()
	if validationErr != nil {
		panic(validationErr)
	}

	return config
}

// fmt.Println(path)

// resp, err := schema.ParseFile(file)
// if err != nil {
// 	panic(err)
// }

// fmt.Println(resp)
