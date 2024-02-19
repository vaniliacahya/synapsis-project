package helper

import (
	"github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

func LogrusDefiner() {
	formatter := runtime.Formatter{ChildFormatter: &logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	}}
	formatter.Line = true
	logrus.SetFormatter(&formatter)
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}
