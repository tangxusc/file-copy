package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/file-copy/pkg/config"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
)

func Start(ctx context.Context) {
	go func() {
		engine := gin.Default()

		for _, value := range routers {
			value(engine)
		}

		logrus.Debugf("正在启动web服务器,端口:%v", config.Instance.Port)
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", config.Instance.Port),
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
	if config.Instance.Debug {
		go func() {
			logrus.Debugf("已启动性能分析端口:8081,可通过 http://localhost:8081/debug/pprof/ 访问")
			_ = http.ListenAndServe(":8081", nil)
		}()
	}
}

var routers = make([]func(*gin.Engine), 0)

func AddRoute(f func(*gin.Engine)) {
	routers = append(routers, f)
}

func init() {
	AddRoute(func(engine *gin.Engine) {
		engine.GET("/", lsHandler)
	})
}

func lsHandler(ctx *gin.Context) {
	dir := ctx.DefaultQuery("dir", "target")
	sub := ctx.DefaultQuery("sub", "/")
	var lists []*FileItem
	var e error
	switch dir {
	case "target":
		lists, e = showDir(config.Instance.Target, sub)
	case "source":
		lists, e = showDir(config.Instance.Source, sub)
	}
	if e != nil {
		ctx.Data(http.StatusInternalServerError, "", []byte(e.Error()))
		return
	}
	ctx.JSON(http.StatusOK, lists)
}

type FileItem struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"is_dir"`
}

func showDir(dir string, sub string) ([]*FileItem, error) {
	infos, e := ioutil.ReadDir(filepath.Join(dir, sub))
	if e != nil {
		return nil, e
	}
	lists := make([]*FileItem, len(infos))
	for key, value := range infos {
		lists[key] = &FileItem{
			value.Name(),
			value.Size(),
			value.IsDir(),
		}
	}
	return lists, e
}
