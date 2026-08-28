package transport

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"stagebeam/internal/service"
)

type Stream struct {
	service *service.ShowService
}

func NewStream(showService *service.ShowService) *Stream {
	return &Stream{service: showService}
}

func (stream *Stream) ProcessLine(line []byte) ([]byte, error) {
	if len(line) == 0 {
		return nil, fmt.Errorf("empty command line")
	}
	return ExecuteJSON(stream.service, line)
}

func (stream *Stream) Serve(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 256), 1024*1024)
	encoder := json.NewEncoder(writer)
	for scanner.Scan() {
		response, err := stream.ProcessLine(scanner.Bytes())
		if err != nil {
			if encodeErr := encoder.Encode(Response{OK: false, Message: err.Error()}); encodeErr != nil {
				return encodeErr
			}
			continue
		}
		var value any
		if err := json.Unmarshal(response, &value); err != nil {
			return err
		}
		if err := encoder.Encode(value); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func EncodeCommand(command Command) ([]byte, error) {
	if err := ValidateCommand(command); err != nil {
		return nil, err
	}
	return json.Marshal(command)
}
