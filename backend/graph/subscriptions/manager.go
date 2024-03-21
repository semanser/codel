package subscriptions

import (
	gmodel "github.com/semanser/ai-coder/graph/model"
)

var (
	taskAddedSubscriptions = make(map[uint]chan *gmodel.Task)
)

type Subscription[T any] interface {
	subscribe() (T, error)
}

type SubscriptionType string

var (
	SubscriptionTaskAdded = "taskAdded"
)

func subscribe[B any](flowID uint, subscriptions map[uint]chan B) (channel chan B, unsubscribe func()) {
	ch := make(chan B)

	if _, ok := subscriptions[flowID]; !ok {
		subscriptions[flowID] = make(chan B)
	}

	unsubscribe = func() {
		if ch == nil {
			return
		}

		if c, ok := subscriptions[flowID]; ok {
			if c == ch {
				close(c)
				subscriptions[flowID] = nil
			}
		}

		if len(subscriptions[flowID]) == 0 {
			delete(subscriptions, flowID)
		}
	}

	return ch, unsubscribe
}
