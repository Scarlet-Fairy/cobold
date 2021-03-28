package redis

import "fmt"

func pubChannel(jobID string) string {
	return fmt.Sprintf("/job/%s", jobID)
}
