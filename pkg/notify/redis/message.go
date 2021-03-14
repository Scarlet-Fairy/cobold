package redis

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type message struct {
	Error error `json:"error"`
}

func encodeMessageToJson(payload message) (string, error) {
	serialized, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Wrap(err, "could not serialized publish message to json")
	}

	return string(serialized), nil
}
