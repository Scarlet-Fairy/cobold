package docker

import (
	"bufio"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
)

type Output struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	Payload string `json:"payload"`
}

func serializeOutputStreamBuffer(buff io.Reader) ([]byte, error) {
	streams := make([]Stream, 0)
	rd := bufio.NewReader(buff)

	for {
		line, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "could not read build output")
		}

		streams = append(streams, Stream{
			Payload: string(line),
		})
	}

	outStruct := Output{
		Streams: streams,
	}
	stringifiedOutput, err := json.Marshal(outStruct)
	if err != nil {
		return nil, errors.Wrap(err, "could not stringify the build output")
	}

	return stringifiedOutput, nil
}
