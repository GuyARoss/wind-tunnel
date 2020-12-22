package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("./composition.yml")
	if err != nil {
		panic(err)
	}

	config := &CompositionConfiguration{}
	err = config.marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(config)

	// generate client wrappers
	//  -- schema defintion to struct
	//  -- setup calling of the client
	//  -- Dockerfilex
}

// fmt.Println(path)

// // validate schema
// file, err := os.Open("./examples/node/schema/definitions.windtunnel")
// if err != nil {
// 	log.Fatal(err)
// }
// defer file.Close()

// resp, err := schema.ParseFile(file)
// if err != nil {
// 	panic(err)
// }

// fmt.Println(resp)
