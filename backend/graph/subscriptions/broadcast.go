package subscriptions

import (
	gmodel "github.com/semanser/ai-coder/graph/model"
)

func BroadcastTaskAdded(flowID int64, task *gmodel.Task) {
	if ch, ok := taskAddedSubscriptions[flowID]; ok {
		ch <- task
	}
}

func BroadcastFlowUpdated(flowID int64, flow *gmodel.Flow) {
	if ch, ok := flowUpdatedSubscriptions[flowID]; ok {
		ch <- flow
	}
}
