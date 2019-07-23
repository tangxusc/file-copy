package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tangxusc/file-copy/pkg/bus"
	"github.com/tangxusc/file-copy/pkg/metrics"
	"github.com/tangxusc/file-copy/pkg/monitor"
	"github.com/tangxusc/file-copy/pkg/web"
)

var source string
var target string
var debug bool
var port string

func NewCommand(ctx context.Context) *cobra.Command {
	var command = &cobra.Command{
		Use:   "start",
		Short: "start file copy",
		RunE: func(cmd *cobra.Command, args []string) error {
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
				logrus.SetReportCaller(true)
				logrus.Debug("已开启debug模式...")
			} else {
				logrus.SetLevel(logrus.WarnLevel)
			}
			//0,eventbus
			bus.Listen(ctx)
			//1,web server
			web.Start(ctx, port)

			//3,prometheus
			metrics.Start(ctx)

			//2,文件复制
			e := monitor.Start(ctx, source, target)
			return e
		},
	}
	logrus.SetFormatter(&logrus.TextFormatter{})
	command.PersistentFlags().BoolVarP(&debug, "debug", "v", false, "debug mod")
	command.PersistentFlags().StringVarP(&source, "source", "s", "/source", "source file dir")
	command.PersistentFlags().StringVarP(&target, "target", "t", "/etc/docker/certs.d/", "target file dir")
	command.PersistentFlags().StringVarP(&port, "port", "p", "8080", "web server port, example: 8080")
	_ = viper.BindPFlag("debug", command.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("source", command.PersistentFlags().Lookup("source"))
	_ = viper.BindPFlag("target", command.PersistentFlags().Lookup("target"))
	_ = viper.BindPFlag("port", command.PersistentFlags().Lookup("port"))
	_ = command.MarkPersistentFlagRequired("source")
	_ = command.MarkPersistentFlagRequired("target")

	return command
}
