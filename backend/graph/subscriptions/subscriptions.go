package subscriptions

import (
	"context"

	gmodel "github.com/semanser/ai-coder/graph/model"
)

func TaskAdded(ctx context.Context, flowId int64) (<-chan *gmodel.Task, error) {
	ch, unsubscribe := subscribe(flowId, taskAddedSubscriptions)

	go func() {
		// Handle deregistration of the channel here. Note the `defer`
		defer func() {
			unsubscribe()
		}()

		for {
			<-ctx.Done() // This runs when context gets cancelled. Subscription closes.
			// Handle deregistration of the channel here. `close(ch)`
			return
		}
	}()

	return ch, nil
}

func FlowUpdated(ctx context.Context, flowId int64) (<-chan *gmodel.Flow, error) {
	ch, unsubscribe := subscribe(flowId, flowUpdatedSubscriptions)
	go func() {
		defer func() {
			unsubscribe()
		}()
		for {
			<-ctx.Done()
			return
		}
	}()
	return ch, nil
}

func TerminalLogsAdded(ctx context.Context, flowId int64) (<-chan *gmodel.Log, error) {
	ch, unsubscribe := subscribe(flowId, terminalLogsAddedSubscriptions)
	go func() {
		defer func() {
			unsubscribe()
		}()
		for {
			<-ctx.Done()
			return
		}
	}()
	return ch, nil
}

func BrowserUpdated(ctx context.Context, flowId int64) (<-chan *gmodel.Browser, error) {
	ch, unsubscribe := subscribe(flowId, browserSubscriptions)
	go func() {
		defer func() {
			unsubscribe()
		}()
		for {
			<-ctx.Done()
			return
		}
	}()
	return ch, nil
}
