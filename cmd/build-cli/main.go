package main

import (
	"fmt"
	"io/ioutil"

	"github.com/GuyARoss/windtunnel/pkg/schema"
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

	// generate client wrappers
	//  -- schema defintion to struct
	//  -- setup calling of the client
	//  -- Dockerfilex
}

func readConfiguration() *CompositionConfiguration {
	// @@cli use flag for path
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
