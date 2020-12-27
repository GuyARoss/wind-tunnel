package bundle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

type xxx_Pipein struct{}

type xxx_Pipeout struct{}

type pipeEnv struct {
	pipeName string
}

func (r *pipeEnv) writeToPipe(input *xxx_Pipein) error {
	serialOut, serialErr := json.Marshal(input)
	if serialErr != nil {
		return serialErr
	}

	syscall.Mkfifo(r.pipeName, 0666)

	file, fileErr := os.OpenFile(r.pipeName, os.O_RDWR, os.ModeNamedPipe)
	defer file.Close()

	if fileErr != nil {
		return fileErr
	}

	file.WriteString(fmt.Sprintf("%s\n", string(serialOut)))

	return nil
}

func (r *pipeEnv) readFromPipe() (*xxx_Pipeout, error) {
	syscall.Mkfifo(r.pipeName, 0666)

	file, err := os.OpenFile(r.pipeName, os.O_RDONLY, os.ModeNamedPipe)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	data, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return nil, readErr
	}

	deserial := &xxx_Pipeout{}
	deserialErr := json.Unmarshal(data, deserial)

	return deserial, deserialErr
}

func startUpStage() error {
	return nil
}
