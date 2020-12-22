package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	schema "github.com/GuyARoss/windtunnel/pkg/windtunnel-schema"
)

func main() {
	var path string
	flag.StringVar(&path, "p", "", "")

	fmt.Println(path)

	// validate schema
	file, err := os.Open("./examples/node/schema/definitions.windtunnel")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	resp, err := schema.ParseFile(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	// generate client wrappers
	//  -- schema defintion to struct
	//  -- setup calling of the client
	//  -- Dockerfilex
}
