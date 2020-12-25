package bundle

import (
	"os"
	"syscall"
)

type stage interface {
	Prestage(input interface{}) (interface{}, error)
	PostStage(input interface{}) (interface{}, error)

	PrimaryStage(input interface{}) (interface{}, error)
}

type TransitData struct {
}

func writeToPipe() (*TransitData, error) {
	syscall.Mkfifo("tmpPipe", 0666)

	file, writeErr := os.OpenFile("tmpPipe", os.O_RDWR, os.ModeNamedPipe)
	if writeErr != nil {

	}

	file, err := os.OpenFile("tmpPipe", os.O_RDONLY, os.ModeNamedPipe)
}

func startUpStage() {

}
