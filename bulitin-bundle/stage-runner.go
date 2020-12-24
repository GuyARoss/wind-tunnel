package bundle

// @@ use named pipe apply input data, return output

type stage interface {
	Prestage(input interface{}) (interface{}, error)
	PostStage(input interface{}) (interface{}, error)

	PrimaryStage(input interface{}) (interface{}, error)
}
