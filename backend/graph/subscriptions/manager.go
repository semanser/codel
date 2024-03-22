package subscriptions

import (
	gmodel "github.com/semanser/ai-coder/graph/model"
)

var (
	taskAddedSubscriptions   = make(map[uint]chan *gmodel.Task)
	flowUpdatedSubscriptions = make(map[uint]chan *gmodel.Flow)
)

type Subscription[T any] interface {
	subscribe() (T, error)
}

func subscribe[B any](flowID uint, subscriptions map[uint]chan B) (channel chan B, unsubscribe func()) {
	ch := make(chan B)

	if _, ok := subscriptions[flowID]; !ok {
		subscriptions[flowID] = ch
	}

	unsubscribe = func() {
		if ch == nil {
			return
		}

		if c, ok := subscriptions[flowID]; ok {
			if c == ch {
				subscriptions[flowID] = nil
			}
		}

		if len(subscriptions[flowID]) == 0 {
			delete(subscriptions, flowID)
		}
	}

	return ch, unsubscribe
}
