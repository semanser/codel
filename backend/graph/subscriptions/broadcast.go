package subscriptions

import (
	"context"

	gmodel "github.com/semanser/ai-coder/graph/model"
)

func BroadcastTaskAdded(ctx context.Context, task *gmodel.Task) {
	if ch, ok := taskAddedSubscriptions[task.ID]; ok {
		ch <- task
	}
}
