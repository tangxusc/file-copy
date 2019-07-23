package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Start(ctx context.Context, port string) {
	go func() {
		engine := gin.Default()

		for _, value := range routers {
			value(engine)
		}

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: engine,
		}

		go func() {
			select {
			case <-ctx.Done():
				e := srv.Shutdown(ctx)
				if e != nil {
					panic(e.Error())
				}
			}
		}()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

var routers = make([]func(*gin.Engine), 0)

func AddRoute(f func(*gin.Engine)) {
	routers = append(routers, f)
}
