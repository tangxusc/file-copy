package bus

import (
	"context"
	"github.com/sirupsen/logrus"
)

var EventBus = make(chan interface{}, 10)
var subscribes = make([]chan interface{}, 0)

func Listen(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(EventBus)
				return
			case event := <-EventBus:
				dispatchEvent(event)
			}
		}
	}()
}

func dispatchEvent(event interface{}) {
	logrus.Debugf("收到事件:%v", event)
	for _, value := range subscribes {
		go func(ch chan interface{}) {
			ch <- event
			logrus.Debugf("已发送事件完成:%v", event)
		}(value)
	}
	logrus.Debugf("事件已发送到所有chan:%v", event)
}

func Register() <-chan interface{} {
	ch := make(chan interface{}, 10)
	subscribes = append(subscribes, ch)
	return ch
}
