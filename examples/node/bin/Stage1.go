package main

type Stage1In struct {
	name String
	Name string
}

func (r *Stage1In) validate() error {

	//@@todo: fill in
	return nil

}

type Stage1Out struct {
	lastName  String
	FirstName string
	LastName  string
	firstName String
}

func (r *Stage1Out) validate() error {

	//@@todo: fill in
	return nil

}

type stage1 struct {
	codeFile string
}

func (r *stage1) invoke(input *Stage1In) Stage1Out {

	// @@todo: fill in
	// - validate input from "input" param
	// - write data to pipe
	// - read data from pipe
	// - marshal pipe data to "output" type
	// - validate output
	// - return output
	return nil

}
