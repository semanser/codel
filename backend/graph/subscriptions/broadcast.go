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

func BroadcastTerminalLogsAdded(flowID int64, flow *gmodel.Log) {
	if ch, ok := terminalLogsAddedSubscriptions[flowID]; ok {
		ch <- flow
	}
}

func BroadcastBrowserUpdated(flowID int64, browser *gmodel.Browser) {
	if ch, ok := browserSubscriptions[flowID]; ok {
		ch <- browser
	}
}
