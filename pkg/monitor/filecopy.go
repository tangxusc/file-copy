package monitor

import (
	"context"
	"github.com/rjeczalik/notify"
	"github.com/sirupsen/logrus"
	"github.com/tangxusc/file-copy/pkg/bus"
	"github.com/tangxusc/file-copy/pkg/metrics"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var eventInfoChan = make(chan notify.EventInfo, 100)
var targetDir string
var sourceDir string

func Start(ctx context.Context, source, target string) error {
	targetDir = target
	e, path := getSourceDirPath(source)
	if e != nil {
		panic(e.Error())
	}
	sourceDir = path

	e = copyDirFiles(sourceDir, targetDir)
	if e != nil {
		return e
	}
	//监听文件
	e = notify.Watch(source, eventInfoChan, notify.InMovedFrom, notify.InDelete, notify.InCreate, notify.InMovedTo)
	if e != nil {
		return e
	}
	defer notify.Stop(eventInfoChan)

	handlerEvent(ctx)
	return nil
}

func getSourceDirPath(sou string) (error, string) {
	s, e := filepath.Abs(sou)
	if e != nil {
		return e, s
	}
	return nil, s
}

func handlerEvent(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(eventInfoChan)
			return
		case eve := <-eventInfoChan:
			logrus.Debugf("获取到文件事件:%v", eve)
			e := dispatchEvent(eve)
			if e != nil {
				logrus.Warnf("处理文件出现错误,event:%v,错误:%v", eve, e.Error())
			}
		}
	}
}

func dispatchEvent(info notify.EventInfo) error {
	fileName := getFileName(info.Path())
	join := filepath.Join(targetDir, fileName)
	switch info.Event() {
	//新增
	case notify.InCreate, notify.InMovedFrom:
		e := check(join)
		if e != nil {
			return e
		}
		e = copyDirFiles(info.Path(), join)
		if e != nil {
			return e
		}

	//删除
	case notify.InDelete, notify.InMovedTo:
		e := deleteFile(join)
		if e != nil {
			return e
		}
	}
	return nil
}

func copyDirFiles(source, target string) error {
	info, e := os.Stat(source)
	if e != nil {
		return e
	}
	if info.IsDir() {
		e := os.MkdirAll(target, os.ModePerm)
		if e != nil {
			return e
		}
		infos, e := ioutil.ReadDir(source)
		if e != nil {
			return e
		}
		for _, value := range infos {
			e = copyDirFiles(filepath.Join(source, value.Name()), filepath.Join(target, value.Name()))
			if e != nil {
				return e
			}
		}
	} else {
		return copyFile(target, source)
	}
	return nil
}

func deleteFile(info string) error {
	stat, e := os.Stat(info)
	if e != nil && os.IsNotExist(e) {
		return nil
	}
	if e != nil {
		return e
	}
	if stat.IsDir() {
		infos, e := ioutil.ReadDir(info)
		if e != nil {
			return e
		}
		for _, value := range infos {
			e := deleteFile(filepath.Join(info, value.Name()))
			if e != nil {
				return e
			}
		}
	}

	e = os.Remove(info)
	bus.EventBus <- &metrics.FileDeleteEvent{}
	return e
}

func check(info string) error {
	//check target dir exist
	e, b := checkTargetDirExist()
	if e != nil {
		return e
	}
	if !b {
		e := createTargetDir()
		if e != nil {
			return e
		}
	}

	return nil
}

func copyFile(target, source string) error {
	//check target file exist
	e, b := checkTargetFileExist(target)
	if e != nil {
		return e
	}
	if b {
		logrus.Debugf("文件:%v 已经存在,跳过复制", target)
		return nil
	}
	targetFile, e := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if e != nil {
		return e
	}

	defer targetFile.Close()
	sourceFile, e := os.Open(source)
	if e != nil {
		return e
	}
	defer sourceFile.Close()
	written, e := io.Copy(targetFile, sourceFile)
	if e != nil {
		return e
	}
	logrus.Debugf("拷贝文件完成,共拷贝字节:%v", written)
	bus.EventBus <- &metrics.CountAddEvent{}
	return nil
}

func getFileName(path string) string {
	return strings.Replace(path, sourceDir, "", -1)
}

func createTargetDir() error {
	e := os.MkdirAll(targetDir, os.ModePerm)
	if e != nil {
		return e
	}
	return nil
}

func checkTargetFileExist(path string) (error, bool) {
	_, e := os.Stat(path)
	if e != nil && os.IsNotExist(e) {
		return nil, false
	}
	if e != nil {
		return e, false
	}
	return nil, true
}

func checkTargetDirExist() (error, bool) {
	_, e := os.Stat(targetDir)
	if e != nil && os.IsNotExist(e) {
		return nil, false
	}
	if e != nil {
		return e, false
	}
	return nil, true
}
