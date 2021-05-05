package notify

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Message struct {
	Topic string `json:"topic"`
	Error string `json:"error"`
}

func EncodeMessageToJson(payload Message) ([]byte, error) {
	serialized, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "could not serialized publish message to json")
	}

	return serialized, nil
}
