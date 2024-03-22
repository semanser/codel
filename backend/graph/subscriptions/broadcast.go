package subscriptions

import (
	gmodel "github.com/semanser/ai-coder/graph/model"
)

func BroadcastTaskAdded(flowID uint, task *gmodel.Task) {
	if ch, ok := taskAddedSubscriptions[flowID]; ok {
		ch <- task
	}
}

func BroadcastFlowUpdated(flowID uint, flow *gmodel.Flow) {
	if ch, ok := flowUpdatedSubscriptions[flowID]; ok {
		ch <- flow
	}
}
