package subscriptions

import (
	gmodel "github.com/semanser/ai-coder/graph/model"
)

func BroadcastTaskAdded(flowID uint, task *gmodel.Task) {
	if ch, ok := taskAddedSubscriptions[flowID]; ok {
		ch <- task
	}
}
