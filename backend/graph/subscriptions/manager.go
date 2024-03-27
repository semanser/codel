package subscriptions

import (
	gmodel "github.com/semanser/ai-coder/graph/model"
)

var (
	taskAddedSubscriptions         = make(map[int64]chan *gmodel.Task)
	flowUpdatedSubscriptions       = make(map[int64]chan *gmodel.Flow)
	terminalLogsAddedSubscriptions = make(map[int64]chan *gmodel.Log)
	browserSubscriptions           = make(map[int64]chan *gmodel.Browser)
)

type Subscription[T any] interface {
	subscribe() (T, error)
}

func subscribe[B any](flowID int64, subscriptions map[int64]chan B) (channel chan B, unsubscribe func()) {
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
