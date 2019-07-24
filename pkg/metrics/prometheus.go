package metrics

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/file-copy/pkg/bus"
	"github.com/tangxusc/file-copy/pkg/web"
)

var counter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "file",
	Subsystem: "copy",
	Name:      "total",
	Help:      "file total count",
}, nil)

var current = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "file",
	Subsystem: "count",
	Name:      "current",
	Help:      "current file count",
}, nil)

func init() {
	web.AddRoute(func(engine *gin.Engine) {
		engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	})
	prometheus.MustRegister(counter)
	prometheus.MustRegister(current)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func Start(ctx context.Context) {
	register := bus.Register()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-register:
				handlerEvent(event)
			}
		}
	}()
}

func handlerEvent(event interface{}) {
	logrus.Debugf("metrics 收到事件:%v", event)
	switch event.(type) {
	case *CountAddEvent:
		counter.With(nil).Add(1)
	case *FileDeleteEvent:
		current.With(nil).Sub(1)
	}
}
