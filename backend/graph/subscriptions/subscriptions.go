package subscriptions

import (
	"context"

	gmodel "github.com/semanser/ai-coder/graph/model"
)

func TaskAdded(ctx context.Context, flowId uint) (<-chan *gmodel.Task, error) {
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
