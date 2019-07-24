package command

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tangxusc/file-copy/pkg/bus"
	"github.com/tangxusc/file-copy/pkg/config"
	"github.com/tangxusc/file-copy/pkg/metrics"
	"github.com/tangxusc/file-copy/pkg/monitor"
	"github.com/tangxusc/file-copy/pkg/web"
)

func NewCommand(ctx context.Context) *cobra.Command {
	var command = &cobra.Command{
		Use:   "start",
		Short: "start file copy",
		RunE: func(cmd *cobra.Command, args []string) error {
			//0,绑定参数
			config.Bind()

			if config.Instance.Debug {
				logrus.SetLevel(logrus.DebugLevel)
				logrus.SetReportCaller(true)
				logrus.Debug("已开启debug模式...")
			} else {
				logrus.SetLevel(logrus.WarnLevel)
			}
			//1,eventbus
			bus.Listen(ctx)
			//2,web server
			web.Start(ctx)
			//3,prometheus
			metrics.Start(ctx)
			//4,文件复制
			e := monitor.Start(ctx)
			return e
		},
	}
	logrus.SetFormatter(&logrus.TextFormatter{})
	viper.SetEnvPrefix("FILE")
	viper.AutomaticEnv()
	command.PersistentFlags().BoolVarP(&config.Instance.Debug, "debug", "v", false, "debug mod")
	command.PersistentFlags().StringVarP(&config.Instance.Source, "source", "s", "/source", "source file dir")
	command.PersistentFlags().StringVarP(&config.Instance.Target, "target", "t", "/etc/docker/certs.d/", "target file dir")
	command.PersistentFlags().StringVarP(&config.Instance.Port, "port", "p", "8080", "web server port, example: 8080")
	_ = command.MarkPersistentFlagRequired("source")
	_ = command.MarkPersistentFlagRequired("target")
	_ = viper.BindPFlag("debug", command.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("source", command.PersistentFlags().Lookup("source"))
	_ = viper.BindPFlag("target", command.PersistentFlags().Lookup("target"))
	_ = viper.BindPFlag("port", command.PersistentFlags().Lookup("port"))

	return command
}
