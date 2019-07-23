package main

import (
	"context"
	"github.com/tangxusc/file-copy/pkg/command"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	comm := command.NewCommand(ctx)
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		cancel()
	}()
	e := comm.Execute()
	if e != nil {
		panic(e.Error())
	}
}
