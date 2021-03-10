package redis

import "fmt"

func pubChannel(jobID string, topic string) string {
	return fmt.Sprintf("/job/%s/%s", jobID, topic)
}
